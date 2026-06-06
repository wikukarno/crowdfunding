package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend-crowdfunding/auth"
	"backend-crowdfunding/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// assertAnError is a sentinel used by handler tests that need the service layer
// to fail without caring about the specific error.
var assertAnError = errors.New("service failure")

func protectedRouter(authService auth.Service, userService user.Service) *gin.Engine {
	router := gin.New()
	router.GET("/me", AuthMiddleware(authService, userService), func(c *gin.Context) {
		current := c.MustGet("currentUser").(user.User)
		c.JSON(http.StatusOK, gin.H{"id": current.ID})
	})
	return router
}

func TestAuthMiddlewareRejectsMissingToken(t *testing.T) {
	router := protectedRouter(auth.NewService("secret"), &fakeUserService{})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddlewareRejectsInvalidToken(t *testing.T) {
	router := protectedRouter(auth.NewService("secret"), &fakeUserService{})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer garbage.token.value")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddlewareAcceptsValidToken(t *testing.T) {
	authService := auth.NewService("secret")
	token, err := authService.GenerateToken("user-42")
	require.NoError(t, err)

	userService := &fakeUserService{
		getByIDFn: func(id string) (user.User, error) {
			return user.User{ID: id}, nil
		},
	}
	router := protectedRouter(authService, userService)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"id":"user-42"`)
}
