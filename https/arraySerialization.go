package https

import "fmt"

// ArraySerializationOption represents the type for request array serialization options.
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

var ArraySerializationOptionStrings = [...]string{
	"Indexed",
	"UnIndexed",
	"Plain",
	"Csv",
	"Tsv",
	"Psv",
}

func (option ArraySerializationOption) getSeparator() rune {
	switch option {
	case Csv:
		return ','
	case Tsv:
		return '\t'
	case Psv:
		return '|'
	default:
		return -1
	}
}

func (option ArraySerializationOption) joinKey(keyPrefix string, index any) string {
	if index == nil {
		switch option {
		case UnIndexed:
			return fmt.Sprintf("%v[]", keyPrefix)
		default:
			return fmt.Sprintf("%v", keyPrefix)
		}
	}
	indexedKey := fmt.Sprintf("%v", index)
	return fmt.Sprintf("%v[%v]", keyPrefix, indexedKey)
}

func (option ArraySerializationOption) appendMap(result map[string][]string, param map[string][]string) {
	for k, values := range param {
		for _, value := range values {
			option.append(result, k, value)
		}
	}
}

func (option ArraySerializationOption) append(result map[string][]string, key string, value string) {
	separator := option.getSeparator()
	if len(result[key]) > 0 && separator != -1 {
		result[key][0] = fmt.Sprintf("%v%c%v", result[key][0], separator, value)
	} else {
		result[key] = append(result[key], value)
	}
}
