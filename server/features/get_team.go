package features

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
)

type getTeamResponse struct {
	Name        string                 `json:"name"`
	GroupName   string                 `json:"groupName"`
	GroupID     int                    `json:"groupId"`
	Player1     string                 `json:"player1"`
	Player2     string                 `json:"player2"`
	Player3     string                 `json:"player3"`
	TotalScore  int                    `json:"totalScore"`
	GamesPlayed int                    `json:"gamesPlayed"`
	GamesWon    int                    `json:"gamesWon"`
	GamesLost   int                    `json:"gamesLost"`
	GamesDraw   int                    `json:"gamesDraw"`
	Matches     []getTeamMatchResponse `json:"matches"`
}

type getTeamMatchResponse struct {
	OpponentID   int    `json:"opponentId"`
	OpponentName string `json:"opponentName"`
	State        int    `json:"state"`
	Sequence     int    `json:"sequence"`
	Win          int    `json:"win"`
	Loss         int    `json:"loss"`
	Draw         int    `json:"draw"`
	Points       int    `json:"points"`
}

func GetTeam(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	teamID, err := ac.GetParamID()
	if err != nil {
		return err
	}
	var response getTeamResponse
	if response, err = getTeam(ac.DB, teamID); err != nil {
		return err
	}
	return c.JSONPretty(http.StatusOK, response, "  ")
}

func getTeam(db *gorm.DB, ID int) (getTeamResponse, error) {
	const query = `
		select 
		  t.name, t.player1, t.player2, t.player3,
		  g.id as group_id, g.name as group_name
		from teams t
		join groups g on g.id = t.group_id
		where t.id = ?
	`
	var response getTeamResponse
	if err := db.Raw(query, ID).Scan(&response).Error; err != nil {
		return response, err
	}

	const matchQuery = `
		select
		  m.state, 
		  m.sequence,
		  case when m.team1_id = ? then m.team2_id else m.team1_id end as opponent_id,
		  case when m.team1_id = ? then t2.name else t1.name end as opponent_name,
		  r.win, r.loss, r.draw, r.points
		from matches m
		join teams t1 on t1.id = m.team1_id
		join teams t2 on t2.id = m.team2_id
		join match_results r on r.match_id = m.id and r.team_id <> ?
		where m.team1_id = ? or m.team2_id = ?
		order by sequence desc
	`
	rows, err := db.Raw(matchQuery, ID, ID, ID, ID, ID).Rows()
	if err != nil {
		return response, err
	}
	for rows.Next() {
		var m getTeamMatchResponse
		if err := db.ScanRows(rows, &m); err != nil {
			return response, err
		}
		response.Matches = append(response.Matches, m)
	}

	return response, nil
}
