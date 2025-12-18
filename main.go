package main

import (
	"fmt"
	"log"
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
	r.GET("/songs", middleware.Auth(), func(c *gin.Context) {
		handlers.FetchSongs(c, db)
	})

	r.PUT("/share/:id", func(c *gin.Context) {
		handlers.SaveShare(c, db)
	})

	r.DELETE("/songs", func(c *gin.Context) {
		handlers.DeleteSongs(c, db)
	})
	r.POST("/register/:isrc", func(c *gin.Context) {
		handlers.MarkRegistered(c, db)
	})

	// SPOTIFY CALLS
	r.GET("/read/:id", middleware.Auth(), func(c *gin.Context) {
		log.Println("Router - All keys:", c.Keys)
		handlers.FetchTracks(c, db)
	})

	// EXCEL GENERATION
	r.POST("/write", func(c *gin.Context) {
		handlers.WriteShares(c, db)
	})

	err = http.ListenAndServe(":8080", r.Handler())
	if err != nil {
		fmt.Println("error setting up", err.Error())
		os.Exit(1)
	}
}
