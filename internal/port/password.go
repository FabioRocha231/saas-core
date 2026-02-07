package ports

type PasswordHashInterface interface {
	Hash(password string) (string, error)
	Verify(password string, hashedPassword string) bool
}
