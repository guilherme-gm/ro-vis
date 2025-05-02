package dao

import "database/sql"

func ToNullString(val string) sql.NullString {
	return sql.NullString{String: val, Valid: true}
}

func ToNullBool(val bool) sql.NullBool {
	return sql.NullBool{Bool: val, Valid: true}
}
