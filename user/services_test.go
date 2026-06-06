package user

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// mockRepository lets each test drive repository behaviour through closures
// without touching a real database.
type mockRepository struct {
	saveFn        func(User) (User, error)
	findByEmailFn func(string) (User, error)
	findByIDFn    func(string) (User, error)
	updateFn      func(User) (User, error)
}

func (m *mockRepository) Save(u User) (User, error)          { return m.saveFn(u) }
func (m *mockRepository) FindByEmail(e string) (User, error) { return m.findByEmailFn(e) }
func (m *mockRepository) FindByID(id string) (User, error)   { return m.findByIDFn(id) }
func (m *mockRepository) Update(u User) (User, error)        { return m.updateFn(u) }

func TestRegisterUserHashesPasswordAndSetsRole(t *testing.T) {
	var saved User
	repo := &mockRepository{
		saveFn: func(u User) (User, error) {
			saved = u
			return u, nil
		},
	}
	service := NewService(repo)

	newUser, err := service.RegisterUser(RegisterUserInput{
		Name:     "Wiku",
		Email:    "wiku@example.com",
		Password: "supersecret",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, newUser.ID)
	assert.Equal(t, "user", saved.Role)
	assert.NotEqual(t, "supersecret", saved.PasswordHash)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(saved.PasswordHash), []byte("supersecret")))
}

func TestLoginSucceedsWithCorrectPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	repo := &mockRepository{
		findByEmailFn: func(string) (User, error) {
			return User{ID: "user-7", PasswordHash: string(hash)}, nil
		},
	}
	service := NewService(repo)

	loggedIn, err := service.Login(LoginInput{Email: "a@b.com", Password: "correct"})

	require.NoError(t, err)
	assert.Equal(t, "user-7", loggedIn.ID)
}

func TestLoginFailsWithWrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	repo := &mockRepository{
		findByEmailFn: func(string) (User, error) {
			return User{ID: "user-7", PasswordHash: string(hash)}, nil
		},
	}
	service := NewService(repo)

	_, err := service.Login(LoginInput{Email: "a@b.com", Password: "wrong"})

	assert.Error(t, err)
}

func TestLoginFailsWhenUserMissing(t *testing.T) {
	repo := &mockRepository{
		findByEmailFn: func(string) (User, error) { return User{}, nil },
	}
	service := NewService(repo)

	_, err := service.Login(LoginInput{Email: "ghost@b.com", Password: "x"})

	assert.EqualError(t, err, "user not found")
}

func TestIsEmailAvailable(t *testing.T) {
	tests := []struct {
		name      string
		found     User
		findErr   error
		available bool
		wantErr   bool
	}{
		{name: "available when no user", found: User{}, available: true},
		{name: "taken when user exists", found: User{ID: "user-3"}, available: false},
		{name: "propagates repo error", findErr: errors.New("db down"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{
				findByEmailFn: func(string) (User, error) { return tt.found, tt.findErr },
			}
			service := NewService(repo)

			available, err := service.IsEmailAvailable(CheckEmailInput{Email: "a@b.com"})

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.available, available)
		})
	}
}
