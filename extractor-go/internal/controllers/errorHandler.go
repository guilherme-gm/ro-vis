package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		statusCode := http.StatusInternalServerError
		var httpError *HttpError
		if errors.As(err, &httpError) {
			statusCode = httpError.Status
		}

		if statusCode >= 500 {
			fmt.Println("---------- A Server Error happened --------")
			fmt.Printf("Status code: %d\n", statusCode)
			fmt.Printf("Message: %s\n", err.Error())
			if httpError != nil {
				fmt.Printf("Original error: %v\n", httpError.OriginalError)
			}
			fmt.Println("------------------------------------------")
		}

		if statusCode >= 400 && statusCode < 500 {
			fmt.Println("---------- A Client Error happened --------")
			fmt.Printf("Status code: %d\n", statusCode)
			fmt.Printf("Message: %s\n", err.Error())
			if httpError != nil {
				fmt.Printf("Original error: %v\n", httpError.OriginalError)
			}
			fmt.Println("------------------------------------------")
		}

		c.JSON(statusCode, map[string]any{
			"message": http.StatusText(statusCode),
		})
	}
}
