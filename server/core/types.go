package core

// Team represents a foosball team of two players
type Team struct {
	ID      int
	Name    string
	Player1 string
	Player2 string
	GroupID *int // ? replace with name ?
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
