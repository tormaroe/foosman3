package features

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/database"
)

// DeleteTeam deletes a team from a Tournament that hasn't started yet
func DeleteTeam(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	ID, err := ac.GetParamID()
	if err != nil {
		return err
	}

	var team database.Team
	if err := ac.DB.First(&team, ID).Error; err != nil {
		return err
	}

	if err := ac.AssertTournamentNotStarted(team.TournamentID); err != nil {
		return err
	}

	if err := ac.DB.Delete(&team).Error; err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
