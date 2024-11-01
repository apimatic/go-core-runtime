package utilities

import (
	"encoding/json"
	"errors"
	"strings"
)

func ValidateAdditionalProperty[T any](dstMap map[string]T, keysToRemove ...string) error {
	containsKey := func(key string) bool {
		for _, tag := range keysToRemove {
			if tag == key {
				return true
			}
		}
		return false
	}

	for key := range dstMap {
		if strings.TrimSpace(key) == "" {
			return errors.New("an additional property key can not be empty or whitespace")
		}
		if containsKey(key) {
			return errors.New("an additional property key, '" + key + "' conflicts with one of the model's properties")
		}
	}
	return nil
}

// MapAdditionalProperties append additional properties to destination struct map
var MapAdditionalProperties = MapAdditionalProperty[any]

// MapAdditionalProperty append additional properties to destination struct map
func MapAdditionalProperty[T any](destinationMap additionalProperties[any], sourceMap additionalProperties[T]) {
	for key, value := range sourceMap {
		destinationMap[key] = value
	}
}

// UnmarshalAdditionalProperties unmarshal additional properties and remove fields that exists on parent struct
var UnmarshalAdditionalProperties = UnmarshalAdditionalProperty[any]

// UnmarshalAdditionalProperty unmarshal additional properties and remove fields that exists on parent struct
func UnmarshalAdditionalProperty[T any](input []byte, keysToRemove ...string) (map[string]T, error) {
	destinationMap := additionalProperties[T]{}
	err := destinationMap.unmarshalAdditionalProperties(input, keysToRemove)
	return destinationMap, err
}

// additionalProperties helper struct for handling additional properties in models
type additionalProperties[T any] map[string]T

func (srcMap *additionalProperties[T]) unmarshalAdditionalProperties(input []byte, keysToRemove []string) error {
	var dstRawMap map[string]json.RawMessage
	if err := json.Unmarshal(input, &dstRawMap); err != nil {
		return err
	}
	for _, key := range keysToRemove {
		delete(dstRawMap, key)
	}
	for key, value := range dstRawMap {
		var typedVal T
		if err := json.Unmarshal(value, &typedVal); err == nil {
			(*srcMap)[key] = typedVal
		}
	}
	return nil
}
