package features

import (
	"errors"
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

	var tournament database.Tournament
	if err := ac.DB.First(&tournament, team.TournamentID).Error; err != nil {
		return err
	}

	if tournament.State != int(core.New) {
		return errors.New("Can't delete team from a tournament that has started")
	}

	if err := ac.DB.Delete(&team).Error; err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
