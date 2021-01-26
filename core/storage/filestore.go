package storage

import "context"

// Filestore is a interface for file storage
type Filestore interface {
	Get(context.Context, StoredDataType, string) (string, error)
	List(context.Context, StoredDataType) ([]string, error)
	Put(context.Context, StoredDataType, string, string) error
}

// StoredDataType designates which type of data is stored
type StoredDataType string

const (
	// Games storage directory
	Games StoredDataType = "games"
	// Analysis is the data type for katago output
	Analysis StoredDataType = "analysis"
	// Problems is the directory for AI problems that have been generated
	// and turned into SGF.
	Problems StoredDataType = "gen_problems"
)

var storedDataTypes = []StoredDataType{
	Games,
	Analysis,
	Problems,
}
