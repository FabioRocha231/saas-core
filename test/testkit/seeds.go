package testkit

import (
	"context"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
)

func (e *Env) SeedUser(ctx context.Context) (userID string, err error) {
	userID = e.UUID.Generate()
	now := time.Now()
	err = e.UserRepo.Create(context.Background(), &entity.User{
		ID:              userID,
		Name:            "usuario teste",
		Email:           "j0Btq@example.com",
		Password:        "123456",
		Cpf:             "74444217065",
		Phone:           "12345678910",
		Role:            entity.UserRoleCostumer,
		Status:          entity.UserStatusActive,
		EmailVerifiedAt: &now,
		PhoneVerifiedAt: &now,
		LastLoginAt:     &now,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
	return
}

func (e *Env) SeedStore(ctx context.Context, ownerID string) (storeID string, err error) {
	storeID = e.UUID.Generate()
	err = e.StoreRepo.Create(ctx, &entity.Store{
		ID:      storeID,
		Name:    "test",
		Cnpj:    "46848972000131",
		OwnerID: ownerID,
		IsOpen:  true,
		Slug:    "test",
	})
	return
}

func (e *Env) SeedStoreMenu(ctx context.Context, storeID string) (menuID string, err error) {
	menuID = e.UUID.Generate()
	err = e.StoreMenuRepo.Create(ctx, &entity.StoreMenu{
		ID:        menuID,
		Name:      "test",
		StoreID:   storeID,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	return
}
