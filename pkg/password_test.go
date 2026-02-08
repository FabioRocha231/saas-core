package pkg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPasswordHash_HashAndVerify_Success(t *testing.T) {
	p := &PasswordHash{}

	plain := "S3nh@F0rte#2026"
	hash, err := p.Hash(plain)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	// deve validar com a senha correta
	require.True(t, p.Verify(hash, plain))

	// não deve validar com senha errada
	require.False(t, p.Verify(hash, "wrong-password"))
}

func TestPasswordHash_Hash_IsDeterministicFalse(t *testing.T) {
	// bcrypt usa salt, então hashes do mesmo plain devem ser diferentes
	p := &PasswordHash{}

	plain := "same-password"
	h1, err := p.Hash(plain)
	require.NoError(t, err)

	h2, err := p.Hash(plain)
	require.NoError(t, err)

	require.NotEqual(t, h1, h2)
	require.True(t, p.Verify(h1, plain))
	require.True(t, p.Verify(h2, plain))
}

func TestPasswordHash_Verify_InvalidHash(t *testing.T) {
	p := &PasswordHash{}

	// hash inválido deve retornar false (CompareHashAndPassword retorna erro)
	require.False(t, p.Verify("not-a-bcrypt-hash", "anything"))
}

func TestNewPasswordHash_ReturnsInterface(t *testing.T) {
	ph := NewPasswordHash()
	require.NotNil(t, ph)

	hash, err := ph.Hash("abc123!!")
	require.NoError(t, err)
	require.True(t, ph.Verify(hash, "abc123!!"))
}
