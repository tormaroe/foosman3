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
	_, err := d.db.Exec(schema)
	return err
}

// Close the database connection
func (d *Database) Close() {
	d.db.Close()
}

// ----------------------------------------------------------------------------
//
//  INSERTS, UPDATES, and DELETES
// ----------------------------------------------------------------------------

// AddTournament saves a new Tournament entity
func (d *Database) AddTournament(t core.Tournament) error {
	stmt, err := d.db.Prepare(`
		insert into tournament
		(name, table_count, state)
		values
		(?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(t.Name, t.TableCount, core.New)
	return err
}

// UpdateTournament saves changes to a Tournament entity
func (d *Database) UpdateTournament(t core.Tournament) error {
	stmt, err := d.db.Prepare(`
		update tournament
		set name=?, table_count=?, state=?
		where id=?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(t.Name, t.TableCount, t.State, t.ID)
	return err
}

// AddTeam saves a new Team entity
func (d *Database) AddTeam(tournamentID int, t core.Team) error {
	stmt, err := d.db.Prepare(`
		insert into team
		(name, tournament_id, player_1, player_2, player_3) 
		values
		(?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(t.Name, tournamentID, t.Player1, t.Player2, t.Player3)
	return err
}

// UpdateTeam saves changes to a Team entity
func (d *Database) UpdateTeam(t core.Team) error {
	stmt, err := d.db.Prepare(`
		update team 
		set name=?, player_1=?, player_2=?, player_3=? 
		where id=?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(t.Name, t.Player1, t.Player2, t.Player3, t.ID)
	return err
}

// ----------------------------------------------------------------------------
//
//  QUERIES
// ----------------------------------------------------------------------------

// GetTournaments gets all the tournament entities
func (d *Database) GetTournaments() ([]core.Tournament, error) {
	rows, err := d.db.Query(`
		select id, name, table_count, state
		from tournament
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var result []core.Tournament
	for rows.Next() {
		var t core.Tournament
		err = rows.Scan(&t.ID, &t.Name, &t.TableCount, &t.State)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, rows.Err()
}

// GetTournamentTeams gets all the teams for a Tournament from the database
func (d *Database) GetTournamentTeams(ID int) ([]core.Team, error) {
	rows, err := d.db.Query(`
		select id, name, player_1, player_2, player_3
		from team 
		where tournament_id=?
	`, ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var result []core.Team
	for rows.Next() {
		var t core.Team
		err = rows.Scan(&t.ID, &t.Name, &t.Player1, &t.Player2, &t.Player3)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, rows.Err()
}
