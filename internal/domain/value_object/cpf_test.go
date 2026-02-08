package valueobject

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCpf_NewCpf_DigitsOnly(t *testing.T) {
	c := NewCpf("529.982.247-25")
	require.Equal(t, "52998224725", c.Digits())
}

func TestCpf_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "valid cpf (masked input)",
			input:   "529.982.247-25",
			wantErr: nil,
		},
		{
			name:    "valid cpf (digits only)",
			input:   "52998224725",
			wantErr: nil,
		},
		{
			name:    "invalid length short",
			input:   "123",
			wantErr: ErrCpfInvalidLength,
		},
		{
			name:    "invalid length long",
			input:   "1234567890123",
			wantErr: ErrCpfInvalidLength,
		},
		{
			name:    "all same digits invalid",
			input:   "111.111.111-11",
			wantErr: ErrCpfInvalid,
		},
		{
			name:    "invalid check digits",
			input:   "529.982.247-26",
			wantErr: ErrCpfInvalid,
		},
		{
			name:    "non-digits are ignored but still must validate length",
			input:   "abc",
			wantErr: ErrCpfInvalidLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCpf(tt.input)
			err := c.Validate()
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCpf_Masked(t *testing.T) {
	t.Run("masked from valid digits", func(t *testing.T) {
		c := NewCpf("52998224725")
		got, err := c.Masked()
		require.NoError(t, err)
		require.Equal(t, "529.982.247-25", got)
	})

	t.Run("masked returns error for invalid length", func(t *testing.T) {
		c := NewCpf("123")
		_, err := c.Masked()
		require.ErrorIs(t, err, ErrCpfInvalidLength)
	})
}

func TestCalcCpfDigit(t *testing.T) {
	// CPF 529.982.247-25:
	// d1 => base "529982247" (peso 10..2) = 2
	// d2 => base "5299822472" (peso 11..2) = 5
	require.Equal(t, 2, calcCpfDigit("529982247", 10))
	require.Equal(t, 5, calcCpfDigit("5299822472", 11))
}
