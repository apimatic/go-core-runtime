package https

import (
	"reflect"
	"testing"
)

func TestGetSeparator(t *testing.T) {
	testCases := make(map[ArraySerializationOption]rune)
	testCases[Indexed] = -1
	testCases[UnIndexed] = -1
	testCases[Plain] = -1
	testCases[Tsv] = '\t'
	testCases[Csv] = ','
	testCases[Psv] = '|'
	testCases[100] = -1
	testCases[-1] = -1

	for testCase, expectedValue := range testCases {
		actual := testCase.getSeparator()
		if actual != expectedValue {
			t.Errorf("For option %v, expected separator %c but got %c", ArraySerializationOptionStrings[testCase], testCase, actual)
		}
	}
}

func TestJoinKey(t *testing.T) {
	testCases := []struct {
		option    ArraySerializationOption
		keyPrefix string
		index     any
		expected  string
	}{
		{UnIndexed, "prefix", nil, "prefix[]"},
		{Indexed, "prefix", nil, "prefix"},
		{Plain, "prefix", nil, "prefix"},
		{Tsv, "prefix", nil, "prefix"},
		{Csv, "prefix", nil, "prefix"},
		{Psv, "prefix", nil, "prefix"},
	}

	for _, tc := range testCases {
		actual := tc.option.joinKey(tc.keyPrefix, tc.index)
		if actual != tc.expected {
			t.Errorf("For option %s, keyPrefix %s, and index %v, expected %q but got %q",
				ArraySerializationOptionStrings[tc.option], tc.keyPrefix, tc.index, tc.expected, actual)
		}
	}

	// test cases with index
	// when index is provided the value will always be equal to `"prefix[index]"`
	for _, tc := range testCases {
		actual := tc.option.joinKey(tc.keyPrefix, "index")
		if actual != "prefix[index]" {
			t.Errorf("For option %s, keyPrefix %s, and index %v, expected %q but got %q",
				ArraySerializationOptionStrings[tc.option], tc.keyPrefix, "index", "prefix[index]", actual)
		}
	}
}

func TestAppendMap(t *testing.T) {
	testCases := []struct {
		option   ArraySerializationOption
		result   map[string][]string
		param    map[string][]string
		expected map[string][]string
	}{
		// Test case where result map is empty
		{
			option: Csv,
			result: map[string][]string{},
			param: map[string][]string{
				"key1": {"value1"},
			},
			expected: map[string][]string{
				"key1": {"value1"},
			},
		},
		// Test case where result map is non-empty and separator is set
		{
			option: Tsv,
			result: map[string][]string{
				"key1": {"value1"},
			},
			param: map[string][]string{
				"key1": {"value2"},
			},
			expected: map[string][]string{
				"key1": {"value1\tvalue2"},
			},
		},
		// Test case where result map is non-empty and separator is not set
		{
			option: Indexed,
			result: map[string][]string{
				"key1": {"value1"},
			},
			param: map[string][]string{
				"key1": {"value2"},
			},
			expected: map[string][]string{
				"key1": {"value1", "value2"},
			},
		},
	}

	for _, tc := range testCases {
		tc.option.appendMap(tc.result, tc.param)

		// Check if the result map matches the expected result
		if !reflect.DeepEqual(tc.result, tc.expected) {
			t.Errorf("For option %d, expected %v but got %v", tc.option, tc.expected, tc.result)
		}
	}
}

func TestAppend(t *testing.T) {
	testCases := []struct {
		option   ArraySerializationOption
		result   map[string][]string
		key      string
		value    string
		expected map[string][]string
	}{
		// Test case where result map is empty
		{option: Csv, result: map[string][]string{"key1": {"value1"}}, key: "key1", value: "value2", expected: map[string][]string{"key1": {"value1,value2"}}},
		// Test case where result map is non-empty and separator is set
		{Tsv, map[string][]string{"key1": {"value1"}}, "key1", "value2", map[string][]string{"key1": {"value1\tvalue2"}}},
		{Psv, map[string][]string{"key1": {"value1"}}, "key1", "value2", map[string][]string{"key1": {"value1|value2"}}},
		// Test case where result map is non-empty and separator is not set
		{Indexed, map[string][]string{"key1": {"value1"}}, "key1", "value2", map[string][]string{"key1": {"value1", "value2"}}},
	}

	for _, tc := range testCases {
		tc.option.append(tc.result, tc.key, tc.value)

		// Check if the result map matches the expected result
		for k, v := range tc.expected {
			if !equalSlices(tc.result[k], v) {
				t.Errorf("For key %q, expected %v but got %v", k, v, tc.result[k])
			}
		}
	}
}

func equalSlices(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i, v := range slice1 {
		if v != slice2[i] {
			return false
		}
	}
	return true
}
