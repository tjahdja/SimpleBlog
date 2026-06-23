package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/tjahdja/SimpleBlog/internal/entity"
	"github.com/tjahdja/SimpleBlog/internal/service"
)

type CommentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

type CommentCreateRequest struct {
	Content string `json:"content" binding:"required,min=1,max=500"`
}

// CreateComment godoc
// @Summary      Add a comment to a post
// @Description  Publish a new textual reply under an existing blog entry
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                        true  "Post ID"
// @Param        request  body      CommentCreateRequest 	   true  "Comment Body"
// @Success      201      {object}  entity.Comment             "Comment added successfully"
// @Failure      401      {object}  map[string]string          "Unauthorized"
// @Router       /posts/{id}/comments [post]
// @Link         entity.Comment github.com/tjahdja/SimpleBlog/internal/entity.Comment
func (h *CommentHandler) CreateComment(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}

	var req CommentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(uint)

	createdComment, err := h.commentService.CreateComment(c.Request.Context(), req.Content, userID, uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to post comment"})
		return
	}

	c.JSON(http.StatusCreated, createdComment)
}

// GetCommentsByPostID godoc
// @Summary      Get comments for a post
// @Description  Retrieve all user feedback comments associated with a specific post ID
// @Tags         Comments
// @Produce      json
// @Param        id       path      int                  true  "Post ID"
// @Success      200      {array}   entity.Comment       "List of comments"
// @Router       /posts/{id}/comments [get]
func (h *CommentHandler) GetCommentsByPostID(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}

	comments, err := h.commentService.GetCommentsByPostID(c.Request.Context(), uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// DeleteComment godoc
// @Summary      Delete a comment
// @Description  Remove a comment statement by its specific unique comment ID
// @Tags         Comments
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                true  "Comment ID"
// @Success      200      {object}  map[string]string  "Comment removed successfully"
// @Failure      401      {object}  map[string]string  "Unauthorized"
// @Router       /comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID format"})
		return
	}

	userID := c.MustGet("userID").(uint)

	err = h.commentService.DeleteComment(c.Request.Context(), uint(commentID), userID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to delete this comment"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
