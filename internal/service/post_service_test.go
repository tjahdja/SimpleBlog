package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/tjahdja/SimpleBlog/internal/entity"
	"github.com/tjahdja/SimpleBlog/internal/handler"
	"github.com/tjahdja/SimpleBlog/internal/service"
)

type mockPostRepository struct {
	posts     map[uint]*entity.Post
	mockError error
	idCounter uint
}

func newMockPostRepository() *mockPostRepository {
	return &mockPostRepository{
		posts:     make(map[uint]*entity.Post),
		idCounter: 1,
	}
}

func (m *mockPostRepository) CreatePost(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}

	post.ID = m.idCounter
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	m.idCounter++

	m.posts[post.ID] = post
	return post, nil
}

func (m *mockPostRepository) GetPostByID(ctx context.Context, id uint) (*entity.Post, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}
	post, exists := m.posts[id]
	if !exists {
		return nil, errors.New("post not found")
	}
	// Return a copy to avoid unintended shared memory mutations in tests
	postCopy := *post
	return &postCopy, nil
}

func (m *mockPostRepository) ListPosts(ctx context.Context) ([]*entity.Post, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}
	var list []*entity.Post
	for _, post := range m.posts {
		postCopy := *post
		list = append(list, &postCopy)
	}
	return list, nil
}

func (m *mockPostRepository) UpdatePost(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}
	if _, exists := m.posts[post.ID]; !exists {
		return nil, errors.New("post not found")
	}
	post.UpdatedAt = time.Now()
	m.posts[post.ID] = post
	return post, nil
}

func (m *mockPostRepository) DeletePost(ctx context.Context, id uint) error {
	if m.mockError != nil {
		return m.mockError
	}
	if _, exists := m.posts[id]; !exists {
		return errors.New("post not found")
	}
	delete(m.posts, id)
	return nil
}

func TestCreatePost_Success(t *testing.T) {
	repo := newMockPostRepository()
	postService := service.NewGORMPostService(repo)

	input := handler.PostCreateRequest{
		Title:   "My First Blog Entry",
		Content: "Writing software tests is satisfying.",
	}
	userID := uint(42)

	post, err := postService.CreatePost(context.Background(), input.Title, input.Content, userID)

	if err != nil {
		t.Fatalf("expected post creation to succeed, got error: %v", err)
	}

	if post.ID == 0 {
		t.Error("expected generated post ID to be populated, got 0")
	}

	if post.AuthorID != userID {
		t.Errorf("expected author UserID to be %d, got %d", userID, post.AuthorID)
	}

	if post.Title != input.Title {
		t.Errorf("expected title to be %q, got %q", input.Title, post.Title)
	}
}

func TestGetPostByID_NotFound(t *testing.T) {
	repo := newMockPostRepository()
	postService := service.NewGORMPostService(repo)

	_, err := postService.GetPostByID(context.Background(), 999) // ID 999 does not exist

	if err == nil {
		t.Fatal("expected an error when looking up a missing post, got nil")
	}
}

func TestUpdatePost_UnauthorizedUser(t *testing.T) {
	repo := newMockPostRepository()
	postService := service.NewGORMPostService(repo)

	originalPost, _ := repo.CreatePost(context.Background(), &entity.Post{
		Title:    "Original Title",
		Content:  "Original Content",
		AuthorID: 5,
	})

	hackerUserID := uint(99)

	_, err := postService.UpdatePost(context.Background(), originalPost.ID, "Hacked Title", "Hacked Content", hackerUserID)

	if err == nil {
		t.Fatal("expected an error when an unauthorized user attempts to update a post, but got nil")
	}

	dbPost, _ := repo.GetPostByID(context.Background(), originalPost.ID)
	if dbPost.Title != "Original Title" {
		t.Error("security breach: post title was modified by an unauthorized user account")
	}
}

func TestDeletePost_SuccessByOwner(t *testing.T) {
	// Arrange
	repo := newMockPostRepository()
	postService := service.NewGORMPostService(repo)

	ownerID := uint(5)
	post, _ := repo.CreatePost(context.Background(), &entity.Post{
		Title:    "To Be Deleted",
		AuthorID: ownerID,
	})

	// Act
	err := postService.DeletePost(context.Background(), post.ID, ownerID)

	// Assert
	if err != nil {
		t.Fatalf("expected owner to delete post successfully, got error: %v", err)
	}
}
