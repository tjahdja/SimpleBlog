package repository

import (
	"context"
	"errors"
	"log/slog"

	"SimpleBlog/internal/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	GetUserByID(ctx context.Context, id uint) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
}

type GORMUserRepository struct {
	db *gorm.DB
}

func NewGORMUserRepository(db *gorm.DB) UserRepository {
	return &GORMUserRepository{db: db}
}

func (r *GORMUserRepository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		slog.Error("Failed to create user", "err", err.Error())
		return nil, err
	}
	return user, nil
}

func (r *GORMUserRepository) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("User record not found", "id", id)
			return nil, err
		}

		slog.Error("Database error when fetching user by ID", "id", id, "err", err.Error())
		return nil, err
	}
	return &user, nil
}

func (r *GORMUserRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("User record not found", "username", username)
			return nil, err
		}

		slog.Error("Database error when fetching user by username", "username", username, "err", err.Error())
		return nil, err
	}
	return &user, nil
}
