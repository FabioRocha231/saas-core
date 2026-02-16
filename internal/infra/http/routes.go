package http

import (
	"context"
	"os"
	"time"

	memoryaddonoption "github.com/FabioRocha231/saas-core/internal/infra/db/repository/addon_option"
	memorycategoryitem "github.com/FabioRocha231/saas-core/internal/infra/db/repository/category_item"
	memoryitemaddongroup "github.com/FabioRocha231/saas-core/internal/infra/db/repository/item_addon_group"
	memoryitemvariantgroup "github.com/FabioRocha231/saas-core/internal/infra/db/repository/item_variant_group"
	memorymenucategory "github.com/FabioRocha231/saas-core/internal/infra/db/repository/menu_category"
	memorymenuread "github.com/FabioRocha231/saas-core/internal/infra/db/repository/menu_read"
	memoryorder "github.com/FabioRocha231/saas-core/internal/infra/db/repository/order"
	memorypayment "github.com/FabioRocha231/saas-core/internal/infra/db/repository/payment"
	memorysession "github.com/FabioRocha231/saas-core/internal/infra/db/repository/session"
	memorystore "github.com/FabioRocha231/saas-core/internal/infra/db/repository/store"
	memorystoremenu "github.com/FabioRocha231/saas-core/internal/infra/db/repository/store_menu"
	memoryuser "github.com/FabioRocha231/saas-core/internal/infra/db/repository/user"
	memoryvariantoption "github.com/FabioRocha231/saas-core/internal/infra/db/repository/variant_option"
	"github.com/FabioRocha231/saas-core/internal/infra/http/handlers"
	"github.com/FabioRocha231/saas-core/internal/infra/http/middleware"
	"github.com/FabioRocha231/saas-core/internal/infra/seed"
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
	menuCategoryRepo := memorymenucategory.New()
	itemCategoryRepo := memorycategoryitem.New()
	itemAddonGroupRepo := memoryitemaddongroup.New()
	addonOptionRepo := memoryaddonoption.New()
	itemVariantGroupRepo := memoryitemvariantgroup.New()
	variantOptionRepo := memoryvariantoption.New()
	orderRepo := memoryorder.New()
	paymentRepo := memorypayment.New()
	menuReadRepo := memorymenuread.New(
		itemCategoryRepo,
		itemAddonGroupRepo,
		addonOptionRepo,
		itemVariantGroupRepo,
		variantOptionRepo,
	)

	seed.Seed(
		context.Background(),
		userRepo,
		storeRepo,
		storeMenuRepo,
		menuCategoryRepo,
		itemCategoryRepo,
		itemAddonGroupRepo,
		addonOptionRepo,
		itemVariantGroupRepo,
		variantOptionRepo,
		passwordHash,
	)

	jwtService := pkg.NewJwtService(os.Getenv("JWT_SECRET"), 24*time.Hour, "saas-core", uuid)

	storeHandler := handlers.NewStoreHandler(storeRepo, uuid)
	userHandler := handlers.NewUserHandler(userRepo, storeRepo, uuid, passwordHash)
	authHandler := handlers.NewAuthHandler(passwordHash, jwtService, userRepo, sessionRepo, storeRepo)
	storeMenuHandler := handlers.NewStoreMenuHandler(storeRepo, storeMenuRepo, uuid)
	menuCategoryHandler := handlers.NewMenuCategoryHandler(menuCategoryRepo, storeMenuRepo, uuid)
	categoryItemHandler := handlers.NewCategoryItemHandler(itemCategoryRepo, menuCategoryRepo, uuid)
	itemAddonGroupHandler := handlers.NewItemAddonGroupHandler(itemAddonGroupRepo, itemCategoryRepo, uuid)
	addonOptionHandler := handlers.NewAddonOptionHandler(addonOptionRepo, itemAddonGroupRepo, uuid)
	itemVariantGroupHandler := handlers.NewItemVariantGroupHandler(itemVariantGroupRepo, itemCategoryRepo, uuid)
	variantOptionHandler := handlers.NewVariantOptionHandler(variantOptionRepo, itemVariantGroupRepo, uuid)
	orderHandler := handlers.NewOrderHandler(orderRepo, menuReadRepo, uuid)
	paymentHandler := handlers.NewPaymentHandler(orderRepo, paymentRepo, uuid)

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

	// Menu category routes
	protected.POST("/menu/:menuId/category", menuCategoryHandler.Create)
	protected.GET("/menu/categories/:menuId", menuCategoryHandler.ListByMenuID)
	protected.GET("/menu/category/:id", menuCategoryHandler.GetByID)

	// Category item routes
	protected.POST("/menu/category/:categoryId/item", categoryItemHandler.Create)
	protected.GET("/menu/category/item/:id", categoryItemHandler.GetByID)
	protected.GET("/menu/category/items/:categoryId", categoryItemHandler.ListByCategoryID)

	// item addon group routes
	protected.POST("/item/:categoryItemId/addon-group", itemAddonGroupHandler.Create)
	protected.GET("/item/addon-group/:id", itemAddonGroupHandler.GetByID)
	protected.GET("/item/:categoryItemId/addon-groups", itemAddonGroupHandler.ListByCategoryItemID)

	// addon option routes
	protected.POST("/addon-group/:itemAddonGroupId/addon-option", addonOptionHandler.Create)
	protected.GET("/addon-option/:id", addonOptionHandler.GetByID)
	protected.GET("/addon-group/:itemAddonGroupId/addon-options", addonOptionHandler.GetByItemAddonGroupID)

	// Item variant group routes
	protected.POST("/item/:categoryItemId/variant-group", itemVariantGroupHandler.Create)
	protected.GET("/item/variant-group/:id", itemVariantGroupHandler.GetByID)
	protected.GET("/item/:categoryItemId/variant-groups", itemVariantGroupHandler.ListByCategoryItemID)

	// Variant option routes
	protected.POST("/variant-group/:itemVariantGroupId/variant-option", variantOptionHandler.Create)
	protected.GET("/variant-option/:id", variantOptionHandler.GetByID)
	protected.GET("/variant-group/:itemVariantGroupId/variant-options", variantOptionHandler.ListByItemVariantGroupID)

	// order routes
	protected.POST("/store/:storeId/order", orderHandler.Create)
	protected.POST("/order/:orderId/item", orderHandler.AddItem)
	protected.GET("/order/:orderId", orderHandler.GetByID)
	protected.PATCH("/order/:orderId/item/:itemId", orderHandler.UpdateItemQty)
	protected.DELETE("/order/:orderId/item/:itemId", orderHandler.RemoveItem)
	protected.PATCH("/order/:orderId/place", orderHandler.PlaceOrder)

	//payment routes
	protected.POST("/order/:orderId/payments", paymentHandler.CreateForOrder)
	protected.GET("/payments/:paymentId", paymentHandler.GetByID)
	protected.POST("/payments/:paymentId/confirm", paymentHandler.Confirm)
	protected.POST("/payments/:paymentId/fail", paymentHandler.Fail)
}
