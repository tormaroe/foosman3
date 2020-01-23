package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Init initializes a foosman3 database
func Init(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if new, err := isNew(db); err == nil && new {
		err = createSchema(db)
	}
	log.Println("Database initialized")
	return db, err
}

func isNew(db *sql.DB) (bool, error) {
	row := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='team'")
	var name string
	err := row.Scan(&name)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return true, nil
		default:
			return false, err
		}
	}
	return false, nil
}

func createSchema(db *sql.DB) error {
	log.Println("Creating database schema")
	_, err := db.Exec(schema)
	return err
}
