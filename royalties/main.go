package royalties

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/arod1213/auto_ingestion/models"
	"github.com/arod1213/auto_ingestion/spotify"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Payment struct {
	ID        uint       `gorm:"primaryKey;autoIncrement"`
	Hash      string     `gorm:"type:text;not null"`
	Earnings  float64    `gorm:"type:real;not null"`
	Payor     string     `gorm:"type:text;not null"`
	Date      *time.Time `gorm:"-"`
	Territory *string    `gorm:"type:text"`

	SongID *uint       `gorm:"type:varchar(15)"`
	Song   models.Song `gorm:"foreignKey:SongID;references:ID;constraint:OnDelete:CASCADE;"`

	UserID uint        `gorm:"type:varchar(15);not null"`
	User   models.User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`
}

type payor struct {
	Name string `json:"name"`
}

type ExtPayment struct {
	// id        string `json:"id"`
	// hash        string `json:"hash"`
	Earnings float64 `json:"earnings"`
	Payor    payor   `json:"payor"`
	Artist   *string `json:"artist"`
	Title    string  `json:"title"`
	Isrc     *string `json:"isrc"`
	Iswc     *string `json:"iswc"`
	Upc      *uint64 `json:"upc"`
	// Date      *time.Time `json:"date"`
	Territory *string `json:"territory"`
}

func (p ExtPayment) IntoPayment(db *gorm.DB, userID uint) (*Payment, error) {
	var song models.Song

	if p.Isrc != nil {
		err := db.Where("isrc = ?", *p.Isrc).Find(&song).Error
		if err != nil {
			return nil, err
		}
		fmt.Println("FOUND ON ISRC")
	} else if p.Iswc != nil {
		err := db.Where("iswc = ?", *p.Iswc).Find(&song).Error
		if err != nil {
			return nil, err
		}
		fmt.Println("FOUND ON ISWC")
	} else {
		query := db.Joins("LEFT JOIN shares on shares.song_id = songs.id").Where("songs.title LIKE ?", "%"+p.Title+"%")
		if p.Artist != nil {
			query = query.Where("songs.artist LIKE ?", "%"+*p.Artist+"%")
		}
		err := query.Where("shares.user_id = ?", userID).Find(&song).Error
		if err != nil {
			return nil, err
		}
		fmt.Println("FOUND ON NAME")
	}

	payment := Payment{
		SongID:   &song.ID,
		Earnings: p.Earnings,
		Payor:    p.Payor.Name,
		// Date:      p.Date,
		Territory: p.Territory,
	}
	return &payment, nil
}

func ReadCsv(filePath string) error {
	base := os.Getenv("API_URL_ROYALTY")
	url := fmt.Sprintf("%s/read/payment", base)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	writer.Close()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var payment ExtPayment
	err = json.Unmarshal(respBody, &payment)
	if err != nil {
		return err
	}
	spotify.Pretty(payment)
	return nil
}

func ForwardFiles(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			c.String(400, "expected multipart/form-data with files")
			return
		}

		var userID uint = 1
		base := os.Getenv("API_URL_ROYALTY")
		url := fmt.Sprintf("%s/read/payment", base)

		req, err := http.NewRequest("POST", url, c.Request.Body)
		if err != nil {
			c.String(400, "failed to create request")
		}
		req.Header.Set("Content-Type", contentType)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.String(400, "failed to send request")
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		var payment []ExtPayment
		err = json.Unmarshal(respBody, &payment)
		if err != nil {
			c.String(400, "failed to read payments")
		}

		var i uint = 0
		var p []Payment = make([]Payment, len(payment))
		for _, pay := range payment {
			x, err := pay.IntoPayment(db, userID)
			if err != nil {
				continue
			}
			p[i] = *x
			i++
		}

		tx := db.Begin()
		err = tx.Create(&p).Error
		if err != nil {
			tx.Rollback()
			c.String(400, "failed to create payments")
		}
		tx.Commit()
	}
}
