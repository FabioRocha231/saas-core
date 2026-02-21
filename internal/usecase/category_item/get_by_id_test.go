package usecase

import (
	"context"
	"testing"

	"github.com/FabioRocha231/saas-core/test/testkit"
	"github.com/stretchr/testify/assert"
)

func TestGetCategoryItemByIDUsecase(t *testing.T) {
	testEnv := testkit.NewEnv()

	storeID, mockErr := testEnv.SeedStore(context.Background(), testEnv.UUID.Generate())
	assert.NoError(t, mockErr)

	menuID, mockMenuErr := testEnv.SeedStoreMenu(context.Background(), storeID)
	assert.NoError(t, mockMenuErr)

	menuCategoryID, mockMenuCategoryErr := testEnv.SeedMenuCategory(context.Background(), menuID)
	assert.NoError(t, mockMenuCategoryErr)

	categoryItemID, mockCategoryItemErr := testEnv.SeedCategoryItem(context.Background(), menuCategoryID)
	assert.NoError(t, mockCategoryItemErr)

	uc := NewGetCategoryItemByIDUsecase(testEnv.CategoryItemRepo, testEnv.UUID)

	t.Run("Should return error if the id is not provided", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), GetCategoryItemByIDInput{})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: category item id are required")
	})

	t.Run("should return error if the inputed id is invalid", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), GetCategoryItemByIDInput{ID: "123"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid_argument: invalid category item id")
	})

	t.Run("Should return error if the category item does not exist", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), GetCategoryItemByIDInput{ID: "a2b24ebb-b79d-450c-a43c-5bfe2a9e7a01"})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "not_found: item not found")
	})

	t.Run("Should return category item", func(t *testing.T) {
		input := GetCategoryItemByIDInput{ID: categoryItemID}
		categoryItem, err := uc.Execute(context.Background(), input)
		assert.NoError(t, err)
		assert.Equal(t, categoryItem.ID, input.ID)
		assert.Equal(t, categoryItem.CategoryID, menuCategoryID)
		assert.Equal(t, categoryItem.Name, "test")
		assert.Equal(t, categoryItem.Description, "description test")
		assert.Equal(t, categoryItem.BasePrice, int64(0))
		assert.Equal(t, categoryItem.ImageURL, "")
		assert.Equal(t, categoryItem.IsActive, true)
	})
}
