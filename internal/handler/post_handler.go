package handler

import (
	"errors"
	"net/http"
	"strconv"

	_ "github.com/tjahdja/SimpleBlog/internal/entity"
	"github.com/tjahdja/SimpleBlog/internal/service"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService service.PostService
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

// PostCreateRequest defines the expected input for creating a post
type PostCreateRequest struct {
	Title   string `json:"title" binding:"required,min=3,max=100"`
	Content string `json:"content" binding:"required"`
}

type PostUpdateRequest struct {
	Title   string `json:"title" binding:"required,min=3,max=100"`
	Content string `json:"content" binding:"required"`
}

// ListPosts godoc
// @Summary      Get all blog posts
// @Description  Retrieve a timeline list of all published blog entries
// @Tags         Posts
// @Produce      json
// @Success      200      {array}   entity.Post        "List of posts"
// @Failure      500      {object}  map[string]string  "Internal server error"
// @Router       /posts [get]
func (h *PostHandler) ListPosts(c *gin.Context) {
	posts, err := h.postService.ListPosts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// GetPostByID godoc
// @Summary      Get a single blog post
// @Description  Retrieve full details of a specific blog post using its ID
// @Tags         Posts
// @Produce      json
// @Param        id       path      int                true  "Post ID"
// @Success      200      {object}  entity.Post        "Post details"
// @Failure      404      {object}  map[string]string  "Post not found"
// @Router       /posts/{id} [get]
func (h *PostHandler) GetPostByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}

	post, err := h.postService.GetPostByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// CreatePost godoc
// @Summary      Create a new blog entry
// @Description  Publish a brand new story post to the platform feed
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      PostCreateRequest  true  "Post Content Details"
// @Success      201      {object} 	entity.Post              "Post created successfully"
// @Failure      401      {object}  map[string]string        "Unauthorized token context error"
// @Router       /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req PostCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract the authenticated user ID set by the JWT Middleware
	userID := c.MustGet("userID").(uint)

	createdPost, err := h.postService.CreatePost(c.Request.Context(), req.Title, req.Content, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, createdPost)
}

// UpdatePost godoc
// @Summary      Update an existing post
// @Description  Modify the title or content of a blog entry by ID
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                     true  "Post ID"
// @Param        request  body      PostUpdateRequest 		true  "Updated Content"
// @Success      200      {object}  entity.Post             "Updated post details"
// @Failure      401      {object}  map[string]string       "Unauthorized"
// @Failure      404      {object}  map[string]string       "Post not found"
// @Router       /posts/{id} [put]
// @Link         entity.Post github.com/tjahdja/SimpleBlog/internal/entity.Post
func (h *PostHandler) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}

	var req PostUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(uint)

	updatedPost, err := h.postService.UpdatePost(c.Request.Context(), uint(id), req.Title, req.Content, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to update this post"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, updatedPost)
}

// DeletePost godoc
// @Summary      Delete a blog post
// @Description  Remove a blog entry completely from the database by ID
// @Tags         Posts
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                true  "Post ID"
// @Success      200      {object}  map[string]string  "Post deleted successfully"
// @Failure      401      {object}  map[string]string  "Unauthorized"
// @Failure      404      {object}  map[string]string  "Post not found"
// @Router       /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}

	userID := c.MustGet("userID").(uint)

	err = h.postService.DeletePost(c.Request.Context(), uint(id), userID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to delete this post"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
