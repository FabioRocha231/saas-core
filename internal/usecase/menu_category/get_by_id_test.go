package usecase

import (
	"context"
	"testing"

	"github.com/FabioRocha231/saas-core/test/testkit"
	"github.com/stretchr/testify/assert"
)

func TestGetMenuCategoryByIDUsecase(t *testing.T) {
	testEnv := testkit.NewEnv()

	storeID, mockStoreErr := testEnv.SeedStore(context.Background(), testEnv.UUID.Generate())
	assert.NoError(t, mockStoreErr)

	menuID, mockMenuErr := testEnv.SeedStoreMenu(context.Background(), storeID)
	assert.NoError(t, mockMenuErr)

	menuCategoryID, mockMenuCategoryErr := testEnv.SeedMenuCategory(context.Background(), menuID)
	assert.NoError(t, mockMenuCategoryErr)

	uc := NewGetMenuCategoryByIDUsecase(testEnv.MenuCategoryRepo, testEnv.UUID)

	t.Run("Should return error if the id is not provided", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), GetMenuCategoryByIDInput{})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: menu category id are required")
	})

	t.Run("should return error if the inputed id is invalid", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), GetMenuCategoryByIDInput{ID: "123"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: invalid menu category id")
	})

	t.Run("Should return error if the menu category does not exist", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), GetMenuCategoryByIDInput{ID: "a2b24ebb-b79d-450c-a43c-5bfe2a9e7a01"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "not_found: category not found")
	})

	t.Run("should return menu category if the correct id is provided", func(t *testing.T) {
		output, err := uc.Execute(context.Background(), GetMenuCategoryByIDInput{ID: menuCategoryID})
		assert.NoError(t, err)
		assert.Equal(t, output.ID, menuCategoryID)
		assert.Equal(t, output.MenuID, menuID)
		assert.Equal(t, output.Name, "test")
		assert.Equal(t, output.IsActive, true)
	})
}
