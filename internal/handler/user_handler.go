package handler

import (
	"net/http"

	"github.com/tjahdja/SimpleBlog/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register godoc
// @Summary      Register a new user account
// @Description  Create a brand new credentials profile for the blog infrastructure
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterRequest  true  "User Registration Payload"
// @Success      201      {object}  map[string]interface{} "User successfully registered"
// @Failure      400      {object}  map[string]string      "Invalid json input payload or mapping error"
// @Router       /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Register(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// Login godoc
// @Summary      Authenticate user and return JWT
// @Description  Verify credentials and return a bearer token for protected endpoints
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "User Credentials"
// @Success      200      {object}  map[string]string         "Returns access token"
// @Failure      401      {object}  map[string]string         "Invalid credentials"
// @Router       /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.userService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
