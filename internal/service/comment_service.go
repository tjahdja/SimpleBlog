package service

import (
	"context"
	"log/slog"

	"github.com/tjahdja/SimpleBlog/internal/entity"
	"github.com/tjahdja/SimpleBlog/internal/repository"
)

type CommentService interface {
	CreateComment(ctx context.Context, content string, userID uint, postID uint) (*entity.Comment, error)
	GetCommentsByPostID(ctx context.Context, postID uint) ([]*entity.Comment, error)
	DeleteComment(ctx context.Context, id uint, userID uint) error
	GetCommentByID(ctx context.Context, id uint) (*entity.Comment, error)
}

type GORMCommentService struct {
	repo repository.CommentRepository
}

func NewGORMCommentService(repo repository.CommentRepository) CommentService {
	return &GORMCommentService{repo: repo}
}

func (s *GORMCommentService) CreateComment(ctx context.Context, content string, userID uint, postID uint) (*entity.Comment, error) {
	comment := &entity.Comment{
		Content: content,
		UserID:  userID,
		PostID:  postID,
	}

	createdComment, err := s.repo.CreateComment(ctx, comment)
	if err != nil {
		slog.Error("Database failure creating comment", "err", err.Error())
		return nil, err
	}
	return createdComment, nil
}

func (s *GORMCommentService) GetCommentsByPostID(ctx context.Context, postID uint) ([]*entity.Comment, error) {
	comments, err := s.repo.ListCommentsByPostID(ctx, postID)
	if err != nil {
		slog.Error("Database failure fetching comments for post", "post_id", postID, "err", err.Error())
		return nil, err
	}
	return comments, nil
}

func (s *GORMCommentService) DeleteComment(ctx context.Context, id uint, userID uint) error {
	existingComment, err := s.repo.GetCommentByID(ctx, id)
	if err != nil {
		slog.Warn("Attempted to delete non-existent comment", "id", id)
		return ErrNotFound
	}

	// Safety check: Ensure the user trying to delete is the actual owner
	if existingComment.UserID != userID {
		slog.Warn("Unauthorized comment deletion attempt", "comment_id", id, "user_id", userID)
		return ErrUnauthorized
	}

	if err := s.repo.DeleteComment(ctx, id); err != nil {
		slog.Error("Database failure deleting comment", "id", id, "err", err.Error())
		return err
	}
	return nil
}

func (s *GORMCommentService) GetCommentByID(ctx context.Context, id uint) (*entity.Comment, error) {
	comment, err := s.repo.GetCommentByID(ctx, id)
	if err != nil {
		slog.Warn("Comment not found by ID", "id", id, "err", err.Error())
		return nil, ErrNotFound
	}
	return comment, nil
}
