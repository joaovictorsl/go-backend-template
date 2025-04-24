package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaovictorsl/go-backend-template/internal/http/jwt"
)

func MustAuthenticate(jwtService jwt.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		claims, err := jwtService.ValidateAccessToken(cookie.Value)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("user_id", claims.UserId)
		c.Next()
	}
}
