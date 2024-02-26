package https

// RequestRetryOption represents the type for request retry options.
type ArraySerializationOption int

// Constants for different request retry options.
const (
	Indexed = iota
	UnIndexed
	Plain
	Csv
	Tsv
	Psv
)
