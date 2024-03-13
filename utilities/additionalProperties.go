package utilities

import "encoding/json"

func MapAdditionalProperties(destinationMap additionalProperties, sourceMap additionalProperties) {
	destinationMap.appendMap(sourceMap)
}

func UnmarshalAdditionalProperties(input []byte, keysToRemove ...string) (map[string]any, error) {
	var destinationMap additionalProperties
	err := destinationMap.unmarshalAdditionalProperties(input, keysToRemove)
	return destinationMap, err
}

type additionalProperties map[string]any

func (dstMap *additionalProperties) appendMap(srcMap additionalProperties) {
	for key, value := range srcMap {
		(*dstMap)[key] = value
	}
}

func (dstMap *additionalProperties) unmarshalAdditionalProperties(input []byte, keysToRemove []string) error {
	if err := json.Unmarshal(input, &dstMap); err != nil {
		return err
	}
	for _, key := range keysToRemove {
		delete(*dstMap, key)
	}
	return nil
}
