package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			fmt.Println("Invalid Header")
			c.AbortWithStatusJSON(404, gin.H{"error": "invalid token"})
			return
		}

		jwksKey := os.Getenv("JWKS_KEY")
		if jwksKey == "" {
			fmt.Println("Invalid KEY")
			c.AbortWithStatusJSON(404, gin.H{"error": "invalid token"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		key := []byte(jwksKey)

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return key, nil
		})

		if err != nil || !token.Valid {
			fmt.Println("failed to get valid token", err.Error())
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("failed to get token claims")
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			fmt.Println("failed to get user_id")
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		userID := uint(userIDFloat)
		c.Set("userID", userID)
		c.Next()
	}
}
