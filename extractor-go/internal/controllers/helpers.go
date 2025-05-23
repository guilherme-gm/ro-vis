package controllers

import (
	"database/sql"
	"errors"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

func queryInt(c *gin.Context, name string, defaultVal int) (int, error) {
	val := c.Query(name)
	if val == "" {
		return defaultVal, nil
	}

	intVal, err := strconv.Atoi(val)
	return intVal, err
}

var updateFormatRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func paramAsUpdate(c *gin.Context, name string) (string, error) {
	val := c.Param(name)
	if val == "" {
		return "", errors.New("missing update param")
	}

	if !updateFormatRegex.Match([]byte(val)) {
		return "", errors.New("invalid format")
	}

	return val, nil
}

type recordResponse[T any] struct {
	Update string
	Data   *T
}

type fromToRecordResponse[T any] struct {
	From *recordResponse[T]
	To   *recordResponse[T]
}

func sqlInt32ToPointer(val sql.NullInt32) *int32 {
	if !val.Valid {
		return nil
	}
	return &val.Int32
}

func sqlStringToPointer(val sql.NullString) *string {
	if !val.Valid {
		return nil
	}
	return &val.String
}
