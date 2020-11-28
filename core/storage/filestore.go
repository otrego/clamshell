package storage

// Filestore is a interface for file storage
type Filestore interface {
	Get(StoredDataType, string) (string, error)
	List(StoredDataType, string) ([]string, error)
	Put(StoredDataType, string, string) error
}

// StoredDataType designates which type of data is stored
type StoredDataType string

const (
	// SGFInput storage directory
	SGFInput StoredDataType = "input_sgf"
	// AnalysisOutput is the data type for katago output
	AnalysisOutput StoredDataType = "katago_analysis"
	// AiProblemIdentification is the directory for AI identifed problems
	AiProblemIdentification StoredDataType = "ai_problem_identification"
	// AiProblemSGF is the directory for AI problems that have been generated
	// and turned into SGF.
	AiProblemSGF StoredDataType = "ai_problem_sgf"
)

var storedDataTypes = []StoredDataType{
	SGFInput,
	AnalysisOutput,
	AiProblemIdentification,
	AiProblemSGF,
}
