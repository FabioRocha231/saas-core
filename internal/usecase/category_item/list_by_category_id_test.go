package usecase

import (
	"context"
	"testing"

	"github.com/FabioRocha231/saas-core/test/testkit"
	"github.com/stretchr/testify/assert"
)

func TestListCategoryItemByCategoryIDUsecase(t *testing.T) {
	testEnv := testkit.NewEnv()

	storeID, mockErr := testEnv.SeedStore(context.Background(), testEnv.UUID.Generate())
	assert.NoError(t, mockErr)

	menuID, mockMenuErr := testEnv.SeedStoreMenu(context.Background(), storeID)
	assert.NoError(t, mockMenuErr)

	menuCategoryID, mockMenuCategoryErr := testEnv.SeedMenuCategory(context.Background(), menuID)
	assert.NoError(t, mockMenuCategoryErr)

	categoryItemID, mockCategoryItemErr := testEnv.SeedCategoryItem(context.Background(), menuCategoryID)
	assert.NoError(t, mockCategoryItemErr)

	uc := NewListCategoryItemsByCategoryIDUsecase(testEnv.CategoryItemRepo, testEnv.MenuCategoryRepo, testEnv.UUID)

	t.Run("Should return error if the id is not provided", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), ListCategoryItemsByCategoryIDInput{})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: category id are required")
	})

	t.Run("should return error if the inputed id is invalid", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), ListCategoryItemsByCategoryIDInput{CategoryID: "123"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: invalid category id")
	})

	t.Run("Should return error if the menu category does not exist", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), ListCategoryItemsByCategoryIDInput{CategoryID: "a2b24ebb-b79d-450c-a43c-5bfe2a9e7a01"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "not_found: category not found")
	})

	t.Run("Should return a list of category items", func(t *testing.T) {
		res, err := uc.Execute(context.Background(), ListCategoryItemsByCategoryIDInput{CategoryID: menuCategoryID})
		assert.NoError(t, err)
		assert.Equal(t, len(res.Items), 1)
		assert.Equal(t, res.Items[0].ID, categoryItemID)
	})
}
