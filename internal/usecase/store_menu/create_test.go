package usecase

import (
	"context"
	"testing"

	"github.com/FabioRocha231/saas-core/test/testkit"
	"github.com/stretchr/testify/assert"
)

func TestCreateStoreMenuUsecase(t *testing.T) {
	testEnv := testkit.NewEnv()
	storeID, mockStoreErr := testEnv.SeedStore(context.Background(), testEnv.UUID.Generate())
	assert.NoError(t, mockStoreErr)

	uc := NewCreateStoreMenuUsecase(testEnv.StoreRepo, testEnv.StoreMenuRepo, testEnv.UUID)

	t.Run("Should return error if the name is not provided", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), CreateStoreMenuInput{})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: menu name are required")
	})

	t.Run("Should return error if the store id is not provided", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), CreateStoreMenuInput{Name: "test"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: store id are required")
	})

	t.Run("Should return error if the store id is invalid", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), CreateStoreMenuInput{Name: "test", StoreID: "123"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: invalid store id")
	})

	t.Run("Should return error if the store does not exist", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), CreateStoreMenuInput{Name: "test", StoreID: "a2b24ebb-b79d-450c-a43c-5bfe2a9e7a01"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "not_found: store not found")
	})

	t.Run("should create a new store menu", func(t *testing.T) {
		input := CreateStoreMenuInput{
			Name:    "test",
			StoreID: storeID,
		}
		output, err := uc.Execute(context.Background(), input)
		assert.NoError(t, err)
		assert.NotNil(t, output, "expected not nil but nil returned")
	})
}
