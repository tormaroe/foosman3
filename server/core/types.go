package core

type TournamentState int

const (
	New TournamentState = iota
	GroupsReady
	GroupPlayStarted
	EliminationPlayStarted
	Done
)

type Tournament struct {
	ID         int             `json:"id"`
	Name       string          `json:"name"`
	TableCount int             `json:"tableCount"`
	State      TournamentState `json:"state"`
}

// Team represents a foosball team of two players
type Team struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Player1 string `json:"player1"`
	Player2 string `json:"player2"`
	Player3 string `json:"player3"`
}

// Group represents a set of teams playing eachother in the groups stage
type Group struct {
	ID      int
	Name    string
	Teams   []Team
	Matches []Match
}

type Match struct {
	ID    int
	Team1 Team
	Team2 Team
}
