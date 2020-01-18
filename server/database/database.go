package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tormaroe/foosman3/server/core"
)

// Database holds a connection to a foosman3 database
// and provides methods to operate on it
type Database struct {
	db *sql.DB
}

// Init initializes a foosman3 database
func Init(path string) (*Database, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	d := &Database{
		db: db,
	}
	if new, err := d.isNew(); err == nil && new {
		err = d.createSchema()
	}
	log.Println("Database initialized")
	return d, err
}

func (d *Database) isNew() (bool, error) {
	row := d.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='team'")
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

func (d *Database) createSchema() error {
	log.Println("Creating database schema")

	createStmt := `
		create table team (
			id       integer primary key AUTOINCREMENT,
			name     text not null unique,
			player_1 text not null,
			player_2 text not null,
			group_id integer
		);
		create table group (
			id       integer primary key AUTOINCREMENT,
			name	 text not null unique
		);
		create table match (
			id       integer primary key AUTOINCREMENT,
			created  integer(4) not null default (strftime('%s','now')),
			team_1   integer not null,
			team_2   integer not null,
			group_id integer,
			state    integer default 0 
		);
		create table result (
			id       integer primary key AUTOINCREMENT,
			created  integer(4) not null default (strftime('%s','now')),
			team_id  integer not null,
			match_id integer not null,
			points   integer not null,
			win      integer not null,
		);
		create table log (
			id       integer primary key AUTOINCREMENT,
			created  integer(4) not null default (strftime('%s','now')),
			message  text not null
		);
	`
	_, err := d.db.Exec(createStmt)
	return err
}

// Close the database connection
func (d *Database) Close() {
	d.db.Close()
}

// SaveTeam updates a team. If Team.ID is 0, team is inserted.
func (d *Database) SaveTeam(team core.Team) error {
	if team.ID == 0 {
		stmt, err := d.db.Prepare("insert into team(name, player_1, player_2) values(?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(team.Name, team.Player1, team.Player2)
		return err
	}
	stmt, err := d.db.Prepare("update team set name=?, player_1=?, player_2=? where id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(team.Name, team.Player1, team.Player2, team.ID)
	return err
}

// AllTeams gets all the teams from the database
func (d *Database) AllTeams() ([]core.Team, error) {
	rows, err := d.db.Query("select id, name, player_1, player_2 from team")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var result []core.Team
	for rows.Next() {
		var t core.Team
		err = rows.Scan(&t.ID, &t.Name, &t.Player1, &t.Player2)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, rows.Err()
}
