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

	return response, nil
}
