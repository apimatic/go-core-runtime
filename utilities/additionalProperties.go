package utilities

import (
	"encoding/json"
	"errors"
)

// MapAdditionalProperties is deprecated. Use MergeAdditionalProperties instead.
// This function is maintained for backward compatibility, appending additional properties
// to the destination struct map.
var MapAdditionalProperties = MergeAdditionalProperties[any]

// UnmarshalAdditionalProperties is deprecated. Use ExtractAdditionalProperties instead.
// This function is maintained for backward compatibility, unmarshal additional properties
// from the input and removing fields that exist on the parent struct.
var UnmarshalAdditionalProperties = ExtractAdditionalProperties[any]

// DetectConflictingProperties checks if any of the keys in structProperties exist in the dstMap.
// If a key is found, it returns an error indicating a conflict with one of the model's properties.
func DetectConflictingProperties[T any](dstMap map[string]T, structProperties ...string) error {
	for _, key := range structProperties {
		if _, ok := dstMap[key]; ok {
			return errors.New("an additional property key, '" + key + "' conflicts with one of the model's properties")
		}
	}
	return nil
}

// MergeAdditionalProperties merges additional properties from the source map
// into the destination map. If a key exists in both, the source map value overwrites
// the destination map value.
func MergeAdditionalProperties[T any](destinationMap additionalProperties[any], sourceMap additionalProperties[T]) {
	for key, value := range sourceMap {
		destinationMap[key] = value
	}
}

// ExtractAdditionalProperties unmarshal additional properties from the input and removes
// fields that exist on the parent struct based on the provided keys to remove.
func ExtractAdditionalProperties[T any](input []byte, keysToRemove ...string) (map[string]T, error) {
	rawMap, err := unmarshalAndFilterProperties(input, keysToRemove)

	destinationMap := additionalProperties[T]{}
	// Unmarshal the remaining properties into the destinationMap.
	for key, value := range rawMap {
		var typedVal T
		if err := json.Unmarshal(value, &typedVal); err == nil {
			destinationMap[key] = typedVal
		}
	}
	return destinationMap, err
}

// additionalProperties is a generic helper struct for handling additional properties in models.
// It allows for the storage of key-value pairs where keys are strings and values are of a specified type T.
type additionalProperties[T any] map[string]T

// unmarshalAndFilterProperties unmarshal additional properties from the input JSON byte array,
// removing any keys specified in keysToRemove.
func unmarshalAndFilterProperties(input []byte, keysToRemove []string) (map[string]json.RawMessage, error) {
	// Create a temporary map to hold the raw JSON data.
	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(input, &rawMap); err != nil {
		return rawMap, err
	}

	// Remove specified keys from the temporary map.
	for _, key := range keysToRemove {
		delete(rawMap, key)
	}
	return rawMap, nil
}
