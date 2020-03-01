package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/golobby/config"
	"github.com/golobby/config/feeder"

	"github.com/tormaroe/foosman3/server/api"
	"github.com/tormaroe/foosman3/server/database"

	"github.com/labstack/echo"
)

func main() {
	configFile := os.Getenv("FOOSMAN3_CONFIG")
	if len(configFile) == 0 {
		configFile = "config.json"
	}
	log.Printf("Reading configuration from: %s", configFile)

	c, err := config.New(config.Options{
		Feeder: feeder.Json{Path: configFile},
	})
	if err != nil {
		log.Fatal("Unable to read config:", err)
	}

	dbPath, _ := c.Get("database.path")
	address, _ := c.Get("api.address")
	logMode, _ := c.Get("instrumentation.dbLogMode")

	log.Println("Database path:", dbPath)
	log.Println("API address:", address)
	log.Println("DB log mode:", logMode)

	db, err := database.Init(fmt.Sprintf("%v", dbPath))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.LogMode(logMode.(bool))

	scheduleChan := database.NewScheduleChan()
	startNextMatchChan := database.NewStartMatchChan()

	e := echo.New()
	api.Init(e, db, scheduleChan, startNextMatchChan, new(sync.Mutex))

	e.Logger.Fatal(e.Start(fmt.Sprintf("%v", address)))
}
