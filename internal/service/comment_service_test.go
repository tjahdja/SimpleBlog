package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/tjahdja/SimpleBlog/internal/entity"
	"github.com/tjahdja/SimpleBlog/internal/service"
)

type mockCommentRepository struct {
	comments  map[uint]*entity.Comment
	mockError error
	idCounter uint
}

func newMockCommentRepository() *mockCommentRepository {
	return &mockCommentRepository{
		comments:  make(map[uint]*entity.Comment),
		idCounter: 1,
	}
}

func (m *mockCommentRepository) CreateComment(ctx context.Context, comment *entity.Comment) (*entity.Comment, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}

	comment.ID = m.idCounter
	comment.CreatedAt = time.Now()
	m.idCounter++

	m.comments[comment.ID] = comment
	return comment, nil
}

func (m *mockCommentRepository) ListCommentsByPostID(ctx context.Context, postID uint) ([]*entity.Comment, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}
	var result []*entity.Comment
	for _, comment := range m.comments {
		if comment.PostID == postID {
			commentCopy := *comment
			result = append(result, &commentCopy)
		}
	}
	return result, nil
}

func (m *mockCommentRepository) GetCommentByID(ctx context.Context, id uint) (*entity.Comment, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}
	comment, exists := m.comments[id]
	if !exists {
		return nil, errors.New("comment not found")
	}
	commentCopy := *comment
	return &commentCopy, nil
}

func (m *mockCommentRepository) DeleteComment(ctx context.Context, id uint) error {
	if m.mockError != nil {
		return m.mockError
	}
	if _, exists := m.comments[id]; !exists {
		return errors.New("comment not found")
	}
	delete(m.comments, id)
	return nil
}

func TestCreateComment_Success(t *testing.T) {
	repo := newMockCommentRepository()
	commentService := service.NewGORMCommentService(repo)
	content := "This is a great blog post!"
	postID := uint(1)
	userID := uint(20)
	comment, err := commentService.CreateComment(context.Background(), content, postID, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if comment.Content != content {
		t.Errorf("expected content %q, got %q", content, comment.Content)
	}
	if comment.PostID != postID {
		t.Errorf("expected post ID %d, got %d", postID, comment.PostID)
	}
	if comment.UserID != userID {
		t.Errorf("expected user ID %d, got %d", userID, comment.UserID)
	}
}

func TestDeleteComment_UnauthorizedUser(t *testing.T) {
	repo := newMockCommentRepository()
	commentService := service.NewGORMCommentService(repo)

	comment, _ := repo.CreateComment(context.Background(), &entity.Comment{
		Content: "Nice post!",
		PostID:  1,
		UserID:  10,
	})

	hackerUserID := uint(99)
	err := commentService.DeleteComment(context.Background(), comment.ID, hackerUserID)
	if err == nil {
		t.Fatal("expected an authorization error when a non-owner tries to delete a comment, but got nil")
	}

	_, findErr := repo.GetCommentByID(context.Background(), comment.ID)
	if findErr != nil {
		t.Error("security breach: comment was deleted by an unauthorized user account")
	}
}
