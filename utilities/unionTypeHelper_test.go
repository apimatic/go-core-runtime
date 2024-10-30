package utilities_test

import (
	"encoding/json"
	"github.com/apimatic/go-core-runtime/internal"
	"github.com/apimatic/go-core-runtime/internal/assert"
	"github.com/apimatic/go-core-runtime/utilities"
	"testing"
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
		assert.Nil(t, result)
		assert.False(t, anyTypeHolderSelected)
		assert.EqualError(t, err, u.expectedErrorMessage)
		return
	}
	assert.Nil(t, err)
	assert.True(t, anyTypeHolderSelected)
	assert.IsType(t, u.expectedType, result)
	marshalled, _ := json.Marshal(result)
	if u.expectedValue == "" {
		u.expectedValue = u.testValue
	}
	assert.Equal(t, u.expectedValue, string(marshalled))
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
			types:           []any{&internal.Truck{}, &internal.Car{}},
			isNullableTypes: []bool{false, true},
			testValue:       `null`,
			expectedType:    nil,
		},
	}

	//assertCases(t, tests, UnmarshallOneOf)
	assertCases(t, tests, utilities.UnmarshallAnyOf)
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
			types:        []any{new([]int), &internal.Atom{}},
			testValue:    `{"number_of_electrons":2345,"number_of_protons":1234}`,
			expectedType: &internal.Atom{},
		},
		{
			name:         `(int{},Atom) => int{}`,
			types:        []any{new(map[string]int), &internal.Atom{}},
			testValue:    `{"number_of_":2345,"number_of_protons":1234}`,
			expectedType: new(map[string]int),
		},
		{
			name:         `(Truck,Car) => Car`,
			types:        []any{&internal.Truck{}, &internal.Car{}},
			testValue:    `{"id":2345,"roof":"BIG","type":null}`,
			expectedType: &internal.Car{},
		},
		{
			name:         `(Truck,[]Car) => []Car`,
			types:        []any{&internal.Truck{}, &[]internal.Car{}},
			testValue:    `[{"id":2345,"roof":"BIG","type":null}]`,
			expectedType: &[]internal.Car{},
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
			types:                []any{&internal.Car{}, &internal.Truck{}},
			testValue:            `{"id":2345,"weight":"heavy","roof":"BIG"}`,
			shouldFail:           true,
			expectedErrorMessage: "There are more than one matching types i.e. {*internal.Car and *internal.Truck} on: {\"id\":2345,\"weight\":\"heavy\",\"roof\":\"BIG\"}",
		},
		{
			name:       `(Car,Truck) => FAIL`,
			types:      []any{&internal.Car{}, &internal.Truck{}},
			testValue:  `{"roof":"BIG"}`,
			shouldFail: true,
			expectedErrorMessage: "We could not match any acceptable type from {*internal.Car, *internal.Truck} on: {\"roof\":\"BIG\"}\n\n" +
				"Error 1:\n  => Car \n\t=> required field `Id` is missing\n\n" +
				"Error 2:\n  => Truck \n\t=> required field `Id` is missing\n\t=> required field `Weight` is missing",
		},
		{
			name:       `(Car,Truck) => FAIL2`,
			types:      []any{&internal.Car{}, &internal.Truck{}},
			testValue:  `"car or truck"`,
			shouldFail: true,
			expectedErrorMessage: "We could not match any acceptable type from {*internal.Car, *internal.Truck} on: \"car or truck\"\n\n" +
				"Error 1:\n  => Car \n\t=> json: cannot unmarshal string into Go value of type internal.car\n\n" +
				"Error 2:\n  => Truck \n\t=> json: cannot unmarshal string into Go value of type internal.truck",
		},
		{
			name:       `(Car,Truck) => FAIL3`,
			types:      []any{&internal.Car{}, &internal.Truck{}},
			testValue:  `null`,
			shouldFail: true,
			expectedErrorMessage: "We could not match any acceptable type from {*internal.Car, *internal.Truck} on: null\n\n" +
				"Error 1:\n  => json: cannot unmarshal null into Go value of type *internal.Car\n\n" +
				"Error 2:\n  => json: cannot unmarshal null into Go value of type *internal.Truck",
		},
	}

	assertCases(t, tests, utilities.UnmarshallOneOf)
}

func TestOneOfDiscriminator(t *testing.T) {
	var tests = []UnionTypeCase{
		{
			name:               `(Car,Bike) => Car`,
			types:              []any{&internal.Car{}, &internal.Bike{}},
			discriminators:     []string{"4 wheeler", "2 wheeler"},
			discriminatorField: "type",
			testValue:          `{"id":2345,"roof":"BIG","type":"4 wheeler"}`,
			expectedType:       &internal.Car{},
		},
		{
			name:               `(Car,Bike) => Bike`,
			types:              []any{&internal.Car{}, &internal.Bike{}},
			discriminators:     []string{"4 wheeler", "2 wheeler"},
			discriminatorField: "type",
			testValue:          `{"id":2345,"roof":"BIG","type":"2 wheeler"}`,
			expectedValue:      `{"id":2345,"roof":"BIG","air_level":{},"type":"2 wheeler"}`,
			expectedType:       &internal.Bike{},
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

	assertDiscriminatorCases(t, tests, utilities.UnmarshallOneOfWithDiscriminator)
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
			types:        []any{&internal.Truck{}, &internal.Car{}},
			testValue:    `{"id":2345,"weight":"heavy","roof":"BIG"}`,
			expectedType: &internal.Truck{},
		},
		{
			name:          `(Car,Truck) => Car`,
			types:         []any{&internal.Car{}, &internal.Truck{}},
			testValue:     `{"id":2345,"weight":"heavy","roof":"BIG"}`,
			expectedValue: `{"id":2345,"roof":"BIG","type":null}`,
			expectedType:  &internal.Car{},
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
			types:      []any{&internal.Bike{}, &internal.Atom{}},
			testValue:  `{"id":2345,"roof":"BIG","air_level":{},"type":"2 wheeler","number_of_protons":1234}`,
			shouldFail: true,
			expectedErrorMessage: "We could not match any acceptable type from {*internal.Bike, *internal.Atom} on: {\"id\":2345,\"roof\":\"BIG\",\"air_level\":{},\"type\":\"2 wheeler\",\"number_of_protons\":1234}\n\n" +
				"Error 1:\n  => Bike . Atom \n\t=> required field `NumberOfElectrons` is missing\n\t=> required field `NumberOfProtons` is missing\n\n" +
				"Error 2:\n  => Atom \n\t=> required field `NumberOfElectrons` is missing",
		},
	}

	assertCases(t, tests, utilities.UnmarshallAnyOf)
}

func TestAnyOfDiscriminator(t *testing.T) {
	var tests = []UnionTypeCase{
		{
			name:               `(Car,Bike) => Car`,
			types:              []any{&internal.Car{}, &internal.Bike{}},
			discriminators:     []string{"4 wheeler", "2 wheeler"},
			discriminatorField: "type",
			testValue:          `{"id":2345,"roof":"BIG","type":"unknown"}`,
			expectedType:       &internal.Car{},
		},
		{
			name:               `(Car,Bike) => Car2`,
			types:              []any{&internal.Car{}, &internal.Bike{}},
			discriminators:     []string{"4 wheeler", "2 wheeler"},
			discriminatorField: "type",
			testValue:          `{"id":2345,"roof":"BIG"}`,
			expectedValue:      `{"id":2345,"roof":"BIG","type":null}`,
			expectedType:       &internal.Car{},
		},
		{
			name:               `(Car,Bike) => Car3`,
			types:              []any{&internal.Car{}, &internal.Bike{}},
			discriminators:     []string{"4 wheeler", "2 wheeler"},
			discriminatorField: "",
			testValue:          `{"id":2345,"roof":"BIG","type":"2 wheeler"}`,
			expectedType:       &internal.Car{},
		},
		{
			name:               `(Car,Bike) => Car4`,
			types:              []any{&internal.Car{}, &internal.Bike{}},
			discriminators:     []string{"4 wheeler", ""},
			discriminatorField: "type",
			testValue:          `{"id":2345,"roof":"BIG","type":""}`,
			expectedType:       &internal.Car{},
		},
		{
			name:               `(Car,Bike) => Bike`,
			types:              []any{&internal.Car{}, &internal.Bike{}},
			discriminators:     []string{"", "2 wheeler"},
			discriminatorField: "type",
			testValue:          `{"id":2345,"roof":"BIG","type":"2 wheeler"}`,
			expectedValue:      `{"id":2345,"roof":"BIG","air_level":{},"type":"2 wheeler"}`,
			expectedType:       &internal.Bike{},
		},
	}

	assertDiscriminatorCases(t, tests, utilities.UnmarshallAnyOfWithDiscriminator)
}

func assertCases(t *testing.T, tests []UnionTypeCase, caller func([]byte, ...*utilities.TypeHolder) (any, error)) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var anyTypeHolderSelected bool
			var typeHolders []*utilities.TypeHolder
			for i, tt := range test.types {
				var isNullableType = false
				if len(test.isNullableTypes) > 0 {
					isNullableType = test.isNullableTypes[i]
				}
				typeHolders = append(typeHolders, utilities.NewTypeHolder(tt, isNullableType, &anyTypeHolderSelected))
			}
			result, err := caller([]byte(test.testValue), typeHolders...)
			test.Assert(t, result, err, anyTypeHolderSelected)
		})
	}
}

func assertDiscriminatorCases(t *testing.T, tests []UnionTypeCase, caller func([]byte, string, ...*utilities.TypeHolder) (any, error)) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var anyTypeHolderSelected bool
			var typeHolders []*utilities.TypeHolder
			for i, tt := range test.types {
				var isNullableType = false
				if len(test.isNullableTypes) > 0 {
					isNullableType = test.isNullableTypes[i]
				}
				typeHolders = append(typeHolders, utilities.NewTypeHolderDiscriminator(tt, isNullableType, &anyTypeHolderSelected, test.discriminators[i]))
			}
			result, err := caller([]byte(test.testValue), test.discriminatorField, typeHolders...)
			test.Assert(t, result, err, anyTypeHolderSelected)
		})
	}
}
