package valueobject

import (
	"errors"
	"fmt"
	"unicode"
)

type Cnpj struct {
	value string // sempre armazenado como 14 dígitos (somente números)
}

var (
	ErrCnpjInvalidLength = errors.New("cnpj must have 14 digits")
	ErrCnpjInvalid       = errors.New("cnpj is invalid")
)

func NewCnpj(value string) *Cnpj {
	d := digitsOnly(value)
	return &Cnpj{value: d}
}

// Validate valida tamanho + dígitos verificadores.
func (c *Cnpj) Validate() error {
	if len(c.value) != 14 {
		return ErrCnpjInvalidLength
	}
	if allSameDigits(c.value) {
		return ErrCnpjInvalid
	}

	// calcula 1º dígito verificador
	d1 := calcCnpjDigit(c.value[:12], []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2})
	// calcula 2º dígito verificador
	d2 := calcCnpjDigit(c.value[:12]+fmt.Sprintf("%d", d1), []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2})

	// compara com os dígitos informados
	if c.value[12] != byte('0'+d1) || c.value[13] != byte('0'+d2) {
		return ErrCnpjInvalid
	}
	return nil
}

// Masked retorna o CNPJ mascarado: 00.000.000/0000-00
func (c *Cnpj) Masked() (string, error) {
	if len(c.value) != 14 {
		return "", ErrCnpjInvalidLength
	}
	v := c.value
	return fmt.Sprintf("%s.%s.%s/%s-%s", v[0:2], v[2:5], v[5:8], v[8:12], v[12:14]), nil
}

// Digits retorna o valor "cru" (só números).
func (c *Cnpj) Digits() string {
	return c.value
}

// ---------- helpers ----------

func digitsOnly(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if unicode.IsDigit(r) {
			out = append(out, r)
		}
	}
	return string(out)
}

func allSameDigits(s string) bool {
	if len(s) == 0 {
		return false
	}
	first := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] != first {
			return false
		}
	}
	return true
}

func calcCnpjDigit(base string, weights []int) int {
	sum := 0
	for i := 0; i < len(weights); i++ {
		sum += int(base[i]-'0') * weights[i]
	}
	rem := sum % 11
	if rem < 2 {
		return 0
	}
	return 11 - rem
}