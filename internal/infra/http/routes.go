package http

import (
	memorystore "github.com/FabioRocha231/saas-core/internal/infra/db/repository/store"
	memoryuser "github.com/FabioRocha231/saas-core/internal/infra/db/repository/user"
	"github.com/FabioRocha231/saas-core/internal/infra/http/handlers"
	"github.com/FabioRocha231/saas-core/pkg"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(engine *gin.Engine) {
	uuid := pkg.NewUUID()
	storeHandler := handlers.NewStoreHandler(memorystore.New(), uuid)
	userHandler := handlers.NewUserHandler(memoryuser.New(), uuid)
	engine.POST("/store", storeHandler.Create)
	engine.GET("/store/:id", storeHandler.GetByID)

	engine.POST("/user", userHandler.Create)
	engine.GET("/user/:id", userHandler.GetByID)
}
