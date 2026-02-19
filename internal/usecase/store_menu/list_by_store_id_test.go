package usecase

import (
	"context"
	"testing"

	"github.com/FabioRocha231/saas-core/test/testkit"
	"github.com/stretchr/testify/assert"
)

func TestListStoreMenuByStoreIDUsecase(t *testing.T) {
	testEnv := testkit.NewEnv()

	storeID, mockStoreErr := testEnv.SeedStore(context.Background(), testEnv.UUID.Generate())
	assert.NoError(t, mockStoreErr)

	_, mockStoreMenuErr := testEnv.SeedStoreMenu(context.Background(), storeID)
	assert.NoError(t, mockStoreMenuErr)

	uc := NewListStoreMenuByStoreIDUsecase(testEnv.StoreMenuRepo, testEnv.UUID)

	t.Run("Should return error if the id is not provided", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), ListStoreMenuByStoreIDInput{})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: store id are required")
	})

	t.Run("should return error if the inputed id is invalid", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), ListStoreMenuByStoreIDInput{StoreID: "123"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: invalid store id")
	})

	t.Run("Should return error if the store does not exist", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), ListStoreMenuByStoreIDInput{StoreID: "4644bcf9-c0bd-4144-9000-7b4c912ea213"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "not_found: store not found")
	})

	t.Run("should return store menu if the correct id is provided", func(t *testing.T) {
		output, err := uc.Execute(context.Background(), ListStoreMenuByStoreIDInput{StoreID: storeID})
		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.NotEmpty(t, output.Menus)
	})
}
