package usecase

import (
	"context"
	"testing"

	memorystore "github.com/FabioRocha231/saas-core/internal/infra/db/repository/store"
	"github.com/FabioRocha231/saas-core/test/testkit"
	"github.com/stretchr/testify/assert"
)

func TestCreateStoreUsecase(t *testing.T) {
	testEnv := testkit.NewEnv()
	userID, mockUserError := testEnv.SeedUser(context.Background())
	assert.NoError(t, mockUserError)
	storeRepo := memorystore.New()

	uc := NewCreateStoreUsecase(storeRepo, testEnv.UserRepo, testEnv.UUID)

	t.Run("Should return error if the name is not provided", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), CreateStoreInput{})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: store name are required")
	})

	t.Run("Should return error if the cnpj is not provided", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), CreateStoreInput{Name: "test"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: store cnpj are required")
	})

	t.Run("Should return error if the cnpj is minor than 14 characters", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), CreateStoreInput{Name: "test", Cnpj: "123"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: cnpj must have 14 digits")
	})

	t.Run("Should return error if the cnpj is invalid", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), CreateStoreInput{Name: "test", Cnpj: "12345678901234"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: cnpj is invalid")
	})

	t.Run("Should return error if the owner id is not provided", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), CreateStoreInput{Name: "test", Cnpj: "65.921.814/0001-04"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: store owner id are required")
	})

	t.Run("Should return error if the owner id is invalid", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), CreateStoreInput{Name: "test", Cnpj: "65.921.814/0001-04", OwnerID: "123"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: invalid owner id")
	})

	t.Run("Should return error if the owner does not exist", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), CreateStoreInput{Name: "test", Cnpj: "65.921.814/0001-04", OwnerID: "e1cec1f9-baa7-47a5-8830-f433a2933e39"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "not_found: user not found")
	})

	t.Run("Should create store", func(t *testing.T) {
		output, err := uc.Execute(context.Background(), CreateStoreInput{Name: "test", Cnpj: "65.921.814/0001-04", OwnerID: userID})
		assert.NoError(t, err)
		assert.NotNil(t, output.ID, "store id should not be empty")
	})
}
