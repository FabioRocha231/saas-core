package usecase

import (
	"context"
	"testing"

	"github.com/FabioRocha231/saas-core/test/testkit"
	"github.com/stretchr/testify/assert"
)

func TestListMenuCategoryByMenuIDUsecase(t *testing.T) {
	testEnv := testkit.NewEnv()

	storeID, mockStoreErr := testEnv.SeedStore(context.Background(), testEnv.UUID.Generate())
	assert.NoError(t, mockStoreErr)

	menuID, mockMenuErr := testEnv.SeedStoreMenu(context.Background(), storeID)
	assert.NoError(t, mockMenuErr)

	categoryID, mockMenuCategoryErr := testEnv.SeedMenuCategory(context.Background(), menuID)
	assert.NoError(t, mockMenuCategoryErr)

	uc := NewListMenuCategoriesByMenuIDUsecase(testEnv.StoreMenuRepo, testEnv.MenuCategoryRepo, testEnv.UUID)

	t.Run("Should return error if the id is not provided", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), ListMenuCategoriesByMenuIDInput{})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: menu id are required")
	})

	t.Run("should return error if the inputed id is invalid", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), ListMenuCategoriesByMenuIDInput{MenuID: "123"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: invalid menu id")
	})

	t.Run("Should return error if the store menu does not exist", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), ListMenuCategoriesByMenuIDInput{MenuID: "a2b24ebb-b79d-450c-a43c-5bfe2a9e7a01"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "not_found: menu not found")
	})

	t.Run("should return store menu if the correct id is provided", func(t *testing.T) {
		output, err := uc.Execute(context.Background(), ListMenuCategoriesByMenuIDInput{MenuID: menuID})
		assert.NoError(t, err)

		assert.Equal(t, len(output.Categories), 1)
		assert.Equal(t, output.Categories[0].ID, categoryID)
		assert.Equal(t, output.Categories[0].MenuID, menuID)
		assert.Equal(t, output.Categories[0].Name, "test")
		assert.Equal(t, output.Categories[0].IsActive, true)
	})
}
