// Package katago provides wrappers for analyzing games with katago.
package katago

import (
	"bufio"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
)

const katagoReadyStr = "Started, ready to begin handling requests"

const collectorTimeout = 2 * time.Second
const maxCollectorTimeout = 30 * time.Second

// Analyzer is a katago-analyzer.
type Analyzer struct {
	// Model is the path to the model file.
	Model string

	// Config is the path to the config file.
	Config string

	// Number of analysis threads. Defaults to 16 if not specified.
	AnalysisThreads int

	cmd *exec.Cmd

	// Standard IO Processing for child preocess
	stdoutQuit   chan int
	stderrQuit   chan int
	stdinQuit    chan int
	stdinWrite   chan string
	katagoOutput chan string

	// collectors aggregate output from Katago
	collectorQuit          chan int
	analysisSizeMap        sync.Map
	outputResultCollection map[string]*AnalysisList
	outputResultCondition  *sync.Cond

	// Wait group to wait until Katago has booted and proces the katagoReadyStr
	bootWait sync.WaitGroup
}

var analyzerSingleton *Analyzer
var once sync.Once

// GetAnalyzer returns the Katago Analyzer Singleton
func GetAnalyzer(model string, configPath string, numThreads int) *Analyzer {
	once.Do(func() {
		threads := numThreads
		if numThreads == 0 {
			threads = 8
		}
		outputResultMutex := sync.Mutex{}
		analyzerSingleton = &Analyzer{
			Model:                  model,
			Config:                 configPath,
			AnalysisThreads:        threads,
			stdoutQuit:             make(chan int),
			stderrQuit:             make(chan int),
			stdinQuit:              make(chan int),
			stdinWrite:             make(chan string),
			katagoOutput:           make(chan string),
			collectorQuit:          make(chan int),
			outputResultCollection: make(map[string]*AnalysisList),
			outputResultCondition:  sync.NewCond(&outputResultMutex),
		}

		analyzerSingleton.createCmd()
		glog.Infof("using model %q", analyzerSingleton.Model)
		glog.Infof("using gtp config %q", analyzerSingleton.Config)
	})
	return analyzerSingleton
}

// AnalyzeGame is a synchronous request for Katago to process a game
func (an *Analyzer) AnalyzeGame(q *Query) (*AnalysisList, error) {
	an.bootWait.Wait()
	// We need to know the number of analyzed moves..
	json, err := q.ToJSON()
	if err != nil {
		return nil, err
	}
	jsonStr := string(json)
	if jsonStr[len(jsonStr)-1:] != "\n" {
		glog.Warningf("Appending newline to input: %s", jsonStr)
		jsonStr = jsonStr + "\n"
	}

	// Run the Katago Analysis and wait for the result
	an.analysisSizeMap.Store(q.ID, len(q.AnalyzeTurns))
	an.stdinWrite <- jsonStr
	analysis := an.readCollectedResult(q.ID)
	an.analysisSizeMap.Delete(q.ID)
	return analysis, nil
}

// writeCollectedResult writes results to the synchronized map and broadcasts.
func (an *Analyzer) writeCollectedResult(analysisID string, al *AnalysisList) {
	glog.Infof("writeCollectedResult Outputting AnalysisList for %s", analysisID)
	an.outputResultCondition.L.Lock()
	an.outputResultCollection[analysisID] = al
	an.outputResultCondition.L.Unlock()
	an.outputResultCondition.Broadcast()
}

// readCollectedResult waits for a completed and compiled AnalysisList and pops it.
func (an *Analyzer) readCollectedResult(analysisID string) *AnalysisList {
	glog.Infof("Read collector started")
	an.outputResultCondition.L.Lock()
	_, ok := an.outputResultCollection[analysisID]
	for !ok {
		an.outputResultCondition.Wait()
		_, ok = an.outputResultCollection[analysisID]
		glog.Infof("readCollectedResult Woke! %v", an.outputResultCollection[analysisID])
	}
	glog.Infof("readCollectedResult Was notified and found AnalysisList for %s", analysisID)
	result := an.outputResultCollection[analysisID]
	delete(an.outputResultCollection, analysisID)
	an.outputResultCondition.L.Unlock()
	return result
}

// collector is a go-func for
func (an *Analyzer) resultCollector() {
	// firstResultProcessed := false
	// var output string
	resultStorage := make(map[string]AnalysisList)
	for {
		select {
		case <-an.collectorQuit:
			// Great! Timeout happened; stop collecting because we can.
			glog.Info("resultCollector: Received shutdown signal - returning")
			return
		case output := <-an.katagoOutput:
			// firstResultProcessed = true
			if strings.Contains(output, "error") {
				glog.Warningf("resultCollector: Katago had an error: %v", output)
			} else {
				glog.Infof("Collected %s", output)
				res, err := ParseAnalysis(output)
				if err != nil {
					glog.Warningf("resultCollector: Error decoding analysis: %s", output)
				}

				// Add the current result to the AnalysisList
				al := resultStorage[res.ID]
				resultStorage[res.ID] = append(al, res)

				// If we've collected all the results for this analysis run, notify the caller.
				sizeLookup, ok := an.analysisSizeMap.Load(res.ID)
				if ok {
					expectedNumResults := sizeLookup.(int)
					actualNumResults := len(resultStorage[res.ID])
					glog.Infof("resultCollector: Found %d/%d expected analysis results for %s",
						actualNumResults, expectedNumResults, res.ID)
					if actualNumResults >= expectedNumResults {
						glog.Infof("resultCollector: Found sufficent results to return for %s", res.ID)
						resultSlice := resultStorage[res.ID]
						an.writeCollectedResult(res.ID, &resultSlice)

						delete(resultStorage, res.ID)

					}
				} else {
					glog.Warningf("resultCollector: Expected to find number of results for analysis run, but did not. Run ID: %s", res.ID)
				}
			}
		}
	}
}

// Start the Katago Analyzer and make it ready to receive input
func (an *Analyzer) Start() error {
	glog.Info("Starting katago analyzer")
	analyzerSingleton.bootWait.Add(1)

	go an.resultCollector()
	err := an.startAnalyzerIO()
	if err != nil {
		return err
	}
	err = an.cmd.Start()
	if err != nil {
		return err
	}
	analyzerSingleton.bootWait.Wait()

	glog.Info("Katago Startup Complete")
	return nil
}

// Stop the Katago Analyzer and the goroutines
func (an *Analyzer) Stop() error {
	glog.Infof("Shutting down Katago!")
	an.stdinQuit <- 1
	an.collectorQuit <- 1
	return nil
}

// Cmd creates the Katago analysis command.
func (an *Analyzer) createCmd() {
	threads := an.AnalysisThreads
	an.cmd = exec.Command("katago", "analysis", "-model", an.Model, "-config", an.Config, "-analysis-threads", strconv.Itoa(threads))
}

func (an *Analyzer) startAnalyzerIO() error {
	// Note: Pipes must be created before the command is run.
	stderr, err := an.cmd.StderrPipe()
	if err != nil {
		return err
	}
	go an.stdErrReader(stderr)

	stdout, err := an.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	go an.stdOutReader(stdout)

	stdin, err := an.cmd.StdinPipe()
	if err != nil {
		return err
	}
	go an.stdInWriter(stdin)
	return nil
}

func (an *Analyzer) stdErrReader(stderr io.ReadCloser) {
	defer stderr.Close()
	glog.Infof("katago stderr: started")
	scanner := bufio.NewScanner(stderr)
	var currentText string
	for scanner.Scan() {
		currentText = scanner.Text()

		if strings.Contains(currentText, katagoReadyStr) {
			glog.Infof("katago stderr: Katago ready-string found; ready to process input.")
			an.bootWait.Done()
		}
		glog.Infof("katago stderr: %v\n", currentText)
	}
	glog.Infof("katago stderr: returned")
}

func (an *Analyzer) stdOutReader(stdout io.ReadCloser) {
	defer stdout.Close()
	glog.Infof("katago stdout: started")
	scanner := bufio.NewScanner(stdout)
	var currentText string
	for scanner.Scan() {
		currentText = scanner.Text()
		glog.Infof("katago stdout: Sending result back from reader: %v\n", currentText)
		an.katagoOutput <- currentText
	}
	glog.Infof("katago stdout: returned")
}

func (an *Analyzer) stdInWriter(stdin io.WriteCloser) {
	defer stdin.Close()
	glog.Infof("katago stdin: started")
	for {
		select {
		case <-an.stdinQuit:
			// Make the go-routine shut-down
			// This will also close stdin, which will shut down Katago.
			glog.Infof("katago stdin: returning")
			return
		case writeValue := <-an.stdinWrite:
			glog.Info("katago stdin: Got value to write to Katago!")
			stdin.Write([]byte(writeValue))
		}
	}
}
