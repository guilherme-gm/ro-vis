package controllers

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
)

// GlobalHandler wraps a gin http handler with the ability to automatically bind path/query parameters into a struct.
// The handler must have the following signature:
// func(*gin.Context) or func(*gin.Context, <params>)
// The params struct must have the following fields:
// - Params: binds to path parameters
// - Query: binds to query parameters
// If the handler has a single param, it must be a *gin.Context
// If the handler has two params, the second param must be a struct with the above fields.
// The handler is called with the gin context as the first param, and the second param with the bound values.
func GlobalHandler(handler any) func(*gin.Context) {
	typeOfHandler := reflect.TypeOf(handler)
	if typeOfHandler.Kind() != reflect.Func {
		panic("handler is not a function - type: " + typeOfHandler.String())
	}

	argCount := typeOfHandler.NumIn()
	if argCount == 0 || argCount > 2 {
		panic("handler must have Gin Context as first param and optionally a params as second. no more than that.")
	}

	if argCount == 1 {
		if handlerFn, ok := handler.(func(*gin.Context)); ok {
			return handlerFn
		}

		panic("gin context must be the first parameter.")
	}

	if typeOfHandler.In(0).Elem().Name() != "Context" {
		panic("handler's first param must be a *gin.Context")
	}

	inputParam := typeOfHandler.In(1)
	_, hasParams := inputParam.FieldByName("Params")
	_, hasQuery := inputParam.FieldByName("Query")

	return func(c *gin.Context) {
		newArg := reflect.New(inputParam)
		if hasParams {
			paramsField := newArg.Elem().FieldByName("Params").Addr()
			err := c.ShouldBindUri(paramsField.Interface())
			if err != nil {
				fmt.Println(err)
				c.AbortWithError(400, err)
				return
			}
		}

		if hasQuery {
			queryField := newArg.Elem().FieldByName("Query").Addr()
			err := c.ShouldBindQuery(queryField.Interface())
			if err != nil {
				fmt.Println(err)
				c.AbortWithError(400, err)
				return
			}
		}

		reflect.ValueOf(handler).Call([]reflect.Value{reflect.ValueOf(c), newArg.Elem()})
	}
}
