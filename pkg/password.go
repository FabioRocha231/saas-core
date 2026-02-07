package pkg

import (
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"golang.org/x/crypto/bcrypt"
)

const cost = 12

type PasswordHash struct{}

func NewPasswordHash() ports.PasswordHashInterface {
	return &PasswordHash{}
}

func (p *PasswordHash) Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (p *PasswordHash) Verify(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
