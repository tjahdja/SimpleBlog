package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/tjahdja/SimpleBlog/internal/entity"
	"github.com/tjahdja/SimpleBlog/internal/repository"
)

type UserService interface {
	Register(ctx context.Context, username string, password string) (*entity.User, error)
	Login(ctx context.Context, username string, password string) (string, error)
}

type GORMUserService struct {
	repo      repository.UserRepository
	jwtSecret []byte
}

func NewGORMUserService(repo repository.UserRepository, jwtSecret string) UserService {
	return &GORMUserService{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *GORMUserService) Register(ctx context.Context, username string, password string) (*entity.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Failed to hash password", "err", err)
		return nil, err
	}

	user := &entity.User{
		Username: username,
		Password: string(hashedPassword),
	}

	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		slog.Error("Failed to create user", "err", err)
		return nil, err
	}

	return createdUser, nil
}

func (s *GORMUserService) Login(ctx context.Context, username string, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		slog.Error("Failed to get user by username", "username", username, "err", err)
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		slog.Error("Password mismatch", "err", err)
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours
	})

	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		slog.Error("Failed to sign JWT", "err", err.Error())
		return "", err
	}

	return signedToken, nil
}
