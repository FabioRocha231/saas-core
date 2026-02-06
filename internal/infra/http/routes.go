package http

import (
	memorystore "github.com/FabioRocha231/saas-core/internal/infra/db/repository/store"
	"github.com/FabioRocha231/saas-core/internal/infra/http/handlers"
	"github.com/FabioRocha231/saas-core/pkg"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(engine *gin.Engine) {
	uuid := pkg.NewUUID()
	storeHandler := handlers.NewStoreHandler(memorystore.New(), uuid)
	engine.POST("/store", storeHandler.Create)
	engine.GET("/store/:id", storeHandler.GetByID)
}
