package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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

	var dev bool = false
	st := os.Getenv("IS_DEV")
	if st == "TRUE" {
		dev = true
		log.Println("runnign in dev mode")
	} else {
		log.Println("running in prod mode", st)
	}

	// AUTH
	r.POST("/signup", handlers.Signup(db))
	r.POST("/login", handlers.Login(db))

	// SIMPLE CRUD
	r.GET("/user", middleware.Auth(dev), handlers.GetMe(db))                 // user profile info
	r.PUT("/user/identify", middleware.Auth(dev), handlers.IdentifyUser(db)) // user profile info
	r.PUT("/user/update/:id", middleware.Auth(dev), handlers.UpdateUser(db)) // update user profile

	r.GET("/song/:songID", middleware.Auth(dev), handlers.GetSong(db))   // song overview for user
	r.GET("/songs", middleware.Auth(dev), handlers.FetchSongs(db))       // fetch user songs
	r.POST("/songs", middleware.Auth(dev), handlers.SaveTracks(db))      // create songs
	r.DELETE("/shares", middleware.Auth(dev), handlers.DeleteShares(db)) // provide list of song ids in body
	r.PUT("/share/:id", handlers.SaveShare(db))                          // update share for user
	// r.POST("/register/:isrc", handlers.MarkRegistered(db))            // mark registered

	// DEALS
	r.POST("deals/pub/:songID", middleware.Auth(dev), handlers.CreateDeals(db, false))
	r.POST("deals/master/:songID", middleware.Auth(dev), handlers.CreateDeals(db, true))
	r.DELETE("deal/pub/:songID/:dealID", middleware.Auth(dev), handlers.DeleteDeal(db, false))
	r.DELETE("deal/master/:songID/:dealID", middleware.Auth(dev), handlers.DeleteDeal(db, true))

	r.POST("/payments/scan", middleware.Auth(dev), handlers.RescanPayments(db))     // insert new payments
	r.POST("/payments", middleware.Auth(dev), handlers.SaveRoyalties(db))           // insert new payments
	r.GET("/payments/song/:songID", middleware.Auth(dev), handlers.GetPayments(db)) // read payments for song by user

	// STATEMENTS
	r.GET("/statement/:statementID", middleware.Auth(dev), handlers.FetchStatement(db)) // get overview for statement
	r.GET("/statements", middleware.Auth(dev), handlers.GetStatements(db))              // get user statements
	r.GET("/statements/:userID", func(c *gin.Context) {
		id := c.Param("userID")
		userID, _ := strconv.ParseUint(id, 10, 32)
		c.Set("userID", uint(userID))
	}, handlers.GetStatements(db)) // get user statements NO AUTH

	// EXCEL GENERATION
	r.POST("/download", middleware.Auth(dev), handlers.DownloadRegistrations(db)) // download shares as excel
	r.POST("/download/all", middleware.Auth(dev), handlers.DownloadAllShares(db)) // download all shares as excel

	// SPOTIFY CALLS
	r.GET("/read/preview/:id", handlers.FetchTracks())                                // preview spotify catalog
	r.GET("/read/save/:id", middleware.Auth(dev), handlers.FetchAndSaveTracks2(db))   // import spotify catalog
	r.GET("/read/saveOld/:id", middleware.Auth(dev), handlers.FetchAndSaveTracks(db)) // import spotify catalog

	r.GET("/user/search", middleware.Auth(dev), handlers.SearchUsers(db)) // search public users query param name

	// GENIUS CALLS
	r.GET("/genius/artist", middleware.Auth(dev), handlers.GetMissingSongs(db)) // param for genius artist id
	r.GET("/genius/search", handlers.GeniusSearch(db))                          // query param for keyword
	r.GET("/genius/search/artist", handlers.GeniusSearchArtistIDs(db))          // query param for keyword

	err = http.ListenAndServe(":8080", r.Handler())
	if err != nil {
		fmt.Println("error setting up", err.Error())
		os.Exit(1)
	}
}
