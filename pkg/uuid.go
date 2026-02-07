package pkg

import (
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/google/uuid"
)

type UUID struct{}

func NewUUID() ports.UUIDInterface {
	return &UUID{}
}

func (u *UUID) Generate() string {
	return uuid.New().String()
}

func (u *UUID) Validate(value string) bool {
	_, err := uuid.Parse(value)
	return err == nil
}
