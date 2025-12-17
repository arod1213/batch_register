package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/arod1213/auto_ingestion/database"
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
	db, err := database.Setup()
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

	// SIMPLE CRUD
	r.GET("/songs", func(c *gin.Context) {
		handlers.FetchSongs(c, db)
	})
	r.DELETE("/songs", func(c *gin.Context) {
		handlers.DeleteSongs(c, db)
	})
	r.POST("/register/:isrc", func(c *gin.Context) {
		handlers.MarkRegistered(c, db)
	})

	// SPOTIFY CALLS
	r.GET("/read/:id", func(c *gin.Context) {
		handlers.FetchTracks(c, db)
	})

	// EXCEL GENERATION
	r.POST("/write", func(c *gin.Context) {
		handlers.WriteTracks(c, db)
	})

	err = http.ListenAndServe(":8080", r.Handler())
	if err != nil {
		fmt.Println("error setting up", err.Error())
		os.Exit(1)
	}
}
