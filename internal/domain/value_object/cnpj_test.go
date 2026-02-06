package valueobject

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCnpj_NormalizesToDigitsOnly(t *testing.T) {
	input := "04.252.011/0001-10"

	cnpj := NewCnpj(input)
	require.NotNil(t, cnpj)
	require.Equal(t, "04252011000110", cnpj.Digits())
}

func TestValidate_ReturnsError_WhenInvalidLength(t *testing.T) {
	tests := []string{
		"",                      // vazio
		"123",                   // muito curto
		"04.252.011/0001",       // 12 dígitos
		"04.252.011/0001-10000", // longo demais
		"1111111111111",         // 13 dígitos
		"111111111111111",       // 15 dígitos
		"04.252.011/0001-1",     // 13 dígitos após limpar
		"aa.bb.cc/dddd-ee",      // vira "" após limpar
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			cnpj := NewCnpj(tt)
			require.NotNil(t, cnpj)

			err := cnpj.Validate()
			require.ErrorIs(t, err, ErrCnpjInvalidLength)
		})
	}
}

func TestValidate_ReturnsError_WhenAllDigitsSame(t *testing.T) {
	tests := []string{
		"00.000.000/0000-00",
		"11.111.111/1111-11",
		"22222222222222",
		"99999999999999",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			cnpj := NewCnpj(tt)
			require.NotNil(t, cnpj)

			err := cnpj.Validate()
			require.ErrorIs(t, err, ErrCnpjInvalid)
		})
	}
}

func TestValidate_ReturnsError_WhenCheckDigitsInvalid(t *testing.T) {
	// Base válida com dígitos alterados no final
	input := "04.252.011/0001-11" // deveria ser ...-10

	cnpj := NewCnpj(input)
	require.NotNil(t, cnpj)

	err := cnpj.Validate()
	require.ErrorIs(t, err, ErrCnpjInvalid)
}

func TestValidate_ValidCnpj(t *testing.T) {
	cnpj := NewCnpj("04.252.011/0001-10")
	require.NotNil(t, cnpj)

	require.NoError(t, cnpj.Validate())
}

func TestMasked_ReturnsFormatted_WhenValid(t *testing.T) {
	cnpj := NewCnpj("04252011000110")
	require.NotNil(t, cnpj)

	// opcional: garantir que é válido antes de mascarar
	require.NoError(t, cnpj.Validate())

	masked, err := cnpj.Masked()
	require.NoError(t, err)
	require.Equal(t, "04.252.011/0001-10", masked)
}

func TestMasked_ReturnsError_WhenInvalidLength(t *testing.T) {
	cnpj := NewCnpj("123")
	require.NotNil(t, cnpj)

	masked, err := cnpj.Masked()
	require.Empty(t, masked)
	require.ErrorIs(t, err, ErrCnpjInvalidLength)
}

func TestDigits_ReturnsRawDigits(t *testing.T) {
	cnpj := NewCnpj("04.252.011/0001-10")
	require.NotNil(t, cnpj)

	require.Equal(t, "04252011000110", cnpj.Digits())
}