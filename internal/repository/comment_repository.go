package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/tjahdja/SimpleBlog/internal/entity"

	"gorm.io/gorm"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *entity.Comment) (*entity.Comment, error)
	GetCommentByID(ctx context.Context, id uint) (*entity.Comment, error)
	ListCommentsByPostID(ctx context.Context, postID uint) ([]*entity.Comment, error)
	DeleteComment(ctx context.Context, id uint) error
}

type GORMCommentRepository struct {
	db *gorm.DB
}

func NewGORMCommentRepository(db *gorm.DB) CommentRepository {
	return &GORMCommentRepository{db: db}
}

func (r *GORMCommentRepository) CreateComment(ctx context.Context, comment *entity.Comment) (*entity.Comment, error) {
	if err := r.db.WithContext(ctx).Create(comment).Error; err != nil {
		slog.Error("Failed to create comment", "err", err)
		return nil, err
	}
	return comment, nil
}

func (r *GORMCommentRepository) GetCommentByID(ctx context.Context, id uint) (*entity.Comment, error) {
	var comment entity.Comment
	if err := r.db.WithContext(ctx).First(&comment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("Post record not found", "id", id)
			return nil, err
		}
		slog.Error("Failed to get comment by ID", "id", id, "err", err)
		return nil, err
	}
	return &comment, nil
}

func (r *GORMCommentRepository) ListCommentsByPostID(ctx context.Context, postID uint) ([]*entity.Comment, error) {
	var comments []*entity.Comment
	if err := r.db.WithContext(ctx).Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		slog.Error("Failed to list comments by post ID", "postID", postID, "err", err)
		return nil, err
	}
	return comments, nil
}

func (r *GORMCommentRepository) DeleteComment(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Comment{}, id).Error; err != nil {
		slog.Error("Database error during comment deletion", "id", id, "err", err.Error())
		return err
	}
	return nil
}
