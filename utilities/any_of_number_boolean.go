package utilities

import (
	"encoding/json"
	"errors"
	"strings"
)

// AnyOfNumberVehicle represents a AnyOfNumberVehicle struct.
// This is a container for any-of cases.
type AnyOfNumberVehicle struct {
	value     any
	isNumber  bool
	isVehicle bool
}

// String converts the AnyOfNumberVehicle object to a string representation.
func (s AnyOfNumberVehicle) String() string {
	if bytes, err := json.Marshal(s.value); err == nil {
		return strings.Trim(string(bytes), "\"")
	}
	return ""
}

// MarshalJSON implements the json.Marshaler interface for AnyOfNumberVehicle.
// It customizes the JSON marshaling process for AnyOfNumberVehicle objects.
func (s AnyOfNumberVehicle) MarshalJSON() (
	[]byte,
	error) {
	if s.value == nil {
		return nil, errors.New("No underlying type is set. Please use any of the `models.AnyOfNumberBooleanContainer.From*` functions to initialize the AnyOfNumberVehicle object.")
	}
	return json.Marshal(s.toMap())
}

// toMap converts the AnyOfNumberVehicle object to a map representation for JSON marshaling.
func (s AnyOfNumberVehicle) toMap() any {
	switch obj := s.value.(type) {
	case *int:
		return *obj
	case *Vehicle[bool]:
		return obj.toMap()
	}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for AnyOfNumberVehicle.
// It customizes the JSON unmarshaling process for AnyOfNumberVehicle objects.
func (s *AnyOfNumberVehicle) UnmarshalJSON(input []byte) error {
	result, err := UnmarshallAnyOf(input,
		NewTypeHolder(new(int), false, &s.isNumber),
		NewTypeHolder(&Vehicle[bool]{}, false, &s.isVehicle),
	)

	s.value = result
	return err
}

func (s *AnyOfNumberVehicle) AsNumber() (
	*int,
	bool) {
	if !s.isNumber {
		return nil, false
	}
	return s.value.(*int), true
}

func (s *AnyOfNumberVehicle) AsVehicle() (
	*Vehicle[bool],
	bool) {
	if !s.isVehicle {
		return nil, false
	}
	return s.value.(*Vehicle[bool]), true
}

// internalAnyOfNumberBoolean represents a AnyOfNumberVehicle struct.
// This is a container for any-of cases.
type internalAnyOfNumberBoolean struct{}

var AnyOfNumberBooleanContainer internalAnyOfNumberBoolean

// The internalAnyOfNumberBoolean instance, wrapping the provided int value.
func (s *internalAnyOfNumberBoolean) FromNumber(val int) AnyOfNumberVehicle {
	return AnyOfNumberVehicle{value: &val}
}

// The internalAnyOfNumberBoolean instance, wrapping the provided bool value.
func (s *internalAnyOfNumberBoolean) FromVehicle(val Vehicle[bool]) AnyOfNumberVehicle {
	return AnyOfNumberVehicle{value: &val}
}
