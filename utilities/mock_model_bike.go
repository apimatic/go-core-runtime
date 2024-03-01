package utilities

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/apimatic/go-core-runtime/types"
)

type Bike struct {
	Id       int                  `json:"id"`
	Roof     *string              `json:"roof"`
	AirLevel types.Optional[Atom] `json:"air_level"`
	Type     *string              `json:"type"`
}

func (b *Bike) UnmarshalJSON(input []byte) error {
	var temp bike
	err := json.Unmarshal(input, &temp)
	if err != nil {
		return NewMarshallerError("Bike", err)
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
	Id       *int                 `json:"id"`
	Roof     *string              `json:"roof"`
	AirLevel types.Optional[Atom] `json:"air_level"`
	Type     *string              `json:"type"`
}

func (b *bike) validate(input []byte) error {
	var errs []string
	if b.Id == nil {
		errs = append(errs, "required field `Id` is missing")
	}
	if len(errs) == 0 {
		return nil
	}
	return NewMarshallerError("Bike", errors.New(strings.Join(errs, "\n\t=> ")))
}
