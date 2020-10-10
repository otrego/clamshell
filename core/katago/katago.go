// Package katago provides wrappers for analyzing games with katago.
package katago

import (
	"bufio"
	"errors"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
)

const katagoReadyStr = "Started, ready to begin handling requests"

// Analyzer is a katago-analyzer.
type Analyzer struct {
	// Model is the path to the model file.
	Model string

	// Config is the path to the config file.
	Config string

	// Number of analysis threads. Defaults to 16 if not specified.
	AnalysisThreads int

	cmd          *exec.Cmd
	stdoutQuit   chan int
	stderrQuit   chan int
	stdinQuit    chan int
	stdinWrite   chan string
	katagoOutput chan string

	// Wait group to wait until Katago has booted and proces the katagoReadyStr
	bootWait sync.WaitGroup
}

var analyzerSingleton *Analyzer
var once sync.Once

// GetAnalyzer returns the Katago Analyzer Singleton
func GetAnalyzer(model string, configPath string, numThreads int) (*Analyzer, error) {
	var err error
	once.Do(func() {
		threads := numThreads
		if numThreads == 0 {
			threads = 8
		}
		analyzerSingleton = &Analyzer{
			Model:           model,
			Config:          configPath,
			AnalysisThreads: threads,
			stdoutQuit:      make(chan int),
			stderrQuit:      make(chan int),
			stdinQuit:       make(chan int),
			stdinWrite:      make(chan string),
			katagoOutput:    make(chan string),
		}
		analyzerSingleton.createCmd()
		glog.Infof("using model %q", analyzerSingleton.Model)
		glog.Infof("using gtp config %q", analyzerSingleton.Config)

		analyzerSingleton.bootWait.Add(1)
		err = analyzerSingleton.start()
	})
	return analyzerSingleton, err
}

// AnalyzeGame is a synchronous request for Katago to process a game
func (an *Analyzer) AnalyzeGame(json string) (*string, error) {
	an.bootWait.Wait()
	an.stdinWrite <- json + "\n"
	output := <-an.katagoOutput
	if strings.Contains(output, "error") {
		glog.Warningf("Katago had an error: %v", output)
		return nil, errors.New(output)
	}
	return &output, nil
}

// Start the Katago Analyzer and make it ready to receive input
func (an *Analyzer) start() error {
	glog.Info("Starting katago analyzer")

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
}

func (an *Analyzer) stdInWriter(stdin io.WriteCloser) {
	defer stdin.Close()
	glog.Infof("katago stdin: started")
	for {
		select {
		case <-an.stdinQuit:
			// Make the go-routine shut-down-able.
			return
		case writeValue := <-an.stdinWrite:
			glog.Info("katago stdin: Got value to write to Katago!")
			stdin.Write([]byte(writeValue))
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
