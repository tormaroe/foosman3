package features

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/database"
)

type getTeamResponse struct {
	Team    database.Team
	Matches []database.Match
}

func GetTeam(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	teamID, err := ac.GetParamID()
	if err != nil {
		return err
	}
	var team database.Team
	if team, err = getTeam(ac.DB, teamID); err != nil {
		return err
	}
	var matches []database.Match
	if matches, err = getTeamMatches(ac.DB, teamID); err != nil {
		return err
	}
	response := getTeamResponse{
		Team:    team,
		Matches: matches,
	}
	return c.JSONPretty(http.StatusOK, response, "  ")
}

func getTeam(db *gorm.DB, ID int) (database.Team, error) {
	var response database.Team
	err := db.Preload("Group").Find(&response, ID).Error
	return response, err
}

func getTeamMatches(db *gorm.DB, ID int) ([]database.Match, error) {
	var response []database.Match
	err := db.Preload("Team1").Preload("Team2").Preload("MatchResults").Where("team1_id = ? OR team2_id = ?", ID, ID).Find(&response).Error
	return response, err
}
