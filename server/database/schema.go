package database

type Tournament struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	TableCount int     `gorm:"DEFAULT:0" json:"tableCount"`
	State      int     `gorm:"DEFAULT:0" json:"state"`
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

// const schema = `
// 	create table tournament (
// 		id            integer primary key AUTOINCREMENT,
// 		name          text not null unique,
// 		table_count   integer default 0,
// 		state		  integer default 0
// 	);
// 	create table team (
// 		id            integer primary key AUTOINCREMENT,
// 		tournament_id integer not null,
// 		group_id	  integer,
// 		name          text not null,
// 		player_1      text not null,
// 		player_2      text not null,
// 		player_3      text not null
// 	);
// 	create table [group] (
// 		id            integer primary key AUTOINCREMENT,
// 		tournament_id integer not null,
// 		name	      text not null unique
// 	);
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
