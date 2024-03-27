// Copyright 2013 sigu-399 ( https://github.com/sigu-399 )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// author       sigu-399
// author-github  https://github.com/sigu-399
// author-mail    sigu.399@gmail.com
//
// repository-name  jsonpointer
// repository-desc  An implementation of JSON Pointer - Go language
//
// description    Main and unique file.
//
// created        25-02-2013

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

	emptyPtr := ``
	ptrSeparator := `/`
	invalidStart := `JSON pointer must be empty or start with a "` + ptrSeparator

	var err error
	var referenceTokens []string
	if jsonPtrStr != emptyPtr {
		if !strings.HasPrefix(jsonPtrStr, ptrSeparator) {
			err = errors.New(invalidStart)
		} else {
			refTokens := strings.Split(jsonPtrStr, ptrSeparator)
			referenceTokens = append(referenceTokens, refTokens[1:]...)
		}
	}
	return referenceTokens, err
}

func getValueFromJSONPtr(referenceTokens []string, node any) (any, reflect.Kind, error) {

	kind := reflect.Invalid

	// return full node when tokens are empty
	if len(referenceTokens) == 0 {
		return node, kind, nil
	}

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
	if input == nil {
		return true
	}
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

	switch kind { //nolint:exhaustive
	case reflect.Struct:
		structAny, err := structToAny(node)
		if err != nil {
			return nil, kind, fmt.Errorf("object has no field %q", decodedToken)
		}

		val, _, err := getSingleImpl(structAny, decodedToken)
		return val, kind, err

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
