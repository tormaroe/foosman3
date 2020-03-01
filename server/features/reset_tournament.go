package features

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/database"
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

	// Check requirement
	var playedCount int
	if err := ac.DB.Model(database.Match{}).Where(
		"tournament_id = ? AND state = ?",
		tournamentID,
		int(core.Played),
	).Count(&playedCount).Error; err != nil {
		return err
	}

	if playedCount > 0 {
		log.Println("Can't reset tournament when some matches are played")
		return c.JSON(http.StatusConflict, map[string]string{
			"message": "Can't reset tournament when there are played matches.",
		})
	}

	err = ac.DB.Transaction(func(tx *gorm.DB) error {

		if err := deleteTournamentMatches(tx, tournamentID); err != nil {
			return err
		}

		if err := tx.Exec(
			"update teams set group_id = 0 where tournament_id = ?",
			tournamentID,
		).Error; err != nil {
			return err
		}

		if err := tx.Delete(
			database.Group{},
			"tournament_id = ?",
			tournamentID,
		).Error; err != nil {
			return err
		}

		err := tx.Table("tournaments").Where(
			"id = ?",
			tournamentID,
		).Update("state", int(core.New)).Error

		return err
	})

	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
