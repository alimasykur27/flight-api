package util

import (
	"strconv"
)

func Ptr[T any](v T) *T { return &v }

func DerefPtr[T any](p *T) T {
	defaultValue := *new(T)
	if p == nil {
		return defaultValue
	}
	return *p
}

func ParseInt64Ptr(s string) *int64 {
	if s == "" {
		return nil
	}

	r, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		return nil
	}

	return &r
}

// ToInterfaces converts any slice of values into []interface{}.
func ToInterfaces[T any](in []T) []interface{} {
	out := make([]interface{}, len(in))
	for i, v := range in {
		out[i] = v
	}
	return out
}
