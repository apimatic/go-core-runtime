package utilities

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UnionTypeCase struct {
	name               string
	types              []any
	discriminators     []string
	discriminatorField string
	testValue          string
	expectedValue      string
	expectedType       any
}

func (u *UnionTypeCase) Assert(t *testing.T, result any, err error, isSuccess bool) {
	assert.Nil(t, err)
	assert.True(t, isSuccess)
	assert.IsType(t, u.expectedType, result)
	marshalled, _ := json.Marshal(result)
	if u.expectedValue != "" {
		u.testValue = u.expectedValue
	}
	assert.Equal(t, u.testValue, string(marshalled))
}

func TestOneOfSuccess(t *testing.T) {
	var tests = []UnionTypeCase{
		{
			name:         `(string,int) => string1`,
			types:        []any{new(string), new(int)},
			testValue:    `"some string"`,
			expectedType: new(string),
		},
		{
			name:         `(string,int) => string2`,
			types:        []any{new(string), new(int)},
			testValue:    `"123"`,
			expectedType: new(string),
		},
		{
			name:         `(string,int) => int`,
			types:        []any{new(string), new(int)},
			testValue:    `123`,
			expectedType: new(int),
		},
		{
			name:         `(bool,int) => bool`,
			types:        []any{new(bool), new(int)},
			testValue:    `true`,
			expectedType: new(bool),
		},
		{
			name:         `(bool,int) => int`,
			types:        []any{new(bool), new(int)},
			testValue:    `0`,
			expectedType: new(int),
		},
		{
			name:         `(float,int) => float`,
			types:        []any{new(float32), new(int)},
			testValue:    `2345.123`,
			expectedType: new(float32),
		},
		{
			name:         `(float,int[]) => float`,
			types:        []any{new(float32), new([]int)},
			testValue:    `2345`,
			expectedType: new(float32),
		},
		{
			name:         `(int,int[]) => int[]`,
			types:        []any{new(int), new([]int)},
			testValue:    `[2345,1234]`,
			expectedType: new([]int),
		},
		{
			name:         `(int{},int[]) => int{}`,
			types:        []any{new(map[string]int), new([]int)},
			testValue:    `{"keyA":2345,"keyB":1234}`,
			expectedType: new(map[string]int),
		},
		{
			name:         `(int[],Atom) => Atom`,
			types:        []any{new([]int), &Atom{}},
			testValue:    `{"number_of_electrons":2345,"number_of_protons":1234}`,
			expectedType: &Atom{},
		},
		{
			name:         `(int{},Atom) => int{}`,
			types:        []any{new(map[string]int), &Atom{}},
			testValue:    `{"number_of_":2345,"number_of_protons":1234}`,
			expectedType: new(map[string]int),
		},
		{
			name:         `(Truck,Car) => Car`,
			types:        []any{&Truck{}, &Car{}},
			testValue:    `{"id":2345,"roof":"BIG","type":null}`,
			expectedType: &Car{},
		},
	}

	assertSuccessCases(t, tests, UnmarshallOneOf)
}

func TestOneOfDiscriminatorSuccess(t *testing.T) {
	var tests = []UnionTypeCase{
		{
			name:               `(Car,Bike) => Car`,
			types:              []any{&Car{}, &Bike{}},
			discriminators:     []string{"4 wheeler", "2 wheeler"},
			discriminatorField: "type",
			testValue:          `{"id":2345,"roof":"BIG","type":"4 wheeler"}`,
			expectedType:       &Car{},
		},
		{
			name:               `(Car,Bike) => Bike`,
			types:              []any{&Car{}, &Bike{}},
			discriminators:     []string{"4 wheeler", "2 wheeler"},
			discriminatorField: "type",
			testValue:          `{"id":2345,"roof":"BIG","type":"2 wheeler"}`,
			expectedType:       &Bike{},
		},
		{
			name:               `(string,bool) => bool`,
			types:              []any{new(string), new(bool)},
			discriminators:     []string{"my str", "my bool"},
			discriminatorField: "type",
			testValue:          `false`,
			expectedType:       new(bool),
		},
	}

	assertSuccessDiscriminatorCases(t, tests, UnmarshallOneOfWithDiscriminator)
}

func TestAnyOfSuccess(t *testing.T) {
	var tests = []UnionTypeCase{
		{
			name:         `(string,int) => string1`,
			types:        []any{new(string), new(int)},
			testValue:    `"some string"`,
			expectedType: new(string),
		},
		{
			name:         `(string,int) => string2`,
			types:        []any{new(string), new(int)},
			testValue:    `"123"`,
			expectedType: new(string),
		},
		{
			name:         `(string,int) => int`,
			types:        []any{new(string), new(int)},
			testValue:    `123`,
			expectedType: new(int),
		},
		{
			name:         `(bool,int) => bool`,
			types:        []any{new(bool), new(int)},
			testValue:    `true`,
			expectedType: new(bool),
		},
		{
			name:         `(bool,int) => int`,
			types:        []any{new(bool), new(int)},
			testValue:    `2345`,
			expectedType: new(int),
		},
		{
			name:         `(float,int) => float`,
			types:        []any{new(float32), new(int)},
			testValue:    `2345`,
			expectedType: new(float32),
		},
		{
			name:         `(int,float) => float`,
			types:        []any{new(int), new(float32)},
			testValue:    `2345.123`,
			expectedType: new(float32),
		},
		{
			name:         `(Truck,Car) => Truck`,
			types:        []any{&Truck{}, &Car{}},
			testValue:    `{"id":2345,"weight":"heavy","roof":"BIG"}`,
			expectedType: &Truck{},
		},
		{
			name:          `(Car,Truck) => Car`,
			types:         []any{&Car{}, &Truck{}},
			testValue:     `{"id":2345,"weight":"heavy","roof":"BIG"}`,
			expectedValue: `{"id":2345,"roof":"BIG","type":null}`,
			expectedType:  &Car{},
		},
	}

	assertSuccessCases(t, tests, UnmarshallAnyOf)
}

func TestAnyOfDiscriminatorSuccess(t *testing.T) {
	var tests = []UnionTypeCase{
		{
			name:               `(Car,Bike) => Car`,
			types:              []any{&Car{}, &Bike{}},
			discriminators:     []string{"4 wheeler", "2 wheeler"},
			discriminatorField: "type",
			testValue:          `{"id":2345,"roof":"BIG","type":"unknown"}`,
			expectedType:       &Car{},
		},
		{
			name:               `(Car,Bike) => Car2`,
			types:              []any{&Car{}, &Bike{}},
			discriminators:     []string{"4 wheeler", "2 wheeler"},
			discriminatorField: "type",
			testValue:          `{"id":2345,"roof":"BIG"}`,
			expectedValue:      `{"id":2345,"roof":"BIG","type":null}`,
			expectedType:       &Car{},
		},
	}

	assertSuccessDiscriminatorCases(t, tests, UnmarshallAnyOfWithDiscriminator)
}

func assertSuccessCases(t *testing.T, tests []UnionTypeCase, caller func([]byte, []*TypeHolder) (any, error)) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var isSuccess bool
			var typeHolders []*TypeHolder
			for _, tt := range test.types {
				typeHolders = append(typeHolders, &TypeHolder{Value: tt, Flag: &isSuccess})
			}
			result, err := caller([]byte(test.testValue), typeHolders)
			test.Assert(t, result, err, isSuccess)
		})
	}
}

func assertSuccessDiscriminatorCases(t *testing.T, tests []UnionTypeCase, caller func([]byte, []*TypeHolder, string) (any, error)) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var isSuccess bool
			var typeHolders []*TypeHolder
			for i, tt := range test.types {
				typeHolders = append(typeHolders, &TypeHolder{Value: tt, Flag: &isSuccess, Discriminator: test.discriminators[i]})
			}
			result, err := caller([]byte(test.testValue), typeHolders, test.discriminatorField)
			test.Assert(t, result, err, isSuccess)
		})
	}
}
