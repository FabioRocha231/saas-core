package seed

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

const (
	SeedUserID  = "11111111-1111-1111-1111-111111111111"
	SeedStoreID = "22222222-2222-2222-2222-222222222222"
	SeedMenuID  = "33333333-3333-3333-3333-333333333333"

	SeedCatBurgers = "44444444-4444-4444-4444-444444444444"
	SeedCatDrinks  = "55555555-5555-5555-5555-555555555555"

	SeedItemCheddar = "66666666-6666-6666-6666-666666666666"
	SeedItemClassic = "77777777-7777-7777-7777-777777777777"
	SeedItemCoke    = "88888888-8888-8888-8888-888888888888"

	SeedAddonGroupCheddarAdds  = "99999999-9999-9999-9999-999999999999"
	SeedAddonGroupCheddarSauce = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"

	// Addon Options (UUIDs conhecidos)
	SeedAddonOptAddsBacon   = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	SeedAddonOptAddsCheddar = "cccccccc-cccc-cccc-cccc-cccccccccccc"
	SeedAddonOptAddsOnion   = "dddddddd-dddd-dddd-dddd-dddddddddddd"

	SeedAddonOptSauceHouse = "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
	SeedAddonOptSauceBBQ   = "ffffffff-ffff-ffff-ffff-ffffffffffff"
)

func Seed(
	ctx context.Context,
	userRepo repository.UserRepository,
	storeRepo repository.StoreRepository,
	menuRepo repository.StoreMenuRepository,
	categoryRepo repository.MenuCategoryRepository,
	itemRepo repository.CategoryItemRepository,
	addonGroupRepo repository.ItemAddonGroupRepository,
	addonOptionRepo repository.AddonOptionRepository,
	password ports.PasswordHashInterface,
) {
	if os.Getenv("APP_ENV") != "dev" {
		return
	}

	now := time.Now()

	// 1) USER (idempotente por email)
	const email = "teste@gmail.com"

	u, err := userRepo.GetByMail(ctx, email)
	if err != nil {
		hash, err := password.Hash("123456")
		if err != nil {
			log.Printf("seed: password hash error: %v", err)
			return
		}

		u = &entity.User{
			ID:       SeedUserID,
			Name:     "Nome teste",
			Cpf:      "85608186001",
			Email:    email,
			Phone:    "11999999999",
			Password: hash,
			Role:     entity.UserRoleStoreOwner,
			Status:   entity.UserStatusActive,
		}

		if err := userRepo.Create(ctx, u); err != nil {
			log.Printf("seed: create user error: %v", err)
			return
		}
	}

	// 2) STORE (se owner não tem store, cria a store fixa)
	stores, err := storeRepo.ListByOwnerID(ctx, u.ID)
	if err != nil {
		log.Printf("seed: list stores error: %v", err)
		return
	}

	var s *entity.Store
	if len(stores) > 0 {
		s = stores[0]
	} else {
		s = &entity.Store{
			ID:      SeedStoreID,
			Name:    "Loja Teste",
			Slug:    slugify("Loja Teste"),
			IsOpen:  true,
			Cnpj:    "19131243000197",
			OwnerID: u.ID,
		}

		if err := storeRepo.Create(ctx, s); err != nil {
			log.Printf("seed: create store error: %v", err)
			return
		}
	}

	// 3) MENU (se store não tem menu, cria o menu fixo)
	menus, err := menuRepo.ListByStoreID(ctx, s.ID)
	if err != nil {
		log.Printf("seed: list menus error: %v", err)
		return
	}

	var m *entity.StoreMenu
	if len(menus) > 0 {
		m = menus[0]
	} else {
		m = &entity.StoreMenu{
			ID:        SeedMenuID,
			StoreID:   s.ID,
			Name:      "Cardápio Principal",
			IsActive:  true,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := menuRepo.Create(ctx, m); err != nil {
			log.Printf("seed: create menu error: %v", err)
			return
		}
	}

	// 4) CATEGORIES (se menu já tem categoria, não duplica)
	cats, err := categoryRepo.ListByMenuID(ctx, m.ID)
	if err != nil {
		log.Printf("seed: list categories error: %v", err)
		return
	}

	if len(cats) == 0 {
		seedCategoriesItemsGroupsAndOptions(ctx, categoryRepo, itemRepo, addonGroupRepo, addonOptionRepo, now, m.ID)
	}

	log.Printf("seed ok: user=%s store=%s menu=%s", u.ID, s.ID, m.ID)
}

func seedCategoriesItemsGroupsAndOptions(
	ctx context.Context,
	categoryRepo repository.MenuCategoryRepository,
	itemRepo repository.CategoryItemRepository,
	addonGroupRepo repository.ItemAddonGroupRepository,
	addonOptionRepo repository.AddonOptionRepository,
	now time.Time,
	menuID string,
) {
	// Burgers
	c1 := &entity.MenuCategory{
		ID:        SeedCatBurgers,
		MenuID:    menuID,
		Name:      "Burgers",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := categoryRepo.Create(ctx, c1); err != nil {
		log.Printf("seed: create category error: %v", err)
		return
	}

	// Itens
	cheddar := &entity.CategoryItem{
		ID:          SeedItemCheddar,
		CategoryID:  c1.ID,
		Name:        "Cheddar Bacon",
		Description: "Pão brioche, burger 180g, cheddar, bacon e molho da casa.",
		BasePrice:   3990,
		ImageURL:    "",
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	_ = itemRepo.Create(ctx, cheddar)

	_ = itemRepo.Create(ctx, &entity.CategoryItem{
		ID:          SeedItemClassic,
		CategoryID:  c1.ID,
		Name:        "Classic",
		Description: "Pão, burger 180g, queijo e molho especial.",
		BasePrice:   2990,
		ImageURL:    "",
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	})

	// Addon Groups do Cheddar Bacon
	_ = addonGroupRepo.Create(ctx, &entity.ItemAddonGroup{
		ID:             SeedAddonGroupCheddarAdds,
		CategoryItemID: cheddar.ID,
		Name:           "Adicionais",
		Required:       false,
		MinSelect:      0,
		MaxSelect:      3,
		Order:          1,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	})

	_ = addonGroupRepo.Create(ctx, &entity.ItemAddonGroup{
		ID:             SeedAddonGroupCheddarSauce,
		CategoryItemID: cheddar.ID,
		Name:           "Molhos",
		Required:       false,
		MinSelect:      0,
		MaxSelect:      2,
		Order:          2,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	})

	// Addon Options - Adicionais
	_ = addonOptionRepo.Create(ctx, &entity.AddonOption{
		ID:           SeedAddonOptAddsBacon,
		AddonGroupID: SeedAddonGroupCheddarAdds,
		Name:         "Bacon",
		Price:        500,
		Order:        1,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	_ = addonOptionRepo.Create(ctx, &entity.AddonOption{
		ID:           SeedAddonOptAddsCheddar,
		AddonGroupID: SeedAddonGroupCheddarAdds,
		Name:         "Cheddar extra",
		Price:        400,
		Order:        2,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	_ = addonOptionRepo.Create(ctx, &entity.AddonOption{
		ID:           SeedAddonOptAddsOnion,
		AddonGroupID: SeedAddonGroupCheddarAdds,
		Name:         "Cebola caramelizada",
		Price:        350,
		Order:        3,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	})

	// Addon Options - Molhos
	_ = addonOptionRepo.Create(ctx, &entity.AddonOption{
		ID:           SeedAddonOptSauceHouse,
		AddonGroupID: SeedAddonGroupCheddarSauce,
		Name:         "Molho da casa",
		Price:        0,
		Order:        1,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	_ = addonOptionRepo.Create(ctx, &entity.AddonOption{
		ID:           SeedAddonOptSauceBBQ,
		AddonGroupID: SeedAddonGroupCheddarSauce,
		Name:         "Barbecue",
		Price:        150,
		Order:        2,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	})

	// Bebidas
	c2 := &entity.MenuCategory{
		ID:        SeedCatDrinks,
		MenuID:    menuID,
		Name:      "Bebidas",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := categoryRepo.Create(ctx, c2); err != nil {
		log.Printf("seed: create category error: %v", err)
		return
	}

	_ = itemRepo.Create(ctx, &entity.CategoryItem{
		ID:          SeedItemCoke,
		CategoryID:  c2.ID,
		Name:        "Coca-Cola Lata",
		Description: "350ml",
		BasePrice:   650,
		ImageURL:    "",
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
}

func slugify(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
