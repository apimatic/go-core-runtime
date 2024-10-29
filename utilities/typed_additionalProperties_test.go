package utilities

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func Test_Float64Vehicle(t *testing.T) {
	// Creating an instance of Vehicle
	testObj := Vehicle[float64]{
		Year:  2022,
		Make:  ToPointer("Porsche"),
		Model: ToPointer("Taycan turbo GT"),
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

	var deserializedObject Vehicle[float64]
	// Deserializing JSON string to struct
	if err := json.Unmarshal([]byte(jsonString), &deserializedObject); err != nil {
		t.Error(err)
	}

	// Verifying if the deserialized object matches the original
	if !reflect.DeepEqual(testObj, deserializedObject) {
		t.Error("Test_Float64_Vehicle struct failed.")
	}
}

func Test_Float64VehicleWhiteSpace(t *testing.T) {
	// Creating an instance of Vehicle with a whitespace key
	testObj := Vehicle[float64]{
		Year:  2022,
		Make:  ToPointer("Porsche"),
		Model: ToPointer("Taycan turbo GT"),
		AdditionalProperties: map[string]float64{
			"      ": 528,
		},
	}

	// Serializing testObj to JSON
	if _, err := json.Marshal(testObj); err != nil {
		fmt.Println(err)
	} else {
		t.Error("Test_Float64_Vehicle for WhiteSpace")
	}
}

func Test_Float64VehicleConflict(t *testing.T) {
	// Creating an instance of Vehicle with a conflicting key
	testObj := Vehicle[float64]{
		Year:  2022,
		Make:  ToPointer("Porsche"),
		Model: ToPointer("Taycan turbo GT"),
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
	testObj := Vehicle[Bike]{
		Year:  2022,
		Make:  ToPointer("Porsche"),
		Model: ToPointer("Taycan turbo GT"),
		AdditionalProperties: map[string]Bike{
			"bike": Bike{
				Id:   2013,
				Roof: ToPointer("Chopper"),
				Type: ToPointer("Yamaha V Max"),
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

	var deserializedObject Vehicle[Bike]
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
	testObj := Vehicle[AnyOfNumberVehicle]{
		Year:  2022,
		Make:  ToPointer("Porsche"),
		Model: ToPointer("Taycan turbo GT"),
		AdditionalProperties: map[string]AnyOfNumberVehicle{
			"top_speed": AnyOfNumberBooleanContainer.FromNumber(290),
			"fav_bike": AnyOfNumberBooleanContainer.FromVehicle(Vehicle[bool]{
				Year:  2013,
				Make:  ToPointer("Yamaha"),
				Model: ToPointer("Chopper V Max"),
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

	var deserializedObject Vehicle[AnyOfNumberVehicle]
	// Deserializing JSON string to struct
	if err := json.Unmarshal([]byte(jsonString), &deserializedObject); err != nil {
		t.Error(err)
	}

	var testMap, objMap = make(map[string]any), make(map[string]any)
	objBytes, _ := json.Marshal(deserializedObject)

	json.Unmarshal(serializedObject, &testMap)
	json.Unmarshal(objBytes, &objMap)

	// Verifying if the deserialized object matches the original
	if !reflect.DeepEqual(testMap, objMap) {
		t.Error("Test_Float64_Vehicle struct failed.")
	}
}
