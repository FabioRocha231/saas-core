package http

import (
	"os"
	"time"

	memorysession "github.com/FabioRocha231/saas-core/internal/infra/db/repository/session"
	memorystore "github.com/FabioRocha231/saas-core/internal/infra/db/repository/store"
	memorystoremenu "github.com/FabioRocha231/saas-core/internal/infra/db/repository/store_menu"
	memoryuser "github.com/FabioRocha231/saas-core/internal/infra/db/repository/user"
	"github.com/FabioRocha231/saas-core/internal/infra/http/handlers"
	"github.com/FabioRocha231/saas-core/internal/infra/http/middleware"
	"github.com/FabioRocha231/saas-core/pkg"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(engine *gin.Engine) {
	uuid := pkg.NewUUID()
	passwordHash := pkg.NewPasswordHash()
	storeRepo := memorystore.New()
	userRepo := memoryuser.New()
	sessionRepo := memorysession.New()
	storeMenuRepo := memorystoremenu.New()
	jwtService := pkg.NewJwtService(os.Getenv("JWT_SECRET"), 24*time.Hour, "saas-core", uuid)

	storeHandler := handlers.NewStoreHandler(storeRepo, uuid)
	userHandler := handlers.NewUserHandler(userRepo, storeRepo, uuid, passwordHash)
	authHandler := handlers.NewAuthHandler(passwordHash, jwtService, userRepo, sessionRepo, storeRepo)
	storeMenuHandler := handlers.NewStoreMenuHandler(storeRepo, storeMenuRepo, uuid)

	authMiddleware := middleware.NewAuthMiddleware(jwtService, sessionRepo)

	engine.POST("/user", userHandler.Create)

	engine.POST("/login", authHandler.Login)

	protected := engine.Group("/")
	protected.Use(authMiddleware.Middleware)

	// Store routes
	protected.POST("/store", storeHandler.Create)
	protected.GET("/store/id/:id", storeHandler.GetByID)
	protected.POST("/store/:storeId/menu", storeMenuHandler.Create)
	protected.GET("/store/:storeId/menus", storeMenuHandler.ListByStoreID)

	// User routes
	protected.GET("/user/:id", userHandler.GetByID)
	protected.GET("/user/email/:email", userHandler.GetByEmail)
	protected.GET("/user/cpf/:cpf", userHandler.GetByCpf)

	// Menu Store routes
	protected.GET("/menu/:id", storeMenuHandler.GetByID)
}
