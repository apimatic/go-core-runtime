package https

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func getValueFromJSON(rawJSON []byte, jsonPtr string) any {
	var jsonBody any
	if err := json.Unmarshal(rawJSON, &jsonBody); err != nil {
		return ""
	}
	refTokens, err := parseJsonPtr(jsonPtr)
	if err != nil {
		return ""
	}
	val, kind, err := getValueFromJSONPtr(refTokens, jsonBody)
	if err != nil {
		return ""
	}
	switch kind {
	case reflect.Map:
		obj, err := json.Marshal(val)
		if err != nil {
			return ""
		}

		return string(obj)
	}
	return val
}

func parseJsonPtr(jsonPtrStr string) ([]string, error) {
	ptrSeparator := `/`
	var err error
	var referenceTokens []string

	if jsonPtrStr != `` {
		if !strings.HasPrefix(jsonPtrStr, ptrSeparator) {
			err = errors.New(`JSON pointer must be empty or start with a "` + ptrSeparator)
		} else {
			refTokens := strings.Split(jsonPtrStr, ptrSeparator)
			referenceTokens = append(referenceTokens, refTokens[1:]...)
		}
	}
	return referenceTokens, err
}

func getValueFromJSONPtr(referenceTokens []string, node any) (any, reflect.Kind, error) {

	kind := reflect.Invalid
	for _, token := range referenceTokens {
		decodedToken := Unescape(token)
		r, knd, err := getSingleImpl(node, decodedToken)
		if err != nil {
			return nil, knd, err
		}
		node = r
	}
	rValue := reflect.ValueOf(node)
	kind = rValue.Kind()
	return node, kind, nil
}

// Unescape unescapes a json pointer reference token string to the original representation
func Unescape(token string) string {
	step1 := strings.ReplaceAll(token, `~1`, `/`)
	step2 := strings.ReplaceAll(step1, `~0`, `~`)
	return step2
}

func isNil(input any) bool {
	switch reflect.TypeOf(input).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan:
		return reflect.ValueOf(input).IsNil()
	default:
		return false
	}
}

func getSingleImpl(node any, decodedToken string) (any, reflect.Kind, error) {
	rValue := reflect.Indirect(reflect.ValueOf(node))
	kind := rValue.Kind()
	if isNil(node) {
		return nil, kind, fmt.Errorf("nil value has not field %q", decodedToken)
	}

	switch typed := node.(type) {
	case *any: // case of a pointer to interface, that is not resolved by reflect.Indirect
		return getSingleImpl(*typed, decodedToken)
	}

	switch kind {
	case reflect.Map:
		kv := reflect.ValueOf(decodedToken)
		mv := rValue.MapIndex(kv)

		if mv.IsValid() {
			return mv.Interface(), kind, nil
		}
		return nil, kind, fmt.Errorf("object has no key %q", decodedToken)

	case reflect.Slice:
		tokenIndex, err := strconv.Atoi(decodedToken)
		if err != nil {
			return nil, kind, err
		}
		sLength := rValue.Len()
		if tokenIndex < 0 || tokenIndex >= sLength {
			return nil, kind, fmt.Errorf("index out of bounds array[0,%d] index '%d'", sLength-1, tokenIndex)
		}

		elem := rValue.Index(tokenIndex)
		return elem.Interface(), kind, nil

	default:
		return nil, kind, fmt.Errorf("invalid token reference %q", decodedToken)
	}

}
