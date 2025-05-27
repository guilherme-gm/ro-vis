package controllers

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type PaginateQuery struct {
	Start int32 `form:"start" binding:"min=0"`
}

var updateFormatRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

var UpdateString validator.Func = func(fl validator.FieldLevel) bool {
	val, ok := fl.Field().Interface().(string)
	if ok {
		return updateFormatRegex.Match([]byte(val))
	}
	return true
}
