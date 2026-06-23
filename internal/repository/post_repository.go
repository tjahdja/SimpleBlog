package repository

import (
	"context"
	"errors"
	"log/slog"

	"SimpleBlog/internal/entity"

	"gorm.io/gorm"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *entity.Post) (*entity.Post, error)
	GetPostByID(ctx context.Context, id uint) (*entity.Post, error)
	ListPosts(ctx context.Context) ([]*entity.Post, error)
	UpdatePost(ctx context.Context, post *entity.Post) (*entity.Post, error)
	DeletePost(ctx context.Context, id uint) error
}

type GORMPostRepository struct {
	db *gorm.DB
}

func NewGORMPostRepository(db *gorm.DB) PostRepository {
	return &GORMPostRepository{db: db}
}

func (r *GORMPostRepository) CreatePost(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	if err := r.db.WithContext(ctx).Create(post).Error; err != nil {
		slog.Error("Failed to create post", "err", err)
		return nil, err
	}
	return post, nil
}

func (r *GORMPostRepository) GetPostByID(ctx context.Context, id uint) (*entity.Post, error) {
	var post entity.Post
	if err := r.db.WithContext(ctx).First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("Post record not found", "id", id)
			return nil, err
		}

		slog.Error("Failed to get post by ID", "id", id, "err", err)
		return nil, err
	}
	return &post, nil
}

func (r *GORMPostRepository) ListPosts(ctx context.Context) ([]*entity.Post, error) {
	var posts []*entity.Post
	if err := r.db.WithContext(ctx).Find(&posts).Error; err != nil {
		slog.Error("Failed to list posts", "err", err)
		return nil, err
	}
	return posts, nil
}

func (r *GORMPostRepository) UpdatePost(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	if err := r.db.WithContext(ctx).Save(post).Error; err != nil {
		slog.Error("Database error during post update", "id", post.ID, "err", err.Error())
		return nil, err
	}
	return post, nil
}

func (r *GORMPostRepository) DeletePost(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Post{}, id).Error; err != nil {
		slog.Error("Database error during post deletion", "id", id, "err", err.Error())
		return err
	}
	return nil
}
