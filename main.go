package main

import (
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
	"gitlab.com/syvita/ronin/db"

	//	"github.com/gin-gonic/autotls"
	"log"
	"net/http"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
	RedisAddr = "localhost:6379"
)
func main() {
	database, err := db.NewDatabase(RedisAddr)

	if err != nil {
		log.Fatalf("Failed to connect to redis: %s", err.Error())
	} else {
		fmt.Println("Connected to Redis successfully")
	}

	apiServer := &http.Server{
		Addr:         ":3999",
		Handler:      initApiRouter(database),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	eventServer := &http.Server{
		Addr:         ":3700",
		Handler:      initEventRouter(database),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		return apiServer.ListenAndServe()
	})

	g.Go(func() error {
		return eventServer.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func initApiRouter(database *db.Database) *gin.Engine {
	r := gin.Default()
	r.GET("/points/:username", func (c *gin.Context) {
		username := c.Param("username")
		user, err := database.GetUser(username)
		if err != nil {
			if err == db.ErrNil {
				c.JSON(http.StatusNotFound, gin.H{"error": "No record found for " + username})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user})
	})

	r.POST("/points", func (c *gin.Context) {
		var userJson db.User
		if err := c.ShouldBindJSON(&userJson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := database.SaveUser(&userJson)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": userJson})
	})

	r.GET("/leaderboard", func(c *gin.Context) {
		leaderboard, err := database.GetLeaderboard()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"leaderboard": leaderboard})
	})

	// we can add auto-TLS setup here on the API server boot if we want
	// guessing this should be userconfigurable in some way or another
	//log.Fatal(autotls.Run(r, "example1.com", "example2.com"))

	return r
}

func initEventRouter(database *db.Database) *gin.Engine {
	r := gin.Default()
	r.GET("/points/:username", func (c *gin.Context) {
		username := c.Param("username")
		user, err := database.GetUser(username)
		if err != nil {
			if err == db.ErrNil {
				c.JSON(http.StatusNotFound, gin.H{"error": "No record found for " + username})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user})
	})

	r.POST("/points", func (c *gin.Context) {
		var userJson db.User
		if err := c.ShouldBindJSON(&userJson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := database.SaveUser(&userJson)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": userJson})
	})

	r.GET("/leaderboard", func(c *gin.Context) {
		leaderboard, err := database.GetLeaderboard()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"leaderboard": leaderboard})
	})

	//log.Fatal(autotls.Run(r, "example1.com", "example2.com"))

	return r
}
