package utilities_test

import (
	"encoding/json"
	"fmt"
	"github.com/apimatic/go-core-runtime/internal"
	"github.com/apimatic/go-core-runtime/utilities"
	"reflect"
	"testing"
	"time"
)

func TestMapAdditionalProperties(t *testing.T) {
	destinationMap := make(map[string]any)
	sourceMap := map[string]any{"Key1": "value1", "Key2": "value2"}

	utilities.MapAdditionalProperties(destinationMap, sourceMap)

	if len(destinationMap) != 2 {
		t.Errorf("MapAdditionalProperties: Expected destination map length 2, got %d", len(destinationMap))
	}
}

func TestUnmarshalAdditionalProperties(t *testing.T) {
	input := []byte(`{"key1":"value1","key2":"value2"}`)
	keys := []string{"key1"}

	result, err := utilities.UnmarshalAdditionalProperties(input, keys...)

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

func TestUnmarshalAdditionalPropertiesInvalidInput(t *testing.T) {
	input := []byte(`"value2"`)

	_, err := utilities.UnmarshalAdditionalProperties(input)

	if err == nil {
		t.Errorf("UnmarshalAdditionalProperties: Expected error, got %v", err)
	}
}

func TestMergeTimeTypeAdditionalProperties(t *testing.T) {

	input := []byte(`{"requiredProperty": "2023-11-01", "additionalProperty1": "2024-11-01"}`)
	keysToRemove := []string{"requiredProperty"}
	timeMap, _ := utilities.ExtractTimeTypeAdditionalProperties(input, utilities.DEFAULT_DATE, keysToRemove...)

	destinationMap := map[string]any{"Key1": "value1", "Key2": "value2"}
	utilities.MergeTimeTypeAdditionalProperties(destinationMap, timeMap, utilities.DEFAULT_DATE)

	srcMap := map[string]any{"Key1": "value1", "Key2": "value2", "additionalProperty1": "2024-11-01"}

	// Verifying if the deserialized object matches the original
	if !reflect.DeepEqual(destinationMap, srcMap) {
		t.Error("Test_MergeTimeTypeAdditionalProperties struct failed.")
	}
}

func TestMergeTimeTypeAdditionalPropertiesUnix(t *testing.T) {

	input := []byte(`{"requiredProperty": "2023-11-01", "additionalProperty1": 1484719381}`)
	keysToRemove := []string{"requiredProperty"}
	timeMap, _ := utilities.ExtractTimeTypeAdditionalProperties(input, time.UnixDate, keysToRemove...)

	destinationMap := map[string]any{"Key1": "value1", "Key2": "value2"}
	utilities.MergeTimeTypeAdditionalProperties(destinationMap, timeMap, time.UnixDate)

	srcMap := map[string]any{"Key1": "value1", "Key2": "value2", "additionalProperty1": int64(1484719381)}

	// Verifying if the deserialized object matches the original
	if !reflect.DeepEqual(destinationMap, srcMap) {
		t.Error("Test_MergeTimeTypeAdditionalProperties struct failed.")
	}
}

func Test_Float64Vehicle(t *testing.T) {
	// Creating an instance of Vehicle
	testObj := internal.Vehicle[float64]{
		Year:  2022,
		Make:  internal.ToPointer("Porsche"),
		Model: internal.ToPointer("Taycan turbo GT"),
		AdditionalProperties: map[string]float64{
			"top_speed":                 290,
			"Electric range (BEV, ECE)": 605,
			"Acceleration 0 - 100 km/h": 2.3,
		},
	}
	// Serializing testObj to JSON
	if serializedObject, err := json.Marshal(testObj); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("serializedObject: %v\n", string(serializedObject))
	}

	// JSON string to be deserialized
	jsonString := `{"make":"Porsche", "model":"Taycan turbo GT", "year":2022, "top_speed":290, "Acceleration 0 - 100 km/h":2.3, "Electric range (BEV, ECE)":605, "battery energy": "97.0"}`

	var deserializedObject internal.Vehicle[float64]
	// Deserializing JSON string to struct
	if err := json.Unmarshal([]byte(jsonString), &deserializedObject); err != nil {
		t.Error(err)
	}

	// Verifying if the deserialized object matches the original
	if !reflect.DeepEqual(testObj, deserializedObject) {
		t.Error("Test_Float64_Vehicle struct failed.")
	}
}

func Test_Float64VehicleConflict(t *testing.T) {
	// Creating an instance of Vehicle with a conflicting key
	testObj := internal.Vehicle[float64]{
		Year:  2022,
		Make:  internal.ToPointer("Porsche"),
		Model: internal.ToPointer("Taycan turbo GT"),
		AdditionalProperties: map[string]float64{
			"year": 2024,
		},
	}
	// Serializing testObj to JSON
	if _, err := json.Marshal(testObj); err != nil {
		fmt.Println(err)
	} else {
		t.Error("Test_Float64_Vehicle for Conflict")
	}
}

func Test_BikeVehicle(t *testing.T) {
	// Creating an instance of Vehicle
	testObj := internal.Vehicle[internal.Bike]{
		Year:  2022,
		Make:  internal.ToPointer("Porsche"),
		Model: internal.ToPointer("Taycan turbo GT"),
		AdditionalProperties: map[string]internal.Bike{
			"bike": {
				Id:   2013,
				Roof: internal.ToPointer("Chopper"),
				Type: internal.ToPointer("Yamaha V Max"),
			},
		},
	}
	// Serializing testObj to JSON
	if serializedObject, err := json.Marshal(testObj); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("serializedObject: %v\n", string(serializedObject))
	}

	// JSON string to be deserialized
	jsonString := `{"bike":{"id":2013,"roof":"Chopper","type":"Yamaha V Max"},"make":"Porsche","model":"Taycan turbo GT","year":2022}`

	var deserializedObject internal.Vehicle[internal.Bike]
	// Deserializing JSON string to struct
	if err := json.Unmarshal([]byte(jsonString), &deserializedObject); err != nil {
		t.Error(err)
	}

	// Verifying if the deserialized object matches the original
	if !reflect.DeepEqual(testObj, deserializedObject) {
		t.Error("Test_Float64_Vehicle struct failed.")
	}
}

func Test_AnyOfNumberVehicleVehicle(t *testing.T) {
	// Creating an instance of Vehicle
	testObj := internal.Vehicle[internal.AnyOfNumberVehicle]{
		Year:  2022,
		Make:  internal.ToPointer("Porsche"),
		Model: internal.ToPointer("Taycan turbo GT"),
		AdditionalProperties: map[string]internal.AnyOfNumberVehicle{
			"top_speed": internal.AnyOfNumberBooleanContainer.FromNumber(290),
			"fav_bike": internal.AnyOfNumberBooleanContainer.FromVehicle(internal.Vehicle[bool]{
				Year:  2013,
				Make:  internal.ToPointer("Yamaha"),
				Model: internal.ToPointer("Chopper V Max"),
				AdditionalProperties: map[string]bool{
					"is_chopper": true,
				},
			}),
		},
	}
	// Serializing testObj to JSON
	serializedObject, err := json.Marshal(testObj)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("serializedObject: %v\n", string(serializedObject))
	}

	// JSON string to be deserialized
	jsonString := `{"make":"Porsche", "model":"Taycan turbo GT", "year":2022, "top_speed":290, "fav_bike":{"is_chopper":true,"make":"Yamaha","model":"Chopper V Max","year":2013, "addProp1":"invalid"}, "addProp2":"invalid2"}`

	var deserializedObject internal.Vehicle[internal.AnyOfNumberVehicle]
	// Deserializing JSON string to struct
	if err := json.Unmarshal([]byte(jsonString), &deserializedObject); err != nil {
		t.Error(err)
	}

	var testMap, objMap = make(map[string]any), make(map[string]any)
	objBytes, _ := json.Marshal(deserializedObject)

	_ = json.Unmarshal(serializedObject, &testMap)
	_ = json.Unmarshal(objBytes, &objMap)

	// Verifying if the deserialized object matches the original
	if !reflect.DeepEqual(testMap, objMap) {
		t.Error("Test_Float64_Vehicle struct failed.")
	}
}

func Test_Float64VehicleWhiteSpace(t *testing.T) {
	// Creating an instance of Vehicle with a whitespace key
	testObj := internal.Vehicle[float64]{
		Year:  2022,
		Make:  internal.ToPointer("Porsche"),
		Model: internal.ToPointer("Taycan turbo GT"),
		AdditionalProperties: map[string]float64{
			"      ": 528,
		},
	}

	// Serializing testObj to JSON
	if serializedObject, err := json.Marshal(testObj); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("serializedObject: %v\n", string(serializedObject))
	}

	// JSON string to be deserialized
	jsonString := `{"make":"Porsche", "model":"Taycan turbo GT", "year":2022, "      ":528, "battery energy": "97.0"}`

	var deserializedObject internal.Vehicle[float64]
	// Deserializing JSON string to struct
	if err := json.Unmarshal([]byte(jsonString), &deserializedObject); err != nil {
		t.Error(err)
	}

	// Verifying if the deserialized object matches the original
	if !reflect.DeepEqual(testObj, deserializedObject) {
		t.Error("Test_Float64_Vehicle for WhiteSpace")
	}
}
