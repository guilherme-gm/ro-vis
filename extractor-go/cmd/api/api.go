package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/guilherme-gm/ro-vis/extractor/internal/conf"
	"github.com/guilherme-gm/ro-vis/extractor/internal/controllers"
	"github.com/guilherme-gm/ro-vis/extractor/internal/middleware"
)

func main() {
	conf.LoadApi()

	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("updateStr", controllers.UpdateString)
	}

	if gin.Mode() == gin.ReleaseMode {
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowOrigins = []string{"https://guilherme-gm.github.io"}
		corsConfig.AllowHeaders = []string{"x-server", "Origin", "Content-Type", "Accept"}
		corsConfig.AllowMethods = []string{"GET", "OPTIONS"}
		router.Use(cors.New(corsConfig))
	} else {
		router.Use(cors.Default())
	}
	router.Use(controllers.ErrorHandler())
	router.Use(middleware.ServerSelectorMiddleware())

	itemsController := controllers.ItemsController{}
	questController := controllers.QuestController{}
	updatesController := controllers.UpdatesController{}
	i18nController := controllers.I18nController{}
	mapsController := controllers.MapsController{}

	router.Group("api/items/").
		GET("/", controllers.GlobalHandler(itemsController.List)).
		GET("/update/:update", controllers.GlobalHandler(itemsController.ListForUpdate)).
		GET("/:itemId", controllers.GlobalHandler(itemsController.ListForItem))
	router.Group("api/quests/").
		GET("/", controllers.GlobalHandler(questController.List)).
		GET("/update/:update", controllers.GlobalHandler(questController.ListForUpdate)).
		GET("/:questId", controllers.GlobalHandler(questController.ListForItem))
	router.Group("api/updates/").
		GET("/", controllers.GlobalHandler(updatesController.List))

	router.Group("api/i18n/").
		GET("/", controllers.GlobalHandler(i18nController.List)).
		GET("/update/:update", controllers.GlobalHandler(i18nController.ListForUpdate)).
		POST("/text", controllers.GlobalHandler(i18nController.ListStrings)).
		GET("/:i18nId", controllers.GlobalHandler(i18nController.ListForI18n))

	router.Group("api/maps/").
		GET("/", controllers.GlobalHandler(mapsController.List)).
		GET("/update/:update", controllers.GlobalHandler(mapsController.ListForUpdate)).
		GET("/:mapId", controllers.GlobalHandler(mapsController.ListForItem))

	router.Run(fmt.Sprintf(":%d", conf.ApiConfig.Port))
}
