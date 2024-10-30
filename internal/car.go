package internal

import (
	"encoding/json"
	"errors"
	"github.com/apimatic/go-core-runtime/utilities"
	"strings"
)

type Car struct {
	Id   int     `json:"id"`
	Roof *string `json:"roof"`
	Type *string `json:"type"`
}

func (c *Car) UnmarshalJSON(input []byte) error {
	var temp car
	err := json.Unmarshal(input, &temp)
	if err != nil {
		return utilities.NewMarshalError("Car", err)
	}
	err = temp.validate(input)
	if err != nil {
		return err
	}
	c.Id = *temp.Id
	c.Roof = temp.Roof
	c.Type = temp.Type
	return nil
}

type car struct {
	Id   *int    `json:"id"`
	Roof *string `json:"roof"`
	Type *string `json:"type"`
}

func (c *car) validate(input []byte) error {
	var errs []string
	if c.Id == nil {
		errs = append(errs, "required field `Id` is missing")
	}
	if len(errs) == 0 {
		return nil
	}
	return utilities.NewMarshalError("Car", errors.New(strings.Join(errs, "\n\t=> ")))
}
