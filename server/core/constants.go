package core

type TournamentState int

const (
	New TournamentState = iota
	GroupPlayStarted
	EliminationPlayStarted
	Done
)

type MatchState int

const (
	Planned MatchState = iota
	Scheduled
	InProgress
	Played
)
