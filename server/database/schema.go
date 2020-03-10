package database

import "errors"

import "github.com/tormaroe/foosman3/server/core"

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
	Group        Group
	Tournament   Tournament `json:"-"`
	TournamentID int        `json:"tournamentId"`
	Name         string     `json:"name"`
	Player1      string     `json:"player1"`
	Player2      string     `json:"player2"`
	Player3      string     `json:"player3"`
}

type Group struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	Tournament   Tournament `json:"-"`
	TournamentID int        `json:"tournamentId"`
}

type Match struct {
	ID                 int `json:"id"`
	Team1              Team
	Team1ID            int `json:"team1_id"`
	Team2              Team
	Team2ID            int `json:"team2_id"`
	TournamentID       int
	GroupID            int `json:"groupId"`
	Group              Group
	Table              string `json:"table"`
	State              int    `json:"state"` // TODO: use typed constant?
	Sequence           int
	PlayoffTier        int `json:"playoff_tier"`         // 1, 2, 4, 8, 16, ...
	PlayoffMatchNumber int `json:"playoff_match_number"` // Number within the tier
	MatchResults       []MatchResult
}

type MatchResult struct {
	ID      int
	TeamID  int
	MatchID int
	Points  int
	Win     int
	Loss    int
	Draw    int
}

// const schema = `
// 	create table log (
// 		id            integer primary key AUTOINCREMENT,
// 		created       integer(4) not null default (strftime('%s','now')),
// 		tournament_id integer not null,
// 		message       text not null
// 	);
// `

func (m Match) GetWinnerAndLooser() (Team, Team, bool, error) {
	var winner Team
	var looser Team
	if m.State != int(core.Played) {
		return winner, looser, false, errors.New("Match not played, can't get winner")
	}
	if len(m.MatchResults) < 2 {
		return winner, looser, false, errors.New("Missing match results, can't get winner")
	}
	if m.MatchResults[0].Win == 1 {
		if m.MatchResults[0].TeamID == m.Team1ID {
			winner = m.Team1
			looser = m.Team2
		} else {
			winner = m.Team2
			looser = m.Team1
		}
	} else if m.MatchResults[1].Win == 1 {
		if m.MatchResults[1].TeamID == m.Team1ID {
			winner = m.Team1
			looser = m.Team2
		} else {
			winner = m.Team2
			looser = m.Team1
		}
	} else {
		return winner, looser, true, nil // DRAW
	}
	return winner, looser, false, nil
}
