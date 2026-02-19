package usecase

import (
	"context"
	"testing"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	memorystore "github.com/FabioRocha231/saas-core/internal/infra/db/repository/store"
	"github.com/FabioRocha231/saas-core/test/testkit"
	"github.com/stretchr/testify/assert"
)

func TestGetStoreByIDUsecase(t *testing.T) {
	bs, mockUserError := testkit.BootstrapUser()
	assert.NoError(t, mockUserError)
	storeRepo := memorystore.New()

	storeID := bs.UUID.Generate()
	err := storeRepo.Create(context.Background(), &entity.Store{
		ID:      storeID,
		Name:    "test",
		Cnpj:    "12345678901234",
		OwnerID: bs.UserID,
		IsOpen:  true,
		Slug:    "test",
	})
	assert.NoError(t, err)

	uc := NewGetStoreByIDUsecase(storeRepo, bs.UUID)

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
		assert.Equal(t, output.Store.Cnpj, "12345678901234")
		assert.Equal(t, output.Store.OwnerID, bs.UserID)
		assert.Equal(t, output.Store.IsOpen, true)
		assert.Equal(t, output.Store.Slug, "test")
	})
}
