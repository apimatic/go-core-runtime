package utilities

import (
	"encoding/json"
	"fmt"
	"time"
)

// DefaultTime is a struct that implements time.Time with custom formatting.
type DefaultTime struct{ time.Time }

func NewDefaultTime(t time.Time) DefaultTime {
	return DefaultTime{Time: t}
}

func (t DefaultTime) Value() time.Time { return t.Time }

// String returns DefaultTime as string value by following its defined format.
func (t DefaultTime) String() string { return t.Format(DEFAULT_DATE) }

// MarshalJSON implements json.Marshaller interface to customize JSON marshaling for DefaultTime objects.
func (t DefaultTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format(DEFAULT_DATE))
}

// UnmarshalJSON implements json.Unmarshaler interface to customize JSON unmarshalling for DefaultTime objects.
func (t *DefaultTime) UnmarshalJSON(input []byte) (err error) {
	t.Time, err = unmarshalTime(input, DEFAULT_DATE)
	return err
}

// RFC3339Time is a struct that implements time.Time with custom formatting.
type RFC3339Time struct{ time.Time }

func NewRFC3339Time(t time.Time) RFC3339Time {
	return RFC3339Time{Time: t}
}

func (t RFC3339Time) Value() time.Time { return t.Time }

// String returns RFC3339Time as a string following the RFC3339 standard.
func (t RFC3339Time) String() string { return t.Format(time.RFC3339Nano) }

// MarshalJSON implements json.Marshaller interface to customize JSON marshaling for RFC3339Time objects.
func (t RFC3339Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format(time.RFC3339Nano))
}

// UnmarshalJSON implements json.Unmarshaler interface to customize JSON unmarshalling for RFC3339Time objects.
func (t *RFC3339Time) UnmarshalJSON(input []byte) (err error) {
	t.Time, err = unmarshalTime(input, time.RFC3339Nano)
	return err
}

// RFC1123Time is a struct that implements time.Time with custom formatting.
type RFC1123Time struct{ time.Time }

func NewRFC1123Time(t time.Time) RFC1123Time {
	return RFC1123Time{Time: t}
}

func (t RFC1123Time) Value() time.Time { return t.Time }

// String returns RFC1123Time as a string following the RFC1123 standard.
func (t RFC1123Time) String() string { return t.Format(time.RFC1123) }

// MarshalJSON implements json.Marshaller interface to customize JSON marshaling for RFC1123Time objects.
func (t RFC1123Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format(time.RFC1123))
}

// UnmarshalJSON implements json.Unmarshaler interface to customize JSON unmarshalling for RFC1123Time objects.
func (t *RFC1123Time) UnmarshalJSON(input []byte) (err error) {
	t.Time, err = unmarshalTime(input, time.RFC1123)
	return err
}

func unmarshalTime(input []byte, layout string) (time.Time, error) {
	var temp string
	if err := json.Unmarshal(input, &temp); err != nil {
		return time.Time{}, err
	}
	if timeVal, err := time.Parse(layout, temp); err == nil {
		return timeVal, nil
	} else {
		return time.Time{}, err
	}
}

// UnixDateTime is a struct that implements time.Time with custom formatting.
type UnixDateTime struct{ time.Time }

func NewUnixDateTime(t time.Time) UnixDateTime {
	return UnixDateTime{Time: t}
}

func (t UnixDateTime) Value() time.Time { return t.Time }

// String returns UnixDateTime as a string following the Unix standard.
func (u UnixDateTime) String() string {
	return fmt.Sprintf("%v", u.Unix())
}

// MarshalJSON implements json.Marshaller interface to customize JSON marshaling for UnixDateTime objects.
func (u UnixDateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.Unix())
}

// UnmarshalJSON implements json.Unmarshaler interface to customize JSON unmarshalling for UnixDateTime objects.
func (u *UnixDateTime) UnmarshalJSON(input []byte) error {
	var temp int64
	if err := json.Unmarshal(input, &temp); err != nil {
		return err
	}
	u.Time = time.Unix(temp, 0)
	return nil
}

type TimeTypes interface {
	RFC1123Time | RFC3339Time | UnixDateTime | DefaultTime
	Value() time.Time
}

func ObjMapToTimeMap[T TimeTypes](objMap map[string]T) map[string]time.Time {
	finalMap := map[string]time.Time{}
	for k, v := range objMap {
		finalMap[k] = v.Value()
	}
	return finalMap
}

func ObjSliceToTimeSlice[T TimeTypes](objSlice []T) []time.Time {
	finalSlice := make([]time.Time, len(objSlice))
	for k, v := range objSlice {
		finalSlice[k] = v.Value()
	}
	return finalSlice
}

func TimeMapToObjMap[T TimeTypes](timeMap map[string]time.Time) map[string]T {
	finalMap := map[string]T{}
	for k, v := range timeMap {
		finalMap[k] = T{v}
	}
	return finalMap
}

func TimeSliceToObjSlice[T TimeTypes](timeSlice []time.Time) []T {
	finalSlice := make([]T, len(timeSlice))
	for k, v := range timeSlice {
		finalSlice[k] = T{v}
	}
	return finalSlice
}
