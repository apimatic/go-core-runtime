package logger

import (
	"encoding/json"
	"reflect"
	"testing"
)

func validateLevelEnumValues(level Level, t *testing.T) {
	bytes, err := json.Marshal(level)
	if err != nil {
		t.Errorf("Unable to marshal level type : %v", err)
	}
	var newLevel Level
	err = json.Unmarshal(bytes, &newLevel)
	if err != nil {
		t.Errorf("Unable to unmarshal bytes into level type : %v", err)
	}

	if !reflect.DeepEqual(level, newLevel) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", level, newLevel)
	}
}

func TestLevelEnumValueERROR(t *testing.T) {
	level := Level(Level_ERROR)
	validateLevelEnumValues(level, t)
}

func TestLevelEnumValueWARN(t *testing.T) {
	level := Level(Level_WARN)
	validateLevelEnumValues(level, t)
}

func TestLevelEnumValueINFO(t *testing.T) {
	level := Level(Level_INFO)
	validateLevelEnumValues(level, t)
}
func TestLevelEnumValueDEBUG(t *testing.T) {
	level := Level(Level_DEBUG)
	validateLevelEnumValues(level, t)
}

func TestLevelEnumValueTRACE(t *testing.T) {
	level := Level(Level_TRACE)
	validateLevelEnumValues(level, t)
}

func TestLevelEnumValueInvalid(t *testing.T) {
	level := Level("Invalid")
	validateLevelEnumValues(level, new(testing.T))
}
