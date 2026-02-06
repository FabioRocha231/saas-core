package http

import (
	"github.com/FabioRocha231/saas-core/internal/infra/http/handlers"
	"github.com/gin-gonic/gin"
)



func RoutesBootstrap(engine *gin.Engine) {
	// Middlewares padrões (logger + recovery)
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	engine.GET("/health", handlers.HealthHandler)
	engine.POST("/store", handlers.CreateStoreHandler)
}