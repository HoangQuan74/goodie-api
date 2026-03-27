package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/HoangQuan74/goodie-api/pkg/auth"
	apperrors "github.com/HoangQuan74/goodie-api/pkg/errors"
	"github.com/HoangQuan74/goodie-api/pkg/response"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	ContextUserID       = "user_id"
	ContextUserRole     = "user_role"
)

func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(AuthorizationHeader)
		if header == "" {
			response.Error(c, apperrors.Unauthorized("missing authorization header"))
			c.Abort()
			return
		}

		if !strings.HasPrefix(header, BearerPrefix) {
			response.Error(c, apperrors.Unauthorized("invalid authorization format"))
			c.Abort()
			return
		}

		token := strings.TrimPrefix(header, BearerPrefix)
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUserRole, claims.Role)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get(ContextUserRole)
		if !exists {
			response.Error(c, apperrors.Unauthorized("missing user role"))
			c.Abort()
			return
		}

		role := userRole.(string)
		for _, r := range roles {
			if r == role {
				c.Next()
				return
			}
		}

		response.Error(c, apperrors.Forbidden("insufficient permissions"))
		c.Abort()
	}
}
