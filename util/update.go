package util

import (
	"database/sql"
	"time"
)

// UpdateNullString applies updates for a pointer to string to sql.NullString.
func UpdateNullString(dst *sql.NullString, src *string) {
	if src != nil {
		*dst = sql.NullString{String: *src, Valid: true}
	}
}

// UpdateNullBool applies updates for a pointer to bool to sql.NullBool.
func UpdateNullBool(dst *sql.NullBool, src *bool) {
	if src != nil {
		*dst = sql.NullBool{Bool: *src, Valid: true}
	}
}

// UpdateNullInt applies updates for a pointer to int to sql.NullInt64.
func UpdateNullInt(dst *sql.NullInt64, src *int) {
	if src != nil {
		*dst = sql.NullInt64{Int64: int64(*src), Valid: true}
	}
}

// UpdateNullInt64 applies updates for a pointer to int64 to sql.NullInt64.
func UpdateNullInt64(dst *sql.NullInt64, src *int64) {
	if src != nil {
		*dst = sql.NullInt64{Int64: *src, Valid: true}
	}
}

// UpdateString applies updates for a pointer to string.
func UpdateString(dst *string, src *string) {
	if src != nil {
		*dst = *src
	}
}

// UpdateBool applies updates for a pointer to bool.
func UpdateBool(dst *bool, src *bool) {
	if src != nil {
		*dst = *src
	}
}

// UpdateInt applies updates for a pointer to int.
func UpdateInt[T int64 | int32 | int16 | int](dst *T, src *T) {
	if src != nil {
		*dst = *src
	}
}

// UpdateTime applies updates for a pointer to time.Time.
func UpdateTime(dst *time.Time, src *time.Time) {
	if src != nil {
		*dst = *src
	}
}
