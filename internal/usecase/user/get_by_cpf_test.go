package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	memoryuser "github.com/FabioRocha231/saas-core/internal/infra/db/repository/user"
	"github.com/FabioRocha231/saas-core/pkg"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByCpf(t *testing.T) {
	userRepo := memoryuser.New()
	uuid := pkg.NewUUID()
	userID := uuid.Generate()
	createUserErr := userRepo.Create(context.Background(), &entity.User{
		ID:              userID,
		Name:            "usuario teste",
		Email:           "j0Btq@example.com",
		Password:        "123456",
		Cpf:             "74444217065",
		Phone:           "12345678910",
		Role:            entity.UserRoleCostumer,
		Status:          entity.UserStatusActive,
		EmailVerifiedAt: &time.Time{},
		PhoneVerifiedAt: &time.Time{},
		LastLoginAt:     &time.Time{},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	})

	assert.NoError(t, createUserErr)

	uc := NewGetUserByCpfUsecase(userRepo, uuid)
	t.Run("Should return error if the cpf is not provided", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), GetUserByCpfInput{})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: invalid cpf")
	})

	t.Run("should return error if the inputed cpf is minor than 11 characters", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), GetUserByCpfInput{Cpf: "123"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: cpf must have 11 digits")
	})

	t.Run("shoul return error if the inputed cpf is invalid", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), GetUserByCpfInput{Cpf: "12345678901"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: cpf is invalid")
	})

	t.Run("should return user not found if the cpf is valid but user does not exist", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), GetUserByCpfInput{Cpf: "12607011078"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "not_found: user not found")
	})

	t.Run("should return the correct user", func(t *testing.T) {
		output, err := uc.Execute(context.Background(), GetUserByCpfInput{Cpf: "74444217065"})
		assert.NoError(t, err)
		assert.Equal(t, output.User.ID, userID)
		assert.Equal(t, output.User.Name, "usuario teste")
		assert.Equal(t, output.User.Email, "j0Btq@example.com")
		assert.Equal(t, output.User.Cpf, "74444217065")
		assert.Equal(t, output.User.Phone, "12345678910")
		assert.Equal(t, output.User.Role, entity.UserRoleCostumer.String())
		assert.Equal(t, output.User.Status, entity.UserStatusActive.String())
		assert.NotNil(t, output.User.EmailVerifiedAt)
		assert.NotNil(t, output.User.PhoneVerifiedAt)
		assert.NotNil(t, output.User.LastLoginAt)
		assert.NotNil(t, output.User.CreatedAt)
		assert.NotNil(t, output.User.UpdatedAt)
	})
}
