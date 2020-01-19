package main

import (
	"fmt"
	"log"

	"github.com/golobby/config"
	"github.com/golobby/config/feeder"

	"github.com/tormaroe/foosman3/server/api"
	"github.com/tormaroe/foosman3/server/database"

	"github.com/labstack/echo"
)

func main() {
	c, err := config.New(config.Options{
		Feeder: feeder.Json{Path: "config.json"},
	})
	if err != nil {
		log.Fatal("Unable to read config:", err)
	}

	dbPath, _ := c.Get("database.path")
	log.Println("Database path:", dbPath)

	db, err := database.Init(fmt.Sprintf("%v", dbPath))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	e := echo.New()
	api.Init(e, db)

	address, _ := c.Get("api.address")
	log.Println("API address:", address)
	e.Logger.Fatal(e.Start(fmt.Sprintf("%v", address)))
}
