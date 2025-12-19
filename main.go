package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/arod1213/auto_ingestion/database"
	"github.com/arod1213/auto_ingestion/handlers"
	"github.com/arod1213/auto_ingestion/middleware"
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

	// AUTH
	r.POST("/signup", handlers.Signup(db))
	r.POST("/login", handlers.Login(db))

	// SIMPLE CRUD
	r.GET("/songs", middleware.Auth(), handlers.FetchSongs(db))
	r.POST("/songs", middleware.Auth(), handlers.SaveTracks(db)) // create songs

	r.GET("/user", middleware.Auth(), handlers.GetMe(db))
	r.PUT("/user/:id", handlers.UpdateUser(db))

	// r.DELETE("/songs", handlers.DeleteSongs(db))
	r.PUT("/share/:id", handlers.SaveShare(db))
	r.POST("/register/:isrc", handlers.MarkRegistered(db))

	// EXCEL GENERATION
	r.POST("/download", middleware.Auth(), handlers.DownloadRegistrations(db))

	// SPOTIFY CALLS
	r.GET("/read/:id", middleware.Auth(), handlers.FetchAndSaveTracks(db))
	r.GET("/read/preview/:id", handlers.FetchTracks())

	err = http.ListenAndServe(":8080", r.Handler())
	if err != nil {
		fmt.Println("error setting up", err.Error())
		os.Exit(1)
	}
}
