package ports

type UUIDInterface interface {
	Generate() string
	Validate(value string) bool
}
