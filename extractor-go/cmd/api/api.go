package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/guilherme-gm/ro-vis/extractor/internal/conf"
)

func main() {
	conf.LoadApi()

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Run(fmt.Sprintf(":%d", conf.ApiConfig.Port))
}
