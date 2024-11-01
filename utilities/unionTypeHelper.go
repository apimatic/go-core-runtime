package utilities

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type TypeHolder struct {
	value          any
	originalValue  any
	isNullableType bool
	isSelected     *bool
	discriminator  string
	typeError      error
}

func NewTypeHolder(val any, isNullableType bool, isSelected *bool) *TypeHolder {
	return &TypeHolder{
		value:          val,
		originalValue:  val,
		isNullableType: isNullableType,
		isSelected:     isSelected,
	}
}

func NewTypeHolderDiscriminator(val any, isNullableType bool, isSelected *bool, discriminator string) *TypeHolder {
	return &TypeHolder{
		value:          val,
		originalValue:  val,
		isNullableType: isNullableType,
		isSelected:     isSelected,
		discriminator:  discriminator,
	}
}

func (t *TypeHolder) selectValue() any {
	*t.isSelected = true
	return t.value
}

func (t *TypeHolder) tryUnmarshall(data []byte) bool {
	if string(data) == `null` {
		if t.isNullableType {
			t.value = nil
		} else {
			typeName := reflect.TypeOf(t.originalValue).String()
			t.typeError = errors.New("json: cannot unmarshal null into Go value of type " + typeName)
		}
		return t.isNullableType
	}
	err := json.Unmarshal(data, t.value)
	t.typeError = err
	return err == nil
}

// UnmarshallAnyOf tries to unmarshal the data into each of the provided types as an AnyOf group
// and return the converted value
func UnmarshallAnyOf(data []byte, types ...*TypeHolder) (any, error) {
	return unmarshallUnionType(data, types, false)
}

// UnmarshallAnyOfWithDiscriminator tries to unmarshal the data into each of the provided types
// as an AnyOf group with discriminators and return the converted value
func UnmarshallAnyOfWithDiscriminator(data []byte, discField string, types ...*TypeHolder) (any, error) {
	return unmarshallUnionType(data, filterTypeHolders(data, types, discField), false)
}

// UnmarshallOneOf tries to unmarshal the data into each of the provided types as a OneOf group
// and return the converted value
func UnmarshallOneOf(data []byte, types ...*TypeHolder) (any, error) {
	return unmarshallUnionType(data, types, true)
}

// UnmarshallOneOfWithDiscriminator tries to unmarshal the data into each of the provided types
// as a OneOf group with discriminators and return the converted value
func UnmarshallOneOfWithDiscriminator(data []byte, discField string, types ...*TypeHolder) (any, error) {
	return unmarshallUnionType(data, filterTypeHolders(data, types, discField), true)
}

// filterTypeHolders filter out the typeHolders from given list based on
// available discriminator field's value in the data
func filterTypeHolders(data []byte, types []*TypeHolder, discField string) []*TypeHolder {
	discValue, ok := extractDiscriminatorValue(data, discField)
	if !ok {
		return types
	}
	for _, t := range types {
		if t.discriminator != "" && t.discriminator == discValue {
			return []*TypeHolder{t}
		}
	}
	return types
}

// extractDiscriminatorValue extracts the discriminator value using the discriminator field
func extractDiscriminatorValue(data []byte, discField string) (any, bool) {
	if discField == "" {
		return nil, false
	}
	dict := map[string]any{}
	err := json.Unmarshal(data, &dict)

	if err != nil {
		return nil, false
	}
	discValue, ok := dict[discField]

	return discValue, ok
}

// unmarshallUnionType tries to unmarshal the byte array into each of the provided types
// and return the converted value
func unmarshallUnionType(data []byte, types []*TypeHolder, matchExactlyOneType bool) (any, error) {
	var selected *TypeHolder
	for _, t := range types {
		if t.tryUnmarshall(data) {
			if !matchExactlyOneType {
				return t.selectValue(), nil
			} else if selected != nil {
				return nil, moreThenOneTypeMatchesError(selected, t, data)
			}
			selected = t
		}
	}
	if matchExactlyOneType && selected != nil {
		return selected.selectValue(), nil
	}
	return nil, noneTypeMatchesError(types, data)
}

func moreThenOneTypeMatchesError(type1 *TypeHolder, type2 *TypeHolder, data []byte) error {
	type1Name := reflect.TypeOf(type1.originalValue).String()
	type2Name := reflect.TypeOf(type2.originalValue).String()
	return errors.New("There are more than one matching types i.e. {" + type1Name + " and " + type2Name + "} on: " + string(data))
}

func noneTypeMatchesError(types []*TypeHolder, data []byte) error {
	names := make([]string, len(types))
	reasons := make([]string, len(types))

	for i, t := range types {
		names[i] = reflect.TypeOf(t.originalValue).String()
		reasons[i] = "\n\nError " + fmt.Sprint(i+1) + ":\n  => " + t.typeError.Error()
	}

	return errors.New("We could not match any acceptable type from {" + strings.Join(names, ", ") + "} on: " + string(data) + strings.Join(reasons, ""))
}
