package logger_test

import (
	"encoding/json"
	"github.com/apimatic/go-core-runtime/logger"
	"reflect"
	"testing"
)

func validateLevelEnumValues(level logger.Level, t *testing.T) {
	bytes, err := json.Marshal(level)
	if err != nil {
		t.Errorf("Unable to marshal level type : %v", err)
	}
	var newLevel logger.Level
	err = json.Unmarshal(bytes, &newLevel)
	if err != nil {
		t.Errorf("Unable to unmarshal bytes into level type : %v", err)
	}

	if !reflect.DeepEqual(level, newLevel) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", level, newLevel)
	}
}

func TestLevelEnumValueERROR(t *testing.T) {
	level := logger.Level_ERROR
	validateLevelEnumValues(level, t)
}

func TestLevelEnumValueWARN(t *testing.T) {
	level := logger.Level_WARN
	validateLevelEnumValues(level, t)
}

func TestLevelEnumValueINFO(t *testing.T) {
	level := logger.Level_INFO
	validateLevelEnumValues(level, t)
}
func TestLevelEnumValueDEBUG(t *testing.T) {
	level := logger.Level_DEBUG
	validateLevelEnumValues(level, t)
}

func TestLevelEnumValueTRACE(t *testing.T) {
	level := logger.Level_TRACE
	validateLevelEnumValues(level, t)
}

func TestLevelEnumValueInvalid(t *testing.T) {
	var level logger.Level
	_ = json.Unmarshal([]byte(`"invalid"`), &level)
	validateLevelEnumValues(level, new(testing.T))
}

func TestLevelEnumValueInvalid2(t *testing.T) {
	var level logger.Level
	_ = json.Unmarshal([]byte(`nil`), &level)
	validateLevelEnumValues(level, new(testing.T))
}

func TestLevelEnumValueInvalid3(t *testing.T) {
	var level logger.Level
	_ = json.Unmarshal([]byte(`{"key": "value"}`), &level)
	validateLevelEnumValues(level, new(testing.T))
}
