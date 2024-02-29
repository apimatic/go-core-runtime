package utilities

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UnionTypeCase struct {
	name         string
	types        []*TypeHolder
	isOneOf      bool
	responseBody string
	expectedType any
}

func TestUnionTypesSuccess(t *testing.T) {
	var isSuccess bool = false
	var tests = []UnionTypeCase{
		{
			name:         `OneOf(string,int) => string1`,
			types:        []*TypeHolder{{Value: new(string), Flag: &isSuccess}, {Value: new(int), Flag: &isSuccess}},
			isOneOf:      true,
			responseBody: `"some string"`,
			expectedType: new(string),
		},
		{
			name:         `AnyOf(string,int) => string1`,
			types:        []*TypeHolder{{Value: new(string), Flag: &isSuccess}, {Value: new(int), Flag: &isSuccess}},
			isOneOf:      false,
			responseBody: `"some string"`,
			expectedType: new(string),
		},
		{
			name:         `OneOf(string,int) => string2`,
			types:        []*TypeHolder{{Value: new(string), Flag: &isSuccess}, {Value: new(int), Flag: &isSuccess}},
			isOneOf:      true,
			responseBody: `"123"`,
			expectedType: new(string),
		},
		{
			name:         `AnyOf(string,int) => string2`,
			types:        []*TypeHolder{{Value: new(string), Flag: &isSuccess}, {Value: new(int), Flag: &isSuccess}},
			isOneOf:      false,
			responseBody: `"123"`,
			expectedType: new(string),
		},
		{
			name:         `OneOf(string,int) => int`,
			types:        []*TypeHolder{{Value: new(string), Flag: &isSuccess}, {Value: new(int), Flag: &isSuccess}},
			isOneOf:      true,
			responseBody: `123`,
			expectedType: new(int),
		},
		{
			name:         `AnyOf(string,int) => int`,
			types:        []*TypeHolder{{Value: new(string), Flag: &isSuccess}, {Value: new(int), Flag: &isSuccess}},
			isOneOf:      false,
			responseBody: `123`,
			expectedType: new(int),
		},
		{
			name:         `OneOf(bool,int) => bool`,
			types:        []*TypeHolder{{Value: new(bool), Flag: &isSuccess}, {Value: new(int), Flag: &isSuccess}},
			isOneOf:      true,
			responseBody: `true`,
			expectedType: new(bool),
		},
		{
			name:         `AnyOf(bool,int) => bool`,
			types:        []*TypeHolder{{Value: new(bool), Flag: &isSuccess}, {Value: new(int), Flag: &isSuccess}},
			isOneOf:      false,
			responseBody: `true`,
			expectedType: new(bool),
		},
		{
			name:         `OneOf(bool,int) => int`,
			types:        []*TypeHolder{{Value: new(bool), Flag: &isSuccess}, {Value: new(int), Flag: &isSuccess}},
			isOneOf:      true,
			responseBody: `2345`,
			expectedType: new(int),
		},
		{
			name:         `AnyOf(bool,int) => int`,
			types:        []*TypeHolder{{Value: new(bool), Flag: &isSuccess}, {Value: new(int), Flag: &isSuccess}},
			isOneOf:      false,
			responseBody: `2345`,
			expectedType: new(int),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedBytes := []byte(test.responseBody)
			result, err := UnmarshallOneOf(expectedBytes, test.types)
			assert.Nil(t, err)
			assert.True(t, isSuccess)
			isSuccess = false
			assert.IsType(t, test.expectedType, result)
			marshalled, _ := json.Marshal(result)
			assert.Equal(t, expectedBytes, marshalled)
		})
	}
}
