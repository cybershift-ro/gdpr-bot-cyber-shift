package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
)

// Configure API
func newAPI(cache_storage *persistence.InMemoryStore) *gin.Engine {
	if cache_storage == nil {
		return nil
	}

	router := gin.Default()

	router.GET("/rss", cache.CachePage(cache_storage, 5*time.Minute, serveRSSFeed))
	router.GET("/json", cache.CachePage(cache_storage, 5*time.Minute, serveJSONFeed))
	router.GET("/atom", cache.CachePage(cache_storage, 5*time.Minute, serveAtomFeed))

	router.GET("/sanctions", cache.CachePage(cache_storage, 5*time.Minute, serveJSONSanctions))

	return router
}

func serveRSSFeed(c *gin.Context) {
	feedMutex.Lock()
	rss, err := sanctionFeed.ToRss()
	feedMutex.Unlock()

	if err != nil {
		c.AbortWithError(http.StatusServiceUnavailable, err)
	}

	// Set the right type
	c.Writer.Header().Set("Content-Type", "application/rss+xml")

	c.String(http.StatusOK, rss)
}

func serveJSONFeed(c *gin.Context) {
	feedMutex.Lock()
	json, err := sanctionFeed.ToJSON()
	feedMutex.Unlock()

	if err != nil {
		c.AbortWithError(http.StatusServiceUnavailable, err)
	}

	// Set the right type
	c.Writer.Header().Set("Content-Type", "application/json")

	c.String(http.StatusOK, json)
}

func serveAtomFeed(c *gin.Context) {
	feedMutex.Lock()
	atom, err := sanctionFeed.ToAtom()
	feedMutex.Unlock()

	if err != nil {
		c.AbortWithError(http.StatusServiceUnavailable, err)
	}

	// Set the right type
	c.Writer.Header().Set("Content-Type", "application/atom+xml")

	c.String(http.StatusOK, atom)
}

func serveJSONSanctions(c *gin.Context) {
	if app == nil {
		c.AbortWithError(http.StatusServiceUnavailable, errors.New("app not ready"))
		return
	}

	sanctions_feed := sanctionsToJSON()

	// Set the right type
	c.Writer.Header().Set("Content-Type", "application/json")

	c.String(http.StatusOK, sanctions_feed)
}
