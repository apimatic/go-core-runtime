package utilities

import (
	"encoding/json"
	"errors"
	"strings"
)

type Bike struct {
	Id   int     `json:"id"`
	Roof *string `json:"roof"`
	Type *string `json:"type"`
}

func (b *Bike) UnmarshalJSON(input []byte) error {
	var temp bike
	err := json.Unmarshal(input, &temp)
	if err != nil {
		return err
	}
	err = temp.validate(input)
	if err != nil {
		return err
	}
	b.Id = *temp.Id
	b.Roof = temp.Roof
	b.Type = temp.Type
	return nil
}

type bike struct {
	Id   *int    `json:"id"`
	Roof *string `json:"roof"`
	Type *string `json:"type"`
}

func (b *bike) validate(input []byte) error {
	var errs []string
	if b.Id == nil {
		errs = append(errs, "required field `Id` is missing")
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.New(strings.Join(errs, "\n\t=> "))
}
