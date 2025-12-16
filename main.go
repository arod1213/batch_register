package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/arod1213/auto_ingestion/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("err is ", err)
		os.Exit(1)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://localhost:5000",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/read/:playlistID", handlers.FetchTracks)
	r.POST("/write", handlers.WriteTracks)

	err = http.ListenAndServe(":8080", r.Handler())
	if err != nil {
		os.Exit(1)
	}
}
