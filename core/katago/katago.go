// Package katago provides wrappers for analyzing games with katago.
package katago

import (
	"io"
	"os/exec"
	"strconv"

	"github.com/golang/glog"
)

// Analyzer is a katago-analyzer.
type Analyzer struct {
	// Model is the path to the model file.
	Model string

	// Config is the path to the config file.
	Config string

	// Number of analysis threads. Defaults to 16 if not specified.
	AnalysisThreads int
}

// Cmd creates the Katago analysis command.
func (an *Analyzer) Cmd() *exec.Cmd {
	threads := an.AnalysisThreads
	if threads == 0 {
		threads = 8
	}
	return exec.Command("katago", "analysis", "-model", an.Model, "-config", an.Config, "-analysis-threads", strconv.Itoa(threads))
}

// Start the katago analyzer.
func (an *Analyzer) Start() (io.WriteCloser, io.ReadCloser, error) {
	glog.Info("Starting katago analyzer")

	cmd := an.Cmd()
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, err
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, nil, err
	}
	return stdinPipe, stdoutPipe, nil
}
