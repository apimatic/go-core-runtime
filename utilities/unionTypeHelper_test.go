package utilities

import (
	"encoding/json"
	"testing"

	"github.com/apimatic/go-core-runtime/internal"
)

type UnionTypeCase struct {
	name                 string
	types                []any
	isNullableTypes      []bool
	discriminators       []string
	discriminatorField   string
	testValue            string
	expectedValue        string
	expectedType         any
	shouldFail           bool
	expectedErrorMessage string
}

func (u *UnionTypeCase) Assert(t *testing.T, result any, err error, anyTypeHolderSelected bool) {
	if u.shouldFail {
		internal.Nil(t, result)
		internal.False(t, anyTypeHolderSelected)
		internal.EqualError(t, err, u.expectedErrorMessage)
		return
	}
	internal.Nil(t, err)
	internal.True(t, anyTypeHolderSelected)
	internal.IsType(t, u.expectedType, result)
	marshalled, _ := json.Marshal(result)
	if u.expectedValue == "" {
		u.expectedValue = u.testValue
	}
	internal.Equal(t, u.expectedValue, string(marshalled))
}

func TestCommonOneOfAndAnyOfCases(t *testing.T) {
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
			name:            `(bool,int) => nil`,
			types:           []any{new(bool), new(int)},
			isNullableTypes: []bool{false, true},
			testValue:       `null`,
			expectedType:    nil,
		},
		{
			name:            `(Truck,Car) => nil`,
			types:           []any{&Truck{}, &Car{}},
			isNullableTypes: []bool{false, true},
			testValue:       `null`,
			expectedType:    nil,
		},
	}

	//assertCases(t, tests, UnmarshallOneOf)
	assertCases(t, tests, UnmarshallAnyOf)
}

func TestOneOf(t *testing.T) {
	var tests = []UnionTypeCase{
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
		{
			name:         `(Truck,[]Car) => []Car`,
			types:        []any{&Truck{}, &[]Car{}},
			testValue:    `[{"id":2345,"roof":"BIG","type":null}]`,
			expectedType: &[]Car{},
		},
		{
			name:                 `(bool,int) => FAIL`,
			types:                []any{new(bool), new(int)},
			isNullableTypes:      []bool{true, true},
			testValue:            `null`,
			shouldFail:           true,
			expectedErrorMessage: "There are more than one matching types i.e. {*bool and *int} on: null",
		},
		{
			name:                 `(float,int) => FAIL`,
			types:                []any{new(float32), new(int)},
			testValue:            `2345`,
			shouldFail:           true,
			expectedErrorMessage: "There are more than one matching types i.e. {*float32 and *int} on: 2345",
		},
		{
			name:       `(float,int) => FAIL2`,
			types:      []any{new(float32), new(int)},
			testValue:  `"2345"`,
			shouldFail: true,
			expectedErrorMessage: "We could not match any acceptable type from {*float32, *int} on: \"2345\"\n\n" +
				"Error 1:\n  => json: cannot unmarshal string into Go value of type float32\n\n" +
				"Error 2:\n  => json: cannot unmarshal string into Go value of type int",
		},
		{
			name:                 `(Car,Truck) => FAIL`,
			types:                []any{&Car{}, &Truck{}},
			testValue:            `{"id":2345,"weight":"heavy","roof":"BIG"}`,
			shouldFail:           true,
			expectedErrorMessage: "There are more than one matching types i.e. {*utilities.Car and *utilities.Truck} on: {\"id\":2345,\"weight\":\"heavy\",\"roof\":\"BIG\"}",
		},
		{
			name:       `(Car,Truck) => FAIL`,
			types:      []any{&Car{}, &Truck{}},
			testValue:  `{"roof":"BIG"}`,
			shouldFail: true,
			expectedErrorMessage: "We could not match any acceptable type from {*utilities.Car, *utilities.Truck} on: {\"roof\":\"BIG\"}\n\n" +
				"Error 1:\n  => Car \n\t=> required field `Id` is missing\n\n" +
				"Error 2:\n  => Truck \n\t=> required field `Id` is missing\n\t=> required field `Weight` is missing",
		},
		{
			name:       `(Car,Truck) => FAIL2`,
			types:      []any{&Car{}, &Truck{}},
			testValue:  `"car or truck"`,
			shouldFail: true,
			expectedErrorMessage: "We could not match any acceptable type from {*utilities.Car, *utilities.Truck} on: \"car or truck\"\n\n" +
				"Error 1:\n  => Car \n\t=> json: cannot unmarshal string into Go value of type utilities.car\n\n" +
				"Error 2:\n  => Truck \n\t=> json: cannot unmarshal string into Go value of type utilities.truck",
		},
		{
			name:       `(Car,Truck) => FAIL3`,
			types:      []any{&Car{}, &Truck{}},
			testValue:  `null`,
			shouldFail: true,
			expectedErrorMessage: "We could not match any acceptable type from {*utilities.Car, *utilities.Truck} on: null\n\n" +
				"Error 1:\n  => json: cannot unmarshal null into Go value of type *utilities.Car\n\n" +
				"Error 2:\n  => json: cannot unmarshal null into Go value of type *utilities.Truck",
		},
	}

	assertCases(t, tests, UnmarshallOneOf)
}

func TestOneOfDiscriminator(t *testing.T) {
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
			expectedValue:      `{"id":2345,"roof":"BIG","air_level":{},"type":"2 wheeler"}`,
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

	assertDiscriminatorCases(t, tests, UnmarshallOneOfWithDiscriminator)
}

func TestAnyOf(t *testing.T) {
	var tests = []UnionTypeCase{
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
		{
			name:       `(bool,int) => FAIL`,
			types:      []any{new(bool), new(int)},
			testValue:  `null`,
			shouldFail: true,
			expectedErrorMessage: "We could not match any acceptable type from {*bool, *int} on: null\n\n" +
				"Error 1:\n  => json: cannot unmarshal null into Go value of type *bool\n\n" +
				"Error 2:\n  => json: cannot unmarshal null into Go value of type *int",
		},
		{
			name:       `(Bike,Atom) => FAIL`,
			types:      []any{&Bike{}, &Atom{}},
			testValue:  `{"id":2345,"roof":"BIG","air_level":{},"type":"2 wheeler","number_of_protons":1234}`,
			shouldFail: true,
			expectedErrorMessage: "We could not match any acceptable type from {*utilities.Bike, *utilities.Atom} on: {\"id\":2345,\"roof\":\"BIG\",\"air_level\":{},\"type\":\"2 wheeler\",\"number_of_protons\":1234}\n\n" +
				"Error 1:\n  => Bike . Atom \n\t=> required field `NumberOfElectrons` is missing\n\t=> required field `NumberOfProtons` is missing\n\n" +
				"Error 2:\n  => Atom \n\t=> required field `NumberOfElectrons` is missing",
		},
	}

	assertCases(t, tests, UnmarshallAnyOf)
}

func TestAnyOfDiscriminator(t *testing.T) {
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
		{
			name:               `(Car,Bike) => Car3`,
			types:              []any{&Car{}, &Bike{}},
			discriminators:     []string{"4 wheeler", "2 wheeler"},
			discriminatorField: "",
			testValue:          `{"id":2345,"roof":"BIG","type":"2 wheeler"}`,
			expectedType:       &Car{},
		},
		{
			name:               `(Car,Bike) => Car4`,
			types:              []any{&Car{}, &Bike{}},
			discriminators:     []string{"4 wheeler", ""},
			discriminatorField: "type",
			testValue:          `{"id":2345,"roof":"BIG","type":""}`,
			expectedType:       &Car{},
		},
		{
			name:               `(Car,Bike) => Bike`,
			types:              []any{&Car{}, &Bike{}},
			discriminators:     []string{"", "2 wheeler"},
			discriminatorField: "type",
			testValue:          `{"id":2345,"roof":"BIG","type":"2 wheeler"}`,
			expectedValue:      `{"id":2345,"roof":"BIG","air_level":{},"type":"2 wheeler"}`,
			expectedType:       &Bike{},
		},
	}

	assertDiscriminatorCases(t, tests, UnmarshallAnyOfWithDiscriminator)
}

func assertCases(t *testing.T, tests []UnionTypeCase, caller func([]byte, ...*TypeHolder) (any, error)) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var anyTypeHolderSelected bool
			var typeHolders []*TypeHolder
			for i, tt := range test.types {
				var isNullableType bool = false
				if len(test.isNullableTypes) > 0 {
					isNullableType = test.isNullableTypes[i]
				}
				typeHolders = append(typeHolders, NewTypeHolder(tt, isNullableType, &anyTypeHolderSelected))
			}
			result, err := caller([]byte(test.testValue), typeHolders...)
			test.Assert(t, result, err, anyTypeHolderSelected)
		})
	}
}

func assertDiscriminatorCases(t *testing.T, tests []UnionTypeCase, caller func([]byte, string, ...*TypeHolder) (any, error)) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var anyTypeHolderSelected bool
			var typeHolders []*TypeHolder
			for i, tt := range test.types {
				var isNullableType bool = false
				if len(test.isNullableTypes) > 0 {
					isNullableType = test.isNullableTypes[i]
				}
				typeHolders = append(typeHolders, NewTypeHolderDiscriminator(tt, isNullableType, &anyTypeHolderSelected, test.discriminators[i]))
			}
			result, err := caller([]byte(test.testValue), test.discriminatorField, typeHolders...)
			test.Assert(t, result, err, anyTypeHolderSelected)
		})
	}
}
