package core

type TournamentState int

const (
	New TournamentState = iota
	GroupsReady
	GroupPlayStarted
	EliminationPlayStarted
	Done
)
