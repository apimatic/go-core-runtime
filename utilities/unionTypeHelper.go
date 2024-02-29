package utilities

import (
	"encoding/json"
	"errors"
)

type TypeHolder struct {
	Value         any
	Flag          *bool
	Discriminator string
}

func (t *TypeHolder) getValue() any {
	*t.Flag = true
	return t.Value
}

func (t *TypeHolder) tryUnmarshall(data []byte) bool {
	err := json.Unmarshal(data, t.Value)
	return err == nil
}

// UnmarshallAnyOf tries to unmarshal the data into each of the provided types as an AnyOf group
// and return the converted value
func UnmarshallAnyOf(data []byte, types []*TypeHolder) (any, error) {
	return unmarshallUnionType(data, types, "", false)
}

// UnmarshallAnyOfWithDiscriminator tries to unmarshal the data into each of the provided types
// as an AnyOf group with discriminators and return the converted value
func UnmarshallAnyOfWithDiscriminator(data []byte, types []*TypeHolder, discField string) (any, error) {
	return unmarshallUnionType(data, types, discField, false)
}

// UnmarshallOneOf tries to unmarshal the data into each of the provided types as a OneOf group
// and return the converted value
func UnmarshallOneOf(data []byte, types []*TypeHolder) (any, error) {
	return unmarshallUnionType(data, types, "", true)
}

// UnmarshallOneOfWithDiscriminator tries to unmarshal the data into each of the provided types
// as a OneOf group with discriminators and return the converted value
func UnmarshallOneOfWithDiscriminator(data []byte, types []*TypeHolder, discField string) (any, error) {
	return unmarshallUnionType(data, types, discField, true)
}

// unmarshallUnionType tries to unmarshal the byte array into each of the provided types
// and return the converted value
func unmarshallUnionType(data []byte, types []*TypeHolder, discField string, isOneOf bool) (any, error) {
	if t := selectDiscriminatedTypeHolder(data, types, discField); t != nil {
		if t.tryUnmarshall(data) {
			return t.getValue(), nil
		}
		return nil, errors.New("failed to unmarshal into the selected discriminated type")
	}

	var selected *TypeHolder
	for _, t := range types {
		if t.tryUnmarshall(data) {
			if !isOneOf {
				return t.getValue(), nil
			} else if selected != nil {
				return nil, errors.New("can not map more then one type")
			}
			selected = t
		}
	}
	if isOneOf && selected != nil {
		return selected.getValue(), nil
	}
	return nil, errors.New("failed to unmarshal into any of the provided types")
}

// Select a specific typeHolder from given types based on available discriminator Field in the data
func selectDiscriminatedTypeHolder(data []byte, types []*TypeHolder, discField string) *TypeHolder {
	if discField == "" {
		return nil
	}
	dict := map[string]any{}
	err := json.Unmarshal(data, &dict)

	if err != nil {
		return nil
	}
	discValue, hasValue := dict[discField]

	if !hasValue {
		return nil
	}
	for _, t := range types {
		if t.Discriminator == discValue {
			return t
		}
	}
	return nil
}
