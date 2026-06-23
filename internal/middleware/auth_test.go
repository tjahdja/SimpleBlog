package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tjahdja/SimpleBlog/internal/middleware"
)

func TestAuthMiddleware_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	r.Use(middleware.AuthMiddleware("test_secret"))
	r.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	ctx.Request, _ = http.NewRequest(http.MethodGet, "/protected", nil)
	r.ServeHTTP(w, ctx.Request)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401 Unauthorized for missing token, got %d", w.Code)
	}
}
