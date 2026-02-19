package usecase

import (
	"context"
	"testing"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByIdUsecase(t *testing.T) {
	uuserRepo, uuid, userID, mockUserErr := BootstrapUser()

	assert.NoError(t, mockUserErr)

	uc := NewGetUserByIdUsecase(uuserRepo, uuid)
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
		output, err := uc.Execute(context.Background(), GetUserByIdInput{ID: userID})
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
		assert.NotNil(t, output.User.CreatedAt)
		assert.NotNil(t, output.User.UpdatedAt)
	})
}
