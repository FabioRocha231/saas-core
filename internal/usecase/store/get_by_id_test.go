package usecase

import (
	"context"
	"testing"

	"github.com/FabioRocha231/saas-core/test/testkit"
	"github.com/stretchr/testify/assert"
)

func TestGetStoreByIDUsecase(t *testing.T) {
	testEnv := testkit.NewEnv()
	userID, mockUserError := testEnv.SeedUser(context.Background())
	assert.NoError(t, mockUserError)
	storeID, mockStoreError := testEnv.SeedStore(context.Background(), userID)
	assert.NoError(t, mockStoreError)

	uc := NewGetStoreByIDUsecase(testEnv.StoreRepo, testEnv.UUID)

	t.Run("Should return error if the id is not provided", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), GetStoreByIDInput{})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: store id are required")
	})

	t.Run("should return error if the inputed id is invalid", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), GetStoreByIDInput{StoreID: "123"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: invalid store id")
	})

	t.Run("should return store not found if the id is valid but store does not exist", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), GetStoreByIDInput{StoreID: "a2b24ebb-b79d-450c-a43c-5bfe2a9e7a01"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "not_found: store not found")
	})

	t.Run("should return store if the correct id is provided", func(t *testing.T) {
		output, err := uc.Execute(context.Background(), GetStoreByIDInput{StoreID: storeID})
		assert.NoError(t, err)
		assert.Equal(t, output.Store.ID, storeID)
		assert.Equal(t, output.Store.Name, "test")
		assert.Equal(t, output.Store.Cnpj, "46848972000131")
		assert.Equal(t, output.Store.OwnerID, userID)
		assert.Equal(t, output.Store.IsOpen, true)
		assert.Equal(t, output.Store.Slug, "test")
	})
}
