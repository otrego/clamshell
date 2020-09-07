// Package katago provides wrappers for analyzing games with katago.
package katago

// Analyzer is a katago-analyzer
type Analyzer struct {
	// Model is the path to the model file.
	Model string

	// Config is the path to the config file.
	Config string
}
