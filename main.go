package main

import (
	"github.com/syvita/ronin/api"
	"github.com/syvita/ronin/events"

	"github.com/syvita/ronin/db"

	"log"
	"os"
)

var RedisAddr = "localhost:6379"

var logger = log.New(os.Stderr, "[MAIN]: ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	database, err := db.NewDatabase(RedisAddr)

	if err != nil {
		logger.Fatalf("Failed to connect to redis: %s", err.Error())
	}

	logger.Println("Connected to Redis successfully")

	events.Database = database
	api.Database = database

	// goroutines >>>>> anything else

	go api.Listen(":3999")
	go events.Listen(":3700")

	for {
	}
}
