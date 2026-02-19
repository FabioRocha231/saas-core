package testkit

import (
	"context"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	memoryuser "github.com/FabioRocha231/saas-core/internal/infra/db/repository/user"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
	"github.com/FabioRocha231/saas-core/pkg"
)

type UserBootstrap struct {
	UserRepo repository.UserRepository
	UUID     ports.UUIDInterface
	UserID   string
}

func BootstrapUser() (*UserBootstrap, error) {
	userRepo := memoryuser.New()
	uuid := pkg.NewUUID()
	userID := uuid.Generate()

	now := time.Now()
	err := userRepo.Create(context.Background(), &entity.User{
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

	if err != nil {
		return nil, err
	}

	return &UserBootstrap{
		UserRepo: userRepo,
		UUID:     uuid,
		UserID:   userID,
	}, nil
}
