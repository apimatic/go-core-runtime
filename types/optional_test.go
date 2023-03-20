package types

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNewOptional(t *testing.T) {
	value := "Optional Value"
	expected := Optional[string]{value: &value, set: true}
	result := NewOptional(&value)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestEmptyOptional(t *testing.T) {
	expected := Optional[any]{value: nil, set: false}
	result := EmptyOptional()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestGetterSetters(t *testing.T) {
	value := "Optional Value"
	expected := Optional[string]{value: &value, set: true}
	result := Optional[string]{value: nil, set: false}
	result.SetValue(&value)
	result.ShouldSetValue(true)

	sameValue := expected.Value() == result.Value()
	if expected.IsValueSet() != result.IsValueSet() && !sameValue {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	value := "Optional Value"
	expected := Optional[string]{value: &value, set: true}
	type tempstruct struct {
		Optional Optional[string] `json:"optional"`
	}
	var result tempstruct
	json.Unmarshal([]byte(`{"optional": "Optional Value"}`), &result)

	if expected.IsValueSet() != result.Optional.IsValueSet() && expected.Value() != result.Optional.Value() {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result.Optional)
	}
}

func TestUnmarshalJSONError(t *testing.T) {
	var result Optional[string]
	err := json.Unmarshal([]byte(`{"optional": "Optional Value"}`), &result)
	if err == nil {
		t.Errorf("Failed:\nExpected: Unmarshalling Error \nGot: %v", result)
	}
}
