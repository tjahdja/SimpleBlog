package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/tjahdja/SimpleBlog/internal/entity"
	"github.com/tjahdja/SimpleBlog/internal/handler"
	"github.com/tjahdja/SimpleBlog/internal/service"
)

type mockUserRepository struct {
	users         map[string]*entity.User
	mockError     error
	lastSavedUser *entity.User
	idCounter     uint
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:     make(map[string]*entity.User),
		idCounter: 1,
	}
}

func (m *mockUserRepository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}

	if _, exists := m.users[user.Username]; exists {
		return nil, errors.New("username already exists")
	}

	user.ID = m.idCounter
	m.idCounter++

	m.users[user.Username] = user
	m.lastSavedUser = user

	return user, nil
}

func (m *mockUserRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}
	user, exists := m.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *mockUserRepository) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func TestRegister_Success(t *testing.T) {
	repo := newMockUserRepository()
	jwtSecret := "test_secret_key"
	userService := service.NewGORMUserService(repo, jwtSecret)

	input := handler.RegisterRequest{
		Username: "testuser",
		Password: "password123",
	}

	_, err := userService.Register(context.Background(), input.Username, input.Password)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.lastSavedUser == nil {
		t.Fatal("expected user to be saved to repository, but it wasn't")
	}

	if repo.lastSavedUser.Username != input.Username {
		t.Errorf("expected saved username to be %s, got %s", input.Username, repo.lastSavedUser.Username)
	}

	if repo.lastSavedUser.Password == input.Password {
		t.Error("expected password to be securely hashed, but it remained plain text")
	}
}

func TestLogin_Success(t *testing.T) {
	// Arrange
	repo := newMockUserRepository()
	jwtSecret := "test_secret_key"
	userService := service.NewGORMUserService(repo, jwtSecret)

	// Pre-populate mock database with a user (manually hash a password or use registration first)
	registerInput := handler.RegisterRequest{
		Username: "loginuser",
		Password: "mysecurepassword",
	}
	_, _ = userService.Register(context.Background(), registerInput.Username, registerInput.Password)

	loginInput := handler.LoginRequest{
		Username: "loginuser",
		Password: "mysecurepassword",
	}

	token, err := userService.Login(context.Background(), loginInput.Username, loginInput.Password)

	if err != nil {
		t.Fatalf("expected successful login, got error: %v", err)
	}

	if token == "" {
		t.Error("expected a generated JWT string token payload, got empty string")
	}
}
