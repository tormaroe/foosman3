package core

type TournamentState int

const (
	New TournamentState = iota
	GroupPlayStarted
	EliminationPlayStarted
	Done
)
