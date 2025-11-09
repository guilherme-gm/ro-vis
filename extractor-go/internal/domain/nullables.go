package domain

import (
	"database/sql"
	"encoding/json"
)

type NullableInt32 sql.NullInt32

func NewNullableInt32(value int32) NullableInt32 {
	return NullableInt32{Int32: value, Valid: true}
}

func NewNullableNullInt32() NullableInt32 {
	return NullableInt32{Int32: 0, Valid: false}
}

func (v NullableInt32) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(v.Int32)
}

type NullableInt64 sql.NullInt64

func (v NullableInt64) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(v.Int64)
}

type NullableString sql.NullString

func NewNullableString(value string) NullableString {
	return NullableString{String: value, Valid: true}
}

func NewNullableNullString() NullableString {
	return NullableString{String: "", Valid: false}
}

func (v NullableString) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(v.String)
}

type NullableBool sql.NullBool

func NewNullableBool(value bool) NullableBool {
	return NullableBool{Bool: value, Valid: true}
}

func NewNullableNullBool() NullableBool {
	return NullableBool{Bool: false, Valid: false}
}

func (v NullableBool) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(v.Bool)
}
