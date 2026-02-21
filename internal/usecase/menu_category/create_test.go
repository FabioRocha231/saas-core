package usecase

import (
	"context"
	"testing"

	"github.com/FabioRocha231/saas-core/test/testkit"
	"github.com/stretchr/testify/assert"
)

func TestCreateMenuCategoryUsecase(t *testing.T) {
	testEnv := testkit.NewEnv()

	storeID, mockStoreErr := testEnv.SeedStore(context.Background(), testEnv.UUID.Generate())
	assert.NoError(t, mockStoreErr)

	menuID, mockMenuErr := testEnv.SeedStoreMenu(context.Background(), storeID)
	assert.NoError(t, mockMenuErr)

	uc := NewCreateMenuCategoryUsecase(testEnv.MenuCategoryRepo, testEnv.StoreMenuRepo, testEnv.UUID)

	t.Run("Should return error if the name is not provided", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), CreateMenuCategoryInput{})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: menu category name are required")
	})

	t.Run("Should return error if the menu id is not provided", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), CreateMenuCategoryInput{Name: "test"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: menu id are required")
	})

	t.Run("Should return error if the menu id is invalid", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), CreateMenuCategoryInput{Name: "test", MenuID: "123"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: invalid menu id")
	})

	t.Run("Should return error if the menu does not exist", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), CreateMenuCategoryInput{Name: "test", MenuID: "a2b24ebb-b79d-450c-a43c-5bfe2a9e7a01"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "not_found: menu not found")
	})

	t.Run("Should create menu category item if the correct menu id is provided", func(t *testing.T) {
		output, err := uc.Execute(context.Background(), CreateMenuCategoryInput{Name: "test", MenuID: menuID, IsActive: true})
		assert.NoError(t, err)
		assert.NotNil(t, output.ID)
	})
}
