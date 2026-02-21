package usecase

import (
	"context"
	"testing"

	"github.com/FabioRocha231/saas-core/test/testkit"
	"github.com/stretchr/testify/assert"
)

func TestCreateCategoryItemUsecase(t *testing.T) {
	testEnv := testkit.NewEnv()

	storeID, mockErr := testEnv.SeedStore(context.Background(), testEnv.UUID.Generate())
	assert.NoError(t, mockErr)

	menuID, mockMenuErr := testEnv.SeedStoreMenu(context.Background(), storeID)
	assert.NoError(t, mockMenuErr)

	menuCategoryID, mockMenuCategoryErr := testEnv.SeedMenuCategory(context.Background(), menuID)
	assert.NoError(t, mockMenuCategoryErr)

	uc := NewCreateCategoryItemUsecase(testEnv.CategoryItemRepo, testEnv.MenuCategoryRepo, testEnv.UUID)

	t.Run("Should return error if the name is not provided", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), CreateCategoryItemInput{})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: category item name are required")
	})

	t.Run("Should return error if the menu category id is not provided", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), CreateCategoryItemInput{Name: "test"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: category item description are required")
	})

	t.Run("Should return error if the menu category id is invalid", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), CreateCategoryItemInput{Name: "test", Description: "description test"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: category id are required")
	})

	t.Run("Should return error if the menu category does not exist", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), CreateCategoryItemInput{Name: "test", Description: "description test", CategoryID: "a2b24ebb-b79d-450c-a43c-5bfe2a9e7a01"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "not_found: category not found")
	})

	t.Run("Should create a category item", func(t *testing.T) {
		output, err := uc.Execute(context.Background(), CreateCategoryItemInput{Name: "test", CategoryID: menuCategoryID, Description: "description test"})
		assert.NoError(t, err)
		assert.NotNil(t, output)
	})
}
