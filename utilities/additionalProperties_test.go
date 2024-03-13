package utilities

import (
	"testing"
)

func TestMapAdditionalProperties(t *testing.T) {
	destinationMap := make(map[string]any)
	sourceMap := map[string]any{"Key1": "value1", "Key2": "value2"}

	MapAdditionalProperties(destinationMap, sourceMap)

	if len(destinationMap) != 2 {
		t.Errorf("MapAdditionalProperties: Expected destination map length 2, got %d", len(destinationMap))
	}
}

func TestUnmarshalAdditionalProperties(t *testing.T) {
	input := []byte(`{"key1":"value1","key2":"value2"}`)
	keys := []string{"key1"}

	result, err := UnmarshalAdditionalProperties(input, keys...)

	if err != nil {
		t.Errorf("UnmarshalAdditionalProperties: Expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Errorf("UnmarshalAdditionalProperties: Expected result map length 1, got %d", len(result))
	}

	if result["key2"] != "value2" {
		t.Error("UnmarshalAdditionalProperties: Key 'key2' expected to be in result map")
	}
}
