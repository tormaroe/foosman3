package core

type TournamentState int

const (
	New TournamentState = iota
	GroupPlayStarted
	GroupPlayDone
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
