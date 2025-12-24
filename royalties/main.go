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
	"time"

	"github.com/arod1213/auto_ingestion/models"
	"github.com/arod1213/auto_ingestion/spotify"
)

type Payment struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Hash string `gorm:"type:text;not null;index:idx_hash_share,unique"`

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
	Hash     string  `json:"hash"`
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
