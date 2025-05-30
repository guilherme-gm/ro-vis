package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/guilherme-gm/ro-vis/extractor/internal/conf"
	"github.com/guilherme-gm/ro-vis/extractor/internal/controllers"
)

func main() {
	conf.LoadApi()

	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("updateStr", controllers.UpdateString)
	}

	router.Use(cors.Default()) // @TODO: Review CORS config
	router.Use(controllers.ErrorHandler())

	itemsController := controllers.ItemsController{}
	updatesController := controllers.UpdatesController{}

	router.Group("api/items/").
		GET("/", controllers.GlobalHandler(itemsController.List)).
		GET("/update/:update", controllers.GlobalHandler(itemsController.ListForUpdate)).
		GET("/:itemId", controllers.GlobalHandler(itemsController.ListForItem))
	router.Group("api/updates/").
		GET("/", controllers.GlobalHandler(updatesController.List))

	router.Run(fmt.Sprintf(":%d", conf.ApiConfig.Port))
}
