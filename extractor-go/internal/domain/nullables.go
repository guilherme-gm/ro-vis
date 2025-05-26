package domain

import (
	"database/sql"
	"encoding/json"
)

type NullableInt32 sql.NullInt32

func (v NullableInt32) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(v.Int32)
}

type NullableString sql.NullString

func (v NullableString) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(v.String)
}
