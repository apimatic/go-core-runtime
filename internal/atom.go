package internal

import (
	"encoding/json"
	"errors"
	"github.com/apimatic/go-core-runtime/utilities"
	"strings"
)

type Atom struct {
	NumberOfElectrons int `json:"number_of_electrons"`
	NumberOfProtons   int `json:"number_of_protons"`
}

func (a *Atom) UnmarshalJSON(input []byte) error {
	var temp atom
	err := json.Unmarshal(input, &temp)
	if err != nil {
		return utilities.NewMarshalError("Atom", err)
	}
	err = temp.validate()
	if err != nil {
		return err
	}
	a.NumberOfElectrons = *temp.NumberOfElectrons
	a.NumberOfProtons = *temp.NumberOfProtons
	return nil
}

type atom struct {
	NumberOfElectrons *int `json:"number_of_electrons"`
	NumberOfProtons   *int `json:"number_of_protons"`
}

func (a *atom) validate() error {
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
	return utilities.NewMarshalError("Atom", errors.New(strings.Join(errs, "\n\t=> ")))
}
