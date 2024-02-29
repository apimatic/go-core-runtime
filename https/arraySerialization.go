package https

import "fmt"

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

func (option ArraySerializationOption) appendValue(result map[string][]string, key string, value string) {
	separator := option.getSeparator()
	if len(result[key]) > 0 && separator != -1 {
		result[key][0] = fmt.Sprintf("%v%c%v", result[key][0], separator, value)
	} else {
		result[key] = append(result[key], value)
	}
}

func (option ArraySerializationOption) appendMap(result map[string][]string, param map[string][]string) {
	for k, v := range param {
		for _, v1 := range v {
			option.appendValue(result, k, v1)
		}
	}
}

func (option ArraySerializationOption) prePareKey(keyPrefix string, obj any) string {

	if (obj == nil) {
		switch option {
		case UnIndexed:
			return fmt.Sprintf("%v[]", keyPrefix)
		case Plain:
			return fmt.Sprintf("%v", keyPrefix)
		}
		return fmt.Sprintf("%v", keyPrefix)
	}
	indexedKey := fmt.Sprintf("%v", obj)
	return fmt.Sprintf("%v[%v]", keyPrefix, indexedKey)
}
