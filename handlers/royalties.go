package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/arod1213/auto_ingestion/middleware"
	"github.com/arod1213/auto_ingestion/royalties"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RescanPayments(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := middleware.GetUserID(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		err = royalties.Reconcile(db, userID)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to save"})
			return
		}
		c.JSON(200, gin.H{"data": "success"})
	}
}

func GetPayments(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := middleware.GetUserID(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		idStr := c.Param("songID")
		songID, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(400, gin.H{"error": "bad request"})
			return
		}

		var payments []royalties.Payment
		err = db.Where("song_id = ? AND user_id = ?", uint(songID), userID).Find(&payments).Error
		if err != nil {
			c.JSON(500, gin.H{"error": "payments not found"})
			return
		}

		c.JSON(200, gin.H{"data": payments})
	}
}

func SaveRoyalties(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			c.String(400, "expected multipart/form-data with files")
			return
		}

		userID, err := middleware.GetUserID(c)
		if err != nil {
			c.String(401, "unauthorized")
			return
		}

		base := os.Getenv("API_URL_ROYALTY")
		url := fmt.Sprintf("%s/read/payment", base)

		req, err := http.NewRequest("POST", url, c.Request.Body)
		if err != nil {
			c.String(500, fmt.Sprintf("failed to create request: %v", err))
			return
		}

		req.Header.Set("Content-Type", contentType)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.String(500, fmt.Sprintf("failed to send request: %v", err))
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var payments []royalties.ExtPayment
		err = json.Unmarshal(body, &payments)
		if err != nil {
			c.String(500, fmt.Sprintf("failed to unmarshal json: %v", err))
			return
		}

		id, err := royalties.SavePayments(db, userID, payments)
		if err != nil {
			c.String(500, fmt.Sprintln("failed to save payments"))
			return
		}
		c.JSON(200, gin.H{"data": id})
	}
}
