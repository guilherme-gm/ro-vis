package controllers

import (
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

func intParam(c *gin.Context, name string) (int, error) {
	val := c.Param(name)
	if val == "" {
		return 0, errors.New("missing required param " + name)
	}

	intVal, err := strconv.Atoi(val)
	return intVal, err
}
