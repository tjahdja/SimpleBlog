package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/tjahdja/SimpleBlog/internal/entity"
	"github.com/tjahdja/SimpleBlog/internal/repository"
)

var (
	ErrUnauthorized = errors.New("unauthorized action")
	ErrNotFound     = errors.New("resource not found")
)

type PostService interface {
	ListPosts(ctx context.Context) ([]*entity.Post, error)
	GetPostByID(ctx context.Context, id uint) (*entity.Post, error)
	CreatePost(ctx context.Context, title string, content string, userID uint) (*entity.Post, error)
	UpdatePost(ctx context.Context, id uint, title string, content string, userID uint) (*entity.Post, error)
	DeletePost(ctx context.Context, id uint, userID uint) error
}

type GORMPostService struct {
	repo repository.PostRepository
}

func NewGORMPostService(repo repository.PostRepository) PostService {
	return &GORMPostService{repo: repo}
}

func (s *GORMPostService) ListPosts(ctx context.Context) ([]*entity.Post, error) {
	posts, err := s.repo.ListPosts(ctx)
	if err != nil {
		slog.Error("Failed to list posts", "err", err)
		return nil, err
	}
	return posts, nil
}

func (s *GORMPostService) GetPostByID(ctx context.Context, id uint) (*entity.Post, error) {
	post, err := s.repo.GetPostByID(ctx, id)
	if err != nil {
		slog.Warn("Post not found by ID", "id", id, "err", err.Error())
		return nil, ErrNotFound
	}
	return post, nil
}

func (s *GORMPostService) CreatePost(ctx context.Context, title string, content string, userID uint) (*entity.Post, error) {
	post := &entity.Post{
		Title:    title,
		Content:  content,
		AuthorID: userID,
	}

	createdPost, err := s.repo.CreatePost(ctx, post)
	if err != nil {
		slog.Error("Failed to create post", "err", err)
		return nil, err
	}
	return createdPost, nil
}

func (s *GORMPostService) UpdatePost(ctx context.Context, id uint, title string, content string, userID uint) (*entity.Post, error) {
	existingPost, err := s.repo.GetPostByID(ctx, id)
	if err != nil {
		slog.Error("Failed to get post by ID", "id", id, "err", err)
		return nil, err
	}

	if existingPost.AuthorID != userID {
		slog.Error("User is not authorized to update this post")
		return nil, entity.ErrUnauthorized
	}

	existingPost.Title = title
	existingPost.Content = content

	updatedPost, err := s.repo.UpdatePost(ctx, existingPost)
	if err != nil {
		slog.Error("Failed to update post", "err", err)
		return nil, err
	}
	return updatedPost, nil
}

func (s *GORMPostService) DeletePost(ctx context.Context, id uint, userID uint) error {
	existingPost, err := s.repo.GetPostByID(ctx, id)
	if err != nil {
		slog.Error("Failed to get post by ID", "id", id, "err", err)
		return err
	}

	if existingPost.AuthorID != userID {
		slog.Error("User is not authorized to delete this post")
		return entity.ErrUnauthorized
	}

	if err := s.repo.DeletePost(ctx, existingPost.ID); err != nil {
		slog.Error("Failed to delete post", "err", err)
		return err
	}
	return nil
}
