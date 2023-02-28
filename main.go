package main

import (
	"log"
	"time"
	"os"

	"github.com/gin-contrib/cache/persistence"
	"github.com/pocketbase/pocketbase"
)

var app *pocketbase.PocketBase

func main() {
	// Make sure to compile all keywords
	compileAllKeywords()

	// Initialize interactive backend
	app = pocketbase.New()

	store := persistence.NewInMemoryStore(time.Second)

	router := newAPI(store)

	if router == nil {
		log.Panic("Can't create API server")
	}

	// Start scapper in a separate routine
	go webScrapper()

	// Load port from environment
	apiPort := os.Getenv("PUBLIC_API_PORT")

	// Choose default port if not set
	if apiPort == "" {
		apiPort = "8091"
	}

	// Start HTTP API in a separate routine
	go router.Run(":" + apiPort)

	// Use the main thread to run the backend
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
