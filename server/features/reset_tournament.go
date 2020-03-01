package features

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
)

// ResetTournament will remove all matches and groups,
// but will keep all teams of the tournament. Teams can then be added or removed,
// and groups can be created, even though the tournament was previously started.
// Requirements: No match results have been registered.
func ResetTournament(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	tournamentID, err := ac.GetParamID()
	if err != nil {
		return err
	}

	ac.DB.Transaction(func(tx *gorm.DB) error {

		if err := deleteTournamentMatches(tx, tournamentID); err != nil {
			return err
		}

		// TODO: Delete groups and remove group foreign keys

		// TODO: Reset tournament.state

		return nil
	})

	return c.NoContent(http.StatusOK)
}
