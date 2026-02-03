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

	r := gin.New()
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
	r.GET("/user", middleware.Auth(), handlers.GetMe(db))                 // user profile info
	r.PUT("/user/identify", middleware.Auth(), handlers.IdentifyUser(db)) // user profile info
	r.PUT("/user/update/:id", middleware.Auth(), handlers.UpdateUser(db)) // update user profile

	r.GET("/song/:songID", middleware.Auth(), handlers.GetSong(db))   // song overview for user
	r.GET("/songs", middleware.Auth(), handlers.FetchSongs(db))       // fetch user songs
	r.POST("/songs", middleware.Auth(), handlers.SaveTracks(db))      // create songs
	r.DELETE("/shares", middleware.Auth(), handlers.DeleteShares(db)) // provide list of song ids in body
	r.PUT("/share/:id", handlers.SaveShare(db))                       // update share for user
	// r.POST("/register/:isrc", handlers.MarkRegistered(db))            // mark registered

	r.POST("/payments/scan", middleware.Auth(), handlers.RescanPayments(db))     // insert new payments
	r.POST("/payments", middleware.Auth(), handlers.SaveRoyalties(db))           // insert new payments
	r.GET("/payments/song/:songID", middleware.Auth(), handlers.GetPayments(db)) // read payments for song by user

	// EXCEL GENERATION
	r.POST("/download", middleware.Auth(), handlers.DownloadRegistrations(db)) // download shares as excel
	r.POST("/download/all", middleware.Auth(), handlers.DownloadAllShares(db)) // download all shares as excel

	// SPOTIFY CALLS
	r.GET("/read/preview/:id", handlers.FetchTracks())                             // preview spotify catalog
	r.GET("/read/save/:id", middleware.Auth(), handlers.FetchAndSaveTracks2(db))   // import spotify catalog
	r.GET("/read/saveOld/:id", middleware.Auth(), handlers.FetchAndSaveTracks(db)) // import spotify catalog

	// GENIUS CALLS
	r.GET("/genius/artist", middleware.Auth(), handlers.GetMissingSongs(db)) // param for genius artist id
	r.GET("/genius/search", handlers.GeniusSearch(db))                       // query param for keyword
	r.GET("/genius/search/artist", handlers.GeniusSearchArtistIDs(db))       // query param for keyword

	err = http.ListenAndServe(":8080", r.Handler())
	if err != nil {
		fmt.Println("error setting up", err.Error())
		os.Exit(1)
	}
}
