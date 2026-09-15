package server

import (
	"strings"

	"github.com/abubakar-sadiq001/learning-go-shop/internal/models"
	"github.com/abubakar-sadiq001/learning-go-shop/internal/utils"
	"github.com/gin-gonic/gin"
)

// authMiddleware is a Gin middleware function that checks for the presence of a valid JWT token in the Authorization header of incoming requests. If the token is valid, it extracts the user information from the token claims and sets it in the Gin context for further use in the request handling.
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization: Bearer JWT_TOKEN
		authHeader := c.GetHeader("authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Authorization header is missing", nil)
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenParts[1], s.config.JWT.Secret)
		if err != nil {
			utils.UnauthorizedResponse(c, "Invalid token", err)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

func (s *Server) adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			utils.ForbiddenResponse(c, "Forbidden", nil)
			c.Abort()
			return
		}

		if role != string(models.UserRoleAdmin) {
			utils.ForbiddenResponse(c, "Forbidden", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
