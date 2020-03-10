package features

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/database"
)

type gTMResult struct {
	Team1     string
	Team2     string
	GroupName string
}

func GetTournamentMatches(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	tournamentID, err := ac.GetParamID()
	if err != nil {
		return err
	}
	matches, err := getTournamentMatches(ac.DB, tournamentID)
	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusOK, matches, "  ")
}

func getTournamentMatches(db *gorm.DB, tournamentID int) ([]database.Match, error) {
	var matches []database.Match
	err := db.Preload("MatchResults").Preload("Team1").Preload("Team2").Preload("Group").Where("tournament_id = ?", tournamentID).Order("sequence desc").Find(&matches).Error
	return matches, err
}

func getPlayedTournamentMatches(db *gorm.DB, tournamentID int) ([]database.Match, error) {
	var matches []database.Match
	err := db.Preload("MatchResults").Preload("Team1").Preload("Team2").Preload("Group").Where("tournament_id = ? and state = ?", tournamentID, int(core.Played)).Order("sequence desc").Find(&matches).Error
	return matches, err
}
