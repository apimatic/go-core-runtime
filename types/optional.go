package types

import "encoding/json"

// Optional struct leverages any type as optional and nullable.
// Optional.set is true when Optional.value is to be used.
type Optional[T any] struct {
	value *T
	set   bool // set is true when its value is to be used
}

func NewOptional[T any](value *T) Optional[T] {
	return Optional[T]{
		value: value,
		set:   true,
	}
}

func EmptyOptional() Optional[any] {
	return Optional[any]{
		value: nil,
		set:   false,
	}
}

func (o *Optional[T]) Value() *T {
	return o.value
}

func (o *Optional[T]) SetValue(value *T) {
	o.value = value
}

func (o *Optional[T]) ShouldSetValue(set bool) {
	o.set = set
}

func (o *Optional[T]) IsValueSet() bool {
	return o.set
}

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
