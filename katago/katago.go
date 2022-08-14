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
	model string

	// Config is the path to the config file.
	config string

	// Number of analysis threads. Defaults to 16 if not specified.
	analysisThreads int

	cmd *exec.Cmd

	// Standard IO Processing for child preocess
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

// New creates a new analyzer. Before being used, the Start() must be called.
func New(model string, configPath string, analysisThreads int) *Analyzer {
	an := &Analyzer{
		model:                  model,
		config:                 configPath,
		analysisThreads:        analysisThreads,
		stdinQuit:              make(chan int),
		stdinWrite:             make(chan string),
		katagoOutput:           make(chan string),
		collectorQuit:          make(chan int),
		outputResultCollection: make(map[string]*AnalysisList),
		outputResultCondition:  sync.NewCond(&sync.Mutex{}),
	}
	an.cmd = an.createCmd()
	return an
}

// Start the Katago Analyzer and ensures that it's ready to receive input.
func (an *Analyzer) Start() error {
	glog.Info("Starting Katago analyzer")
	glog.Infof("Using model %q", an.model)
	glog.Infof("Using gtp config %q", an.config)
	an.bootWait.Add(1)

	go an.resultCollector()
	err := an.startAnalyzerIO()
	if err != nil {
		return err
	}

	glog.V(2).Infof("Executing katago command: %v", an.cmd)

	err = an.cmd.Start()
	if err != nil {
		return err
	}

	an.bootWait.Wait()
	glog.Info("Katago Startup Complete")
	return nil
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
	glog.V(3).Infof("writeCollectedResult Outputting AnalysisList for %s", analysisID)
	an.outputResultCondition.L.Lock()
	an.outputResultCollection[analysisID] = al
	an.outputResultCondition.L.Unlock()
	an.outputResultCondition.Broadcast()
}

// readCollectedResult waits for a completed and compiled AnalysisList and pops it.
func (an *Analyzer) readCollectedResult(analysisID string) *AnalysisList {
	glog.V(2).Infof("Read collector started")

	an.outputResultCondition.L.Lock()
	_, ok := an.outputResultCollection[analysisID]
	for !ok {
		an.outputResultCondition.Wait()
		_, ok = an.outputResultCollection[analysisID]
		glog.V(3).Infof("readCollectedResult Woke! %v", an.outputResultCollection[analysisID])
	}

	glog.V(3).Infof("readCollectedResult Was notified and found AnalysisList for %s", analysisID)
	result := an.outputResultCollection[analysisID]
	delete(an.outputResultCollection, analysisID)
	an.outputResultCondition.L.Unlock()
	return result
}

// collector is a go-func for
func (an *Analyzer) resultCollector() {
	res := make(map[string]AnalysisList)
	for {
		select {
		case <-an.collectorQuit:
			// Great! Timeout happened; stop collecting because we can.
			glog.V(2).Info("resultCollector: Received shutdown signal - returning")
			return
		case output := <-an.katagoOutput:
			an.processKatagoOutput(output, res)
		}
	}
}

// processKatagoOutput takes the raw string-output from katago and decides what
// to do with it.
func (an *Analyzer) processKatagoOutput(output string, resultStorage map[string]AnalysisList) {
	if strings.Contains(output, "error") {
		glog.Warningf("resultCollector: Katago had an error: %v", output)
		return
	}

	glog.V(3).Infof("Collected %s", output)
	res, err := ParseAnalysis(output)
	if err != nil {
		glog.Warningf("resultCollector: Error decoding analysis: %s", output)
	}

	// Add the current result to the AnalysisList
	al := resultStorage[res.ID]
	resultStorage[res.ID] = append(al, res)

	// If we've collected all the results for this analysis run, notify the caller.
	sizeLookup, ok := an.analysisSizeMap.Load(res.ID)

	if !ok {
		glog.Warningf("resultCollector: Expected to find number of results for analysis run, but did not. Run ID: %s", res.ID)
		return
	}

	expectedNumResults := sizeLookup.(int)
	actualNumResults := len(resultStorage[res.ID])
	glog.V(3).Infof("resultCollector: found %d/%d expected analysis results for %s", actualNumResults, expectedNumResults, res.ID)
	if actualNumResults >= expectedNumResults {
		glog.V(3).Infof("resultCollector: found sufficent results to return for %s", res.ID)
		resultSlice := resultStorage[res.ID]
		an.writeCollectedResult(res.ID, &resultSlice)
		delete(resultStorage, res.ID)
	}
}

// Stop the Katago Analyzer and the goroutines
func (an *Analyzer) Stop() error {
	glog.Infof("Shutting down Katago analyzer")
	an.stdinQuit <- 1
	an.collectorQuit <- 1
	return nil
}

// Cmd creates the Katago analysis command.
func (an *Analyzer) createCmd() *exec.Cmd {
	return exec.Command("katago", "analysis", "-model", an.model, "-config", an.config, "-analysis-threads", strconv.Itoa(an.analysisThreads))
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
	glog.V(3).Infof("katago stderr: started")
	scanner := bufio.NewScanner(stderr)
	var currentText string
	for scanner.Scan() {
		currentText = scanner.Text()

		if strings.Contains(currentText, katagoReadyStr) {
			glog.V(2).Infof("katago stderr: Katago ready-string found; ready to process input.")
			an.bootWait.Done()
		}
		glog.V(2).Infof("katago stderr: %v\n", currentText)
	}
	glog.V(2).Infof("katago stderr: returned")
}

func (an *Analyzer) stdOutReader(stdout io.ReadCloser) {
	defer stdout.Close()
	glog.V(2).Infof("katago stdout: started")
	scanner := bufio.NewScanner(stdout)
	var currentText string
	for scanner.Scan() {
		currentText = scanner.Text()
		glog.V(3).Infof("katago stdout: Sending result back from reader: %v\n", currentText)
		an.katagoOutput <- currentText
	}
	glog.V(2).Infof("katago stdout: returned")
}

func (an *Analyzer) stdInWriter(stdin io.WriteCloser) {
	defer stdin.Close()
	glog.V(2).Infof("katago stdin: started")
	for {
		select {
		case <-an.stdinQuit:
			// Make the go-routine shut-down
			// This will also close stdin, which will shut down Katago.
			glog.V(2).Infof("katago stdin: returning")
			return
		case writeValue := <-an.stdinWrite:
			glog.V(2).Info("katago stdin: Got value to write to Katago!")
			stdin.Write([]byte(writeValue))
		}
	}
}
