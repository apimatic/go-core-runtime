package utilities

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type TypeHolder struct {
	Value         any
	Flag          *bool
	Discriminator string
	Error         error
}

func (t *TypeHolder) selectValue() any {
	*t.Flag = true
	return t.Value
}

func (t *TypeHolder) tryUnmarshall(data []byte) bool {
	err := json.Unmarshal(data, t.Value)
	t.Error = err
	return err == nil
}

// UnmarshallAnyOf tries to unmarshal the data into each of the provided types as an AnyOf group
// and return the converted value
func UnmarshallAnyOf(data []byte, types []*TypeHolder) (any, error) {
	return unmarshallUnionType(data, types, false)
}

// UnmarshallAnyOfWithDiscriminator tries to unmarshal the data into each of the provided types
// as an AnyOf group with discriminators and return the converted value
func UnmarshallAnyOfWithDiscriminator(data []byte, types []*TypeHolder, discField string) (any, error) {
	return unmarshallUnionType(data, filterTypeHolders(data, types, discField), false)
}

// UnmarshallOneOf tries to unmarshal the data into each of the provided types as a OneOf group
// and return the converted value
func UnmarshallOneOf(data []byte, types []*TypeHolder) (any, error) {
	return unmarshallUnionType(data, types, true)
}

// UnmarshallOneOfWithDiscriminator tries to unmarshal the data into each of the provided types
// as a OneOf group with discriminators and return the converted value
func UnmarshallOneOfWithDiscriminator(data []byte, types []*TypeHolder, discField string) (any, error) {
	return unmarshallUnionType(data, filterTypeHolders(data, types, discField), true)
}

// filterTypeHolders filter out the typeholders from given list based on
// available discriminator field's value in the data
func filterTypeHolders(data []byte, types []*TypeHolder, discField string) []*TypeHolder {
	if discField == "" {
		return types
	}
	dict := map[string]any{}
	err := json.Unmarshal(data, &dict)

	if err != nil {
		return types
	}
	discValue, ok := dict[discField]

	if !ok {
		return types
	}
	for _, t := range types {
		if t.Discriminator != "" && t.Discriminator == discValue {
			return []*TypeHolder{t}
		}
	}
	return types
}

// unmarshallUnionType tries to unmarshal the byte array into each of the provided types
// and return the converted value
func unmarshallUnionType(data []byte, types []*TypeHolder, isOneOf bool) (any, error) {
	var selected *TypeHolder
	for _, t := range types {
		if t.tryUnmarshall(data) {
			if !isOneOf {
				return t.selectValue(), nil
			} else if selected != nil {
				return nil, moreThenOneTypeMatchesError(selected, t, data)
			}
			selected = t
		}
	}
	if isOneOf && selected != nil {
		return selected.selectValue(), nil
	}
	return nil, noneTypeMatchesError(types, data)
}

func moreThenOneTypeMatchesError(type1 *TypeHolder, type2 *TypeHolder, data []byte) error {
	type1Name := reflect.TypeOf(type1.Value).String()
	type2Name := reflect.TypeOf(type2.Value).String()
	return errors.New("There are more than one matching types i.e. {" + type1Name + " and " + type2Name + "} on: " + string(data))
}

func noneTypeMatchesError(types []*TypeHolder, data []byte) error {
	names := make([]string, len(types))
	reasons := make([]string, len(types))

	for i, t := range types {
		names[i] = reflect.TypeOf(t.Value).String()
		reasons[i] = "\n\nError " + fmt.Sprint(i+1) + ":\n  => " + t.Error.Error()
	}

	return errors.New("We could not match any acceptable type from {" + strings.Join(names, ", ") + "} on: " + string(data) + strings.Join(reasons, ""))
}
