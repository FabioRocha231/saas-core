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
	storeRepo := memorystore.New()
	storeHandler := handlers.NewStoreHandler(storeRepo, uuid)
	userHandler := handlers.NewUserHandler(memoryuser.New(), storeRepo, uuid)
	engine.POST("/store", storeHandler.Create)
	engine.GET("/store/:id", storeHandler.GetByID)

	engine.POST("/user", userHandler.Create)
	engine.GET("/user/:id", userHandler.GetByID)
	engine.GET("/user/email/:email", userHandler.GetByEmail)
	engine.GET("/user/cpf/:cpf", userHandler.GetByCpf)
}
