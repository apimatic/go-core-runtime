// Package types provides utility types and functions.
// Copyright (c) APIMatic. All rights reserved.
package types

import "encoding/json"

// Optional is a generic struct that allows any type to be used as optional and nullable.
// Optional.set is true when Optional.value is to be used.
type Optional[T any] struct {
	value *T
	set   bool // set is true when its value is to be used
}

// NewOptional creates and returns an Optional instance with the given value set.
func NewOptional[T any](value *T) Optional[T] {
	return Optional[T]{
		value: value,
		set:   true,
	}
}

// EmptyOptional creates and returns an empty Optional instance with no value set.
func EmptyOptional[T any]() Optional[T] {
	return Optional[T]{
		value: nil,
		set:   false,
	}
}

// Value returns the value stored in the Optional. It returns nil if no value is set.
func (o *Optional[T]) Value() *T {
	return o.value
}

// SetValue sets the value of the Optional.
func (o *Optional[T]) SetValue(value *T) {
	o.value = value
}

// ShouldSetValue sets whether the value should be used or not.
func (o *Optional[T]) ShouldSetValue(set bool) {
	o.set = set
}

// IsValueSet returns true if a value is set in the Optional, false otherwise.
func (o *Optional[T]) IsValueSet() bool {
	return o.set
}

// UnmarshalJSON unmarshal the JSON input into the Optional value.
func (o *Optional[T]) UnmarshalJSON(input []byte) error {
	var temp *T
	if input != nil {
		err := json.Unmarshal(input, &temp)
		if err != nil {
			return err
		}
		o.value = temp
	}
	o.set = true

	return nil
}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.value)
}
