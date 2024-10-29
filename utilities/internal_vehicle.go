package utilities

import (
	"encoding/json"
	"errors"
	"strings"
)

func ToPointer[T any](value T) *T {
	return &value
}

type Vehicle[T any] struct {
	Year                 int          `json:"year"`
	Make                 *string      `json:"make"`
	Model                *string      `json:"model"`
	AdditionalProperties map[string]T `json:"_"`
}

func (v Vehicle[T]) MarshalJSON() (
	[]byte,
	error) {
	if err := ValidateAdditionalProperty(v.AdditionalProperties,
		"year", "make", "model"); err != nil {
		return nil, err
	}
	return json.Marshal(v.toMap())
}

func (v Vehicle[T]) toMap() map[string]any {
	structMap := make(map[string]any)
	MapAdditionalProperty(structMap, v.AdditionalProperties)
	if v.Make != nil {
		structMap["make"] = *v.Make
	} else {
		structMap["make"] = "ferrari"
	}
	if v.Model != nil {
		structMap["model"] = *v.Model
	} else {
		structMap["model"] = "MONZA SP2"
	}
	structMap["year"] = v.Year
	return structMap
}

func (v *Vehicle[T]) UnmarshalJSON(input []byte) error {
	var temp tempVehicle
	err := json.Unmarshal(input, &temp)
	if err != nil {
		return NewMarshalError("Vehicle", err)
	}
	err = temp.validate()
	if err != nil {
		return err
	}
	additionalProperties, err := UnmarshalAdditionalProperty[T](input, "year", "make", "model")
	if err != nil {
		return err
	}
	v.AdditionalProperties = additionalProperties
	v.Year = *temp.Year
	v.Make = temp.Make
	v.Model = temp.Model
	return nil
}

type tempVehicle struct {
	Year  *int    `json:"year"`
	Make  *string `json:"make"`
	Model *string `json:"model"`
}

func (c *tempVehicle) validate() error {
	var errs []string
	if c.Year == nil {
		errs = append(errs, "required field `Year` is missing")
	}
	if len(errs) == 0 {
		return nil
	}
	return NewMarshalError("Vehicle", errors.New(strings.Join(errs, "\n\t=> ")))
}
