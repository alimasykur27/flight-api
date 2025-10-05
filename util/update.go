package util

import (
	"time"
)

// UpdateString applies updates for a pointer to string.
func UpdateString(dst **string, src *string) {
	if src != nil {
		*dst = src
	}
}

// UpdateBool applies updates for a pointer to bool.
func UpdateBool(dst **bool, src *bool) {
	if src != nil {
		*dst = src
	}
}

// UpdateInt applies updates for a pointer to int.
func UpdateInt[T int64 | int32 | int16 | int](dst **T, src *T) {
	if src != nil {
		*dst = src
	}
}

// UpdateTime applies updates for a pointer to time.Time.
func UpdateTime(dst **time.Time, src *time.Time) {
	if src != nil {
		*dst = src
	}
}
