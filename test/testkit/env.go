package testkit

import (
	memorycategoryitem "github.com/FabioRocha231/saas-core/internal/infra/db/repository/category_item"
	memorymenucategory "github.com/FabioRocha231/saas-core/internal/infra/db/repository/menu_category"
	memorystore "github.com/FabioRocha231/saas-core/internal/infra/db/repository/store"
	memorystoremenu "github.com/FabioRocha231/saas-core/internal/infra/db/repository/store_menu"
	memoryuser "github.com/FabioRocha231/saas-core/internal/infra/db/repository/user"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
	"github.com/FabioRocha231/saas-core/pkg"
)

type Env struct {
	UUID             ports.UUIDInterface
	UserRepo         repository.UserRepository
	StoreRepo        repository.StoreRepository
	StoreMenuRepo    repository.StoreMenuRepository
	MenuCategoryRepo repository.MenuCategoryRepository
	CategoryItemRepo repository.CategoryItemRepository
}

func NewEnv() *Env {
	userRepo := memoryuser.New()
	storeRepo := memorystore.New()
	storeMenuRepo := memorystoremenu.New()
	menuCategoryRepo := memorymenucategory.New()
	categoryItemRepo := memorycategoryitem.New()
	return &Env{
		UUID:             pkg.NewUUID(),
		UserRepo:         userRepo,
		StoreRepo:        storeRepo,
		StoreMenuRepo:    storeMenuRepo,
		MenuCategoryRepo: menuCategoryRepo,
		CategoryItemRepo: categoryItemRepo,
	}
}
