package seed

import (
	"context"
	"log"
	"os"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

func Seed(
	ctx context.Context,
	userRepo repository.UserRepository,
	password ports.PasswordHashInterface,
	uuid ports.UUIDInterface,
) {
	if os.Getenv("APP_ENV") != "dev" {
		return
	}

	const email = "teste@gmail.com"

	_, err := userRepo.GetByMail(ctx, email)
	if err == nil {
		return
	}

	hash, err := password.Hash("123456")
	if err != nil {
		log.Printf("seed: password hash error: %v", err)
		return
	}

	u := &entity.User{
		ID:       uuid.Generate(),
		Name:     "Nome teste",
		Cpf:      "85608186001",
		Email:    email,
		Phone:    "11999999999",
		Password: hash,
		Role:     entity.UserRoleCostumer,
		Status:   entity.UserStatusActive,
	}

	if err := userRepo.Create(ctx, u); err != nil {
		log.Printf("seed: create user error: %v", err)
	}
}
