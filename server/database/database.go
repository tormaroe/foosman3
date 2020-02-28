package database

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/tormaroe/foosman3/server/core"

	// Importing SQLite driver
	_ "github.com/mattn/go-sqlite3"
)

// Init initializes a foosman3 database
func Init(path string) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Tournament{}, &Team{}, &Group{}, &Match{}, &MatchResult{}, &User{})
	log.Println("Database initialized")
	return db, err
}

// AssertTournamentNotStarted will return an error if tournament has
// started.
func AssertTournamentNotStarted(ac *core.FoosmanContext, ID int) error {
	var tournament Tournament
	if err := ac.DB.First(&tournament, ID).Error; err != nil {
		return err
	}

	if tournament.State != int(core.New) {
		return errors.New("Can't delete team from a tournament that has started")
	}

	return nil
}
