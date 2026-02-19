package usecase

import (
	"context"
	"testing"

	"github.com/FabioRocha231/saas-core/test/testkit"
	"github.com/stretchr/testify/assert"
)

func TestGetStoreMenuByIDUsecase(t *testing.T) {
	testEnv := testkit.NewEnv()

	storeID, mockStoreErr := testEnv.SeedStore(context.Background(), testEnv.UUID.Generate())
	assert.NoError(t, mockStoreErr)

	storeMenuID, mockStoreMenuErr := testEnv.SeedStoreMenu(context.Background(), storeID)
	assert.NoError(t, mockStoreMenuErr)

	uc := NewGetStoreMenuByIDUsecase(testEnv.StoreMenuRepo, testEnv.UUID)

	t.Run("Should return error if the id is not provided", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), GetStoreMenuByIDInput{})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: store menu id are required")
	})

	t.Run("should return error if the inputed id is invalid", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), GetStoreMenuByIDInput{StoreMenuID: "123"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: invalid store menu id")
	})

	t.Run("Should return error if the store menu does not exist", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), GetStoreMenuByIDInput{StoreMenuID: "a2b24ebb-b79d-450c-a43c-5bfe2a9e7a01"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "not_found: menu not found")
	})

	t.Run("should return store menu if the correct id is provided", func(t *testing.T) {
		output, err := uc.Execute(context.Background(), GetStoreMenuByIDInput{StoreMenuID: storeMenuID})
		assert.NoError(t, err)

		assert.Equal(t, output.ID, storeMenuID)
		assert.Equal(t, output.StoreID, storeID)
		assert.Equal(t, output.Name, "test")
		assert.Equal(t, output.IsActive, true)
	})
}
