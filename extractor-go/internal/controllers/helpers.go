package controllers

import (
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
