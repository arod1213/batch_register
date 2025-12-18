package handlers

import (
	"fmt"
	"os"
	"time"

	"github.com/arod1213/auto_ingestion/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenInfo struct {
	token   string
	expires string
}

func newToken(user models.User) (string, error) {
	secret := os.Getenv("JWKS_KEY")

	if secret == "" {
		return "", fmt.Errorf("JWKS_KEY environment variable is not set")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("fucking error is", err)
		return "", err
	}
	return tokenString, nil
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var info LoginInfo
		err := c.ShouldBindJSON(&info)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "bad request"})
			return
		}
		var user models.User
		err = db.Where("username = ?", info.Username).Find(&user).Error
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "bad request"})
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(info.Password))
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "bad request"})
			return
		}
		token, err := newToken(user)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "could not sign token"})
			return
		}
		c.JSON(200, gin.H{"token": token, "user_id": user.ID, "username": user.Username})
	}
}

func Signup(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var info LoginInfo
		err := c.ShouldBindJSON(&info)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "bad request"})
			return
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(info.Password), 12)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "bad request"})
			return
		}

		fmt.Println("username is ", info.Username)
		user := models.User{
			Username: info.Username,
			Password: string(hashed),
		}
		tx := db.Begin()
		err = tx.Create(&user).Error
		if err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(400, gin.H{"error": "failed to signup"})
			return
		}

		token, err := newToken(user)
		if err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(500, gin.H{"error": "could not sign token"})
			return
		}

		tx.Commit()
		c.JSON(200, gin.H{"token": token, "user_id": user.ID, "username": user.Username})
	}
}
