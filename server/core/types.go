package core

// Team represents a foosball team of two players
type Team struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Player1 string `json:"player1"`
	Player2 string `json:"player2"`
	GroupID *int   // ? replace with name ?
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
