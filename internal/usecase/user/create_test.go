package usecase

import (
	"testing"

	memorystore "github.com/FabioRocha231/saas-core/internal/infra/db/repository/store"
	memoryuser "github.com/FabioRocha231/saas-core/internal/infra/db/repository/user"
	"github.com/FabioRocha231/saas-core/pkg"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	userRepo := memoryuser.New()
	storeRepo := memorystore.New()
	uuid := pkg.NewUUID()
	passwordHash := pkg.NewPasswordHash()
	u := NewCreateUserUsecase(userRepo, storeRepo, uuid, passwordHash)

	t.Run("test create a user with a empty input", func(t *testing.T) {
		_, err := u.Execute(t.Context(), CreateUserInput{})

		require.NotNil(t, err, "required an error to be returned when creating a user wihout data but no error was returned")
		require.Equal(t, err.Error(), "invalid_argument: invalid name")
	})

	t.Run("test create a user with a name but without others data", func(t *testing.T) {
		_, err := u.Execute(t.Context(), CreateUserInput{Name: "name"})

		require.NotNil(t, err, "required an error to be returned when creating a user wihout data but no error was returned")
		require.Equal(t, err.Error(), "invalid_argument: invalid email")
	})

	t.Run("test create a user with a name and email but without others data", func(t *testing.T) {
		_, err := u.Execute(t.Context(), CreateUserInput{Name: "name", Email: "email"})

		require.NotNil(t, err, "required an error to be returned when creating a user wihout data but no error was returned")
		require.Equal(t, err.Error(), "invalid_argument: invalid password")
	})

	t.Run("test create a user with a name, email and password but without others data", func(t *testing.T) {
		_, err := u.Execute(t.Context(), CreateUserInput{Name: "name", Email: "email", Password: "password"})

		require.NotNil(t, err, "required an error to be returned when creating a user wihout data but no error was returned")
		require.Equal(t, err.Error(), "invalid_argument: invalid cpf")
	})

	t.Run("test create a user with a name, email, password and invalid cpf", func(t *testing.T) {
		_, err := u.Execute(t.Context(), CreateUserInput{Name: "name", Email: "email", Password: "password", Cpf: "12345678928"})

		require.NotNil(t, err, "required an error to be returned when creating a user wihout data but no error was returned")
		require.Equal(t, err.Error(), "invalid_argument: invalid cpf")
	})

	t.Run("test create a user with a name, email, password, valid cpf and wrong user type", func(t *testing.T) {
		_, err := u.Execute(t.Context(), CreateUserInput{Name: "name", Email: "email", Password: "password", Cpf: "237.566.760-30"})

		require.NotNil(t, err, "required an error to be returned when creating a user wihout data but no error was returned")
		require.Equal(t, err.Error(), "invalid_argument: invalid user type")
	})

	t.Run("test create a user with a name, email, password, valid cpf and valid user type", func(t *testing.T) {
		ouput, err := u.Execute(t.Context(), CreateUserInput{Name: "name", Email: "email", Password: "password", Cpf: "237.566.760-30", UserType: "store"})

		require.Nil(t, err, "required an error to be returned when creating a user wihout data but no error was returned")
		require.NotNil(t, ouput, "required an output to be returned when creating a user wihout data but no output was returned")
	})
}
