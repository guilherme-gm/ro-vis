package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/guilherme-gm/ro-vis/extractor/internal/conf"
	"github.com/guilherme-gm/ro-vis/extractor/internal/controllers"
)

func main() {
	conf.LoadApi()

	router := gin.Default()
	router.Use(cors.Default()) // @TODO: Review CORS config

	itemsController := controllers.ItemsController{}
	updatesController := controllers.UpdatesController{}
	router.Group("api/items/").
		GET("/update/:update", itemsController.ListForUpdate)
	router.Group("api/updates/").
		GET("/", updatesController.List)

	router.Run(fmt.Sprintf(":%d", conf.ApiConfig.Port))
}
