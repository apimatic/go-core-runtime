package utilities

import (
	"encoding/json"
	"errors"
	"strings"
)

type Atom struct {
	NumberOfElectrons int `json:"number_of_electrons"`
	NumberOfProtons   int `json:"number_of_protons"`
}

func (c *Atom) UnmarshalJSON(input []byte) error {
	var temp atom
	err := json.Unmarshal(input, &temp)
	if err != nil {
		return NewMarshalError("Atom", err)
	}
	err = temp.validate(input)
	if err != nil {
		return err
	}
	c.NumberOfElectrons = *temp.NumberOfElectrons
	c.NumberOfProtons = *temp.NumberOfProtons
	return nil
}

type atom struct {
	NumberOfElectrons *int `json:"number_of_electrons"`
	NumberOfProtons   *int `json:"number_of_protons"`
}

func (a *atom) validate(input []byte) error {
	var errs []string
	if a.NumberOfElectrons == nil {
		errs = append(errs, "required field `NumberOfElectrons` is missing")
	}
	if a.NumberOfProtons == nil {
		errs = append(errs, "required field `NumberOfProtons` is missing")
	}
	if len(errs) == 0 {
		return nil
	}
	return NewMarshalError("Atom", errors.New(strings.Join(errs, "\n\t=> ")))
}
