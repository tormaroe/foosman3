package database

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

// Init initializes a foosman3 database
func Init(path string) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Tournament{}, &Team{}, &Group{}, &Match{})
	log.Println("Database initialized")
	return db, err
}
