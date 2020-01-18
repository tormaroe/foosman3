package main

import (
	"fmt"
	"log"

	"github.com/golobby/config"
	"github.com/golobby/config/feeder"

	"github.com/tormaroe/foosman3/server/database"
)

func main() {
	c, err := config.New(config.Options{
		Feeder: feeder.Json{Path: "config.json"},
	})
	if err != nil {
		log.Fatal("Unable to read config:", err)
	}

	dbPath, err := c.Get("database.path")
	log.Println("Database path:", dbPath)

	db, err := database.Init(fmt.Sprintf("%v", dbPath))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
