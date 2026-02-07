package http

import (
	memorystore "github.com/FabioRocha231/saas-core/internal/infra/db/repository/store"
	memoryuser "github.com/FabioRocha231/saas-core/internal/infra/db/repository/user"
	"github.com/FabioRocha231/saas-core/internal/infra/http/handlers"
	"github.com/FabioRocha231/saas-core/pkg"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

func RegisterRoutes(engine *gin.Engine) {
	uuid := pkg.NewUUID()
	passwordHash := pkg.NewPasswordHash()
	storeRepo := memorystore.New()
	userRepo := memoryuser.New()
	jwtService := pkg.NewJwtService(os.Getenv("JWT_SECRET"), 24*time.Hour, "saas-core")

	storeHandler := handlers.NewStoreHandler(storeRepo, uuid)
	userHandler := handlers.NewUserHandler(userRepo, storeRepo, uuid, passwordHash)
	authHandler := handlers.NewAuthHandler(passwordHash, jwtService, userRepo)

	engine.POST("/store", storeHandler.Create)

	engine.POST("/user", userHandler.Create)

	engine.POST("/login", authHandler.Login)

	engine.GET("/store/:id", storeHandler.GetByID)

	engine.GET("/user/:id", userHandler.GetByID)
	engine.GET("/user/email/:email", userHandler.GetByEmail)
	engine.GET("/user/cpf/:cpf", userHandler.GetByCpf)
}
