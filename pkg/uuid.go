package pkg

import "github.com/google/uuid"

type UUID struct{}

func NewUUID() *UUID {
	return &UUID{}
}

func (u *UUID) Generate() string {
	return uuid.New().String()
}

func (u *UUID) Validate(value string) bool {
	_, err := uuid.Parse(value)
	return err == nil
}
