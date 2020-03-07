package features

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/database"
)

func GetEliminationMatches(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	tournamentID, err := ac.GetParamID()
	if err != nil {
		return err
	}

	var matches []database.Match
	if err := ac.DB.Preload("Team1").Preload("Team2").Preload("MatchResults").Where(
		"tournament_id = ? and playoff_tier > 0",
		tournamentID,
	).Order("playoff_tier, playoff_match_number").Find(&matches).Error; err != nil {
		return err
	}

	return c.JSONPretty(http.StatusOK, matches, "  ")
}
