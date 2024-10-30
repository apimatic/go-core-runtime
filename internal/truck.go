package internal

import (
	"encoding/json"
	"errors"
	"github.com/apimatic/go-core-runtime/utilities"
	"strings"
)

type Truck struct {
	Id     int     `json:"id"`
	Weight string  `json:"weight"`
	Roof   *string `json:"roof"`
}

func (c *Truck) UnmarshalJSON(input []byte) error {
	var temp truck
	err := json.Unmarshal(input, &temp)
	if err != nil {
		return utilities.NewMarshalError("Truck", err)
	}
	err = temp.validate(input)
	if err != nil {
		return err
	}
	c.Id = *temp.Id
	c.Weight = *temp.Weight
	c.Roof = temp.Roof
	return nil
}

type truck struct {
	Id     *int    `json:"id"`
	Weight *string `json:"weight"`
	Roof   *string `json:"roof"`
}

func (t *truck) validate(input []byte) error {
	var errs []string
	if t.Id == nil {
		errs = append(errs, "required field `Id` is missing")
	}
	if t.Weight == nil {
		errs = append(errs, "required field `Weight` is missing")
	}
	if len(errs) == 0 {
		return nil
	}
	return utilities.NewMarshalError("Truck", errors.New(strings.Join(errs, "\n\t=> ")))
}
