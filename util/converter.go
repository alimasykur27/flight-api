package util

import (
	"database/sql"
	"time"
)

func Ptr[T any](v T) *T { return &v }

func ToSqlNullString(v *string) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *v, Valid: true}
}

func ToSqlNullBool(v *bool) sql.NullBool {
	if v == nil {
		return sql.NullBool{}
	}
	return sql.NullBool{Bool: *v, Valid: true}
}

func ToSqlNullInt16(v *int16) sql.NullInt16 {
	if v == nil {
		return sql.NullInt16{}
	}
	return sql.NullInt16{Int16: *v, Valid: true}
}

func ToSqlNullInt32(v *int32) sql.NullInt32 {
	if v == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: *v, Valid: true}
}

func ToSqlNullInt64[T int | int64](v *T) sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*v), Valid: true}
}

func ToSqlNullTime(v *time.Time) sql.NullTime {
	if v == nil {
		return sql.NullTime{}
	}

	return sql.NullTime{Time: *v, Valid: true}
}
func FromSqlNull[T any](v interface{}) T {
	switch val := v.(type) {
	case sql.NullString:
		if val.Valid {
			return any(val.String).(T)
		}
	case sql.NullBool:
		if val.Valid {
			return any(val.Bool).(T)
		}
	case sql.NullInt16:
		if val.Valid {
			return any(val.Int16).(T)
		}
	case sql.NullInt32:
		if val.Valid {
			return any(val.Int32).(T)
		}
	case sql.NullInt64:
		if val.Valid {
			return any(val.Int64).(T)
		}
	}
	var zero T
	return zero
}
