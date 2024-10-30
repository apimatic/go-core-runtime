package types_test

import (
	"encoding/json"
	"github.com/apimatic/go-core-runtime/types"
	"reflect"
	"testing"
)

func TestNewOptional(t *testing.T) {
	value := "Optional Value"
	expected := types.Optional[string]{}
	expected.SetValue(&value)
	expected.ShouldSetValue(true)
	result := types.NewOptional(&value)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestEmptyOptional(t *testing.T) {
	expected := types.Optional[any]{}
	expected.SetValue(nil)
	expected.ShouldSetValue(false)
	result := types.EmptyOptional[any]()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestGetterSetters(t *testing.T) {
	value := "Optional Value"
	expected := types.NewOptional[string](&value)
	result := types.Optional[string]{}
	result.SetValue(&value)
	result.ShouldSetValue(true)

	sameValue := expected.Value() == result.Value()
	if expected.IsValueSet() != result.IsValueSet() && !sameValue {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	value := "Optional Value"
	expected := types.Optional[string]{}
	expected.SetValue(&value)
	expected.ShouldSetValue(true)
	type tempStruct struct {
		Optional types.Optional[string] `json:"optional"`
	}
	var result tempStruct
	_ = json.Unmarshal([]byte(`{"optional": "Optional Value"}`), &result)

	if expected.IsValueSet() != result.Optional.IsValueSet() && expected.Value() != result.Optional.Value() {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result.Optional)
	}
}

func TestUnmarshalJSONError(t *testing.T) {
	var result types.Optional[string]
	err := json.Unmarshal([]byte(`{"optional": "Optional Value"}`), &result)
	if err == nil {
		t.Errorf("Failed:\nExpected: Unmarshalling Error \nGot: %v", result)
	}
}
