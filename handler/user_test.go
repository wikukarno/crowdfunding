package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend-crowdfunding/auth"
	"backend-crowdfunding/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// fakeUserService implements user.Service. Each test wires up only the methods
// it needs; unused ones stay nil and would panic loudly if called by mistake.
type fakeUserService struct {
	registerFn func(user.RegisterUserInput) (user.User, error)
	loginFn    func(user.LoginInput) (user.User, error)
	emailFn    func(user.CheckEmailInput) (bool, error)
	avatarFn   func(string, string) (user.User, error)
	getByIDFn  func(string) (user.User, error)
}

func (f *fakeUserService) RegisterUser(in user.RegisterUserInput) (user.User, error) {
	return f.registerFn(in)
}
func (f *fakeUserService) Login(in user.LoginInput) (user.User, error) { return f.loginFn(in) }
func (f *fakeUserService) IsEmailAvailable(in user.CheckEmailInput) (bool, error) {
	return f.emailFn(in)
}
func (f *fakeUserService) SaveAvatar(id string, loc string) (user.User, error) {
	return f.avatarFn(id, loc)
}
func (f *fakeUserService) GetUserByID(id string) (user.User, error) { return f.getByIDFn(id) }

func performRequest(r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

func TestRegisterUserReturnsToken(t *testing.T) {
	userService := &fakeUserService{
		registerFn: func(user.RegisterUserInput) (user.User, error) {
			return user.User{ID: "user-1", Name: "Wiku", Email: "wiku@example.com"}, nil
		},
	}
	h := NewUserHandler(userService, auth.NewService("secret"), nil)

	router := gin.New()
	router.POST("/users", h.RegisterUser)

	rec := performRequest(router, http.MethodPost, "/users",
		`{"name":"Wiku","occupation":"engineer","email":"wiku@example.com","password":"secret123"}`)

	require.Equal(t, http.StatusOK, rec.Code)

	var body struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.NotEmpty(t, body.Data.Token)
}

func TestRegisterUserRejectsInvalidPayload(t *testing.T) {
	h := NewUserHandler(&fakeUserService{}, auth.NewService("secret"), nil)

	router := gin.New()
	router.POST("/users", h.RegisterUser)

	rec := performRequest(router, http.MethodPost, "/users", `{"email":"not-an-email"}`)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestLoginHidesCredentialError(t *testing.T) {
	userService := &fakeUserService{
		loginFn: func(user.LoginInput) (user.User, error) {
			return user.User{}, assertAnError
		},
	}
	h := NewUserHandler(userService, auth.NewService("secret"), nil)

	router := gin.New()
	router.POST("/sessions", h.Login)

	rec := performRequest(router, http.MethodPost, "/sessions",
		`{"email":"wiku@example.com","password":"wrong"}`)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid email or password")
	assert.NotContains(t, rec.Body.String(), assertAnError.Error())
}
