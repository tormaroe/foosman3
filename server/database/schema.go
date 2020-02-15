package database

type Tournament struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	TableCount int     `gorm:"DEFAULT:0" json:"tableCount"`
	State      int     `gorm:"DEFAULT:0" json:"state"` // TODO: use typed constant?
	Teams      []Team  `json:"teams"`
	Groups     []Group `json:"groups"`
}

type Team struct {
	ID           int `json:"id"`
	GroupID      int `json:"groupId"`
	Tournament   Tournament
	TournamentID int    `json:"tournamentId"`
	Name         string `json:"name"`
	Player1      string `json:"player1"`
	Player2      string `json:"player2"`
	Player3      string `json:"player3"`
}

type Group struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Tournament   Tournament
	TournamentID int `json:"tournamentId"`
}

type Match struct {
	ID      int `json:"id"`
	Team1   Team
	Team1ID int `json:"team1_id"`
	Team2   Team
	Team2ID int `json:"team2_id"`
	// ??? Add tournament ID
	GroupID int    `json:"groupId"`
	Table   string `json:"table"`
	State   int    `json:"state"` // TODO: use typed constant?
}

type MatchResult struct {
	ID      int
	TeamID  int
	MatchID int
	Points  int
	Win     bool
	Loss    bool
	Draw    bool
}

// const schema = `
// 	create table match (
// 		id            integer primary key AUTOINCREMENT,
// 		created       integer(4) not null default (strftime('%s','now')),
// 		team_1        integer not null,
// 		team_2        integer not null,
// 		group_id      integer,
// 		table_name	  text,
// 		state         integer default 0
// 	);
// 	create table result (
// 		id            integer primary key AUTOINCREMENT,
// 		created       integer(4) not null default (strftime('%s','now')),
// 		team_id       integer not null,
// 		match_id      integer not null,
// 		points        integer not null,
// 		win           integer not null,
// 	);
// 	create table log (
// 		id            integer primary key AUTOINCREMENT,
// 		created       integer(4) not null default (strftime('%s','now')),
// 		tournament_id integer not null,
// 		message       text not null
// 	);
// `
