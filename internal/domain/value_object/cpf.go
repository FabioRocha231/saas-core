package valueobject

import (
	"errors"
	"fmt"
)

type Cpf struct {
	value string
}

var (
	ErrCpfInvalidLength = errors.New("cpf must have 11 digits")
	ErrCpfInvalid       = errors.New("cpf is invalid")
)

func NewCpf(value string) *Cpf {
	d := digitsOnly(value)
	return &Cpf{value: d}
}

func (c *Cpf) Validate() error {
	if len(c.value) != 11 {
		return ErrCpfInvalidLength
	}
	if allSameDigits(c.value) {
		return ErrCpfInvalid
	}

	// calcula 1º dígito verificador (10 pesos: 10..2)
	d1 := calcCpfDigit(c.value[:9], 10)

	// calcula 2º dígito verificador (11 pesos: 11..2) usando d1
	d2 := calcCpfDigit(c.value[:9]+fmt.Sprintf("%d", d1), 11)

	// compara com os dígitos informados
	if c.value[9] != byte('0'+d1) || c.value[10] != byte('0'+d2) {
		return ErrCpfInvalid
	}
	return nil
}

func (c *Cpf) Masked() (string, error) {
	if len(c.value) != 11 {
		return "", ErrCpfInvalidLength
	}
	v := c.value
	return fmt.Sprintf("%s.%s.%s-%s", v[0:3], v[3:6], v[6:9], v[9:11]), nil
}

func (c *Cpf) Digits() string {
	return c.value
}

func calcCpfDigit(base string, startWeight int) int {
	sum := 0
	w := startWeight
	for i := 0; i < len(base); i++ {
		sum += int(base[i]-'0') * w
		w--
	}
	rem := sum % 11
	if rem < 2 {
		return 0
	}
	return 11 - rem
}
