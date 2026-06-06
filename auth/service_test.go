package auth

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndValidateToken(t *testing.T) {
	service := NewService("test-secret")

	token, err := service.GenerateToken("user-42")
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	parsed, err := service.ValidateToken(token)
	require.NoError(t, err)
	assert.True(t, parsed.Valid)

	claims, ok := parsed.Claims.(jwt.MapClaims)
	require.True(t, ok)
	assert.Equal(t, "user-42", claims["user_id"])
}

func TestValidateTokenRejectsGarbage(t *testing.T) {
	service := NewService("test-secret")

	_, err := service.ValidateToken("not-a-real-token")
	assert.Error(t, err)
}

func TestValidateTokenRejectsWrongSecret(t *testing.T) {
	token, err := NewService("secret-a").GenerateToken("user-1")
	require.NoError(t, err)

	parsed, err := NewService("secret-b").ValidateToken(token)
	assert.Error(t, err)
	assert.False(t, parsed.Valid)
}
