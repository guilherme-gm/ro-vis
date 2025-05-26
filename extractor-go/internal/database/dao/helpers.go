package dao

import (
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
)

func ToNullableString(val string) domain.NullableString {
	return domain.NullableString{String: val, Valid: true}
}

func ToNullableInt32(val int32) domain.NullableInt32 {
	return domain.NullableInt32{Int32: val, Valid: true}
}
