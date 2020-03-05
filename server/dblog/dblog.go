package dblog

import "time"

// Needed?
type DBLogger interface {
	Log(tournamentID int, message string)
}

type Message struct {
	ID           int
	Timestamp    time.Time
	TournamentID int
	Text         string
}
