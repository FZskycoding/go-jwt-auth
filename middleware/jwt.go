package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JWT 密鑰，實際應用中應該從環境變數獲取
const SecretKey = "your-256-bit-secret"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//獲取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供授權標頭"})
			c.Abort()
			return
		}

		// 檢查 Bearer token 格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "授權格式錯誤"})
			c.Abort()
			return
		}

		// 提取 token、驗證 token
		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的token"})
			c.Abort()
			return
		}

		// 從 token 中獲取使用者資訊並存入上下文
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
		}

		c.Next()
	}
}
