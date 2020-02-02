package core

type TournamentState int

const (
	New TournamentState = iota
	MatchesPlanned
	GroupPlayStarted
	EliminationPlayStarted
	Done
)

type MatchState int

const (
	Planned MatchState = iota
	Scheduled
	Played
)
