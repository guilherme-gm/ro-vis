package dao

import "database/sql"

func ToNullString(val string) sql.NullString {
	return sql.NullString{String: val, Valid: true}
}

func ToNullBool(val bool) sql.NullBool {
	return sql.NullBool{Bool: val, Valid: true}
}

func ToNullInt32(val int32) sql.NullInt32 {
	return sql.NullInt32{Int32: val, Valid: true}
}
