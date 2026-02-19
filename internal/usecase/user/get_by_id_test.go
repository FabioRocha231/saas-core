package usecase

import (
	"context"
	"testing"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/test/testkit"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByIdUsecase(t *testing.T) {
	bs, mockUserErr := testkit.BootstrapUser()

	assert.NoError(t, mockUserErr)

	uc := NewGetUserByIdUsecase(bs.UserRepo, bs.UUID)
	t.Run("Should return error if the id is not provided", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), GetUserByIdInput{})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: invalid user id")
	})

	t.Run("should return error if the inputed id is invalid", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), GetUserByIdInput{ID: "123"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: invalid user id")
	})

	t.Run("should return user not found if the id is valid but user does not exist", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), GetUserByIdInput{ID: "a2b24ebb-b79d-450c-a43c-5bfe2a9e7a01"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "not_found: user not found")
	})

	t.Run("should return user if the correct id is provided", func(t *testing.T) {
		output, err := uc.Execute(context.Background(), GetUserByIdInput{ID: bs.UserID})
		assert.NoError(t, err)
		assert.Equal(t, output.User.ID, bs.UserID)
		assert.Equal(t, output.User.Name, "usuario teste")
		assert.Equal(t, output.User.Email, "j0Btq@example.com")
		assert.Equal(t, output.User.Cpf, "74444217065")
		assert.Equal(t, output.User.Phone, "12345678910")
		assert.Equal(t, output.User.Role, entity.UserRoleCostumer.String())
		assert.Equal(t, output.User.Status, entity.UserStatusActive.String())
		assert.NotNil(t, output.User.EmailVerifiedAt)
		assert.NotNil(t, output.User.PhoneVerifiedAt)
		assert.NotNil(t, output.User.CreatedAt)
		assert.NotNil(t, output.User.UpdatedAt)
	})
}
