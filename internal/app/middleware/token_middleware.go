package middleware

import (
	utility "aTalkBackEnd/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")

		splitToken := strings.Split(token, "Bearer ")
		if len(splitToken) != 2 {
			// 错误处理
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header111"})
		}
		token = splitToken[1]

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not provided"})
			c.Abort()
			return
		}

		claims, err := utility.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		// 将用户信息存储在context中，以便后续的handlers可以使用
		c.Set("userName", claims.UserID)
		c.Next()
	}
}
