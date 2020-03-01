package features

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/database"
)

// ResetMatch will clear the results of a played match
// and set it to Planned state.
// If the are no matches scheduled in the Tournament,
// the match will be set to InProgress state. Else it
// will become scheduled.
func ResetMatch(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	matchID, err := ac.GetParamID()
	if err != nil {
		return err
	}

	ac.DB.Transaction(func(tx *gorm.DB) error {

		var results []database.MatchResult
		log.Printf("Reset match %d requested", matchID)
		if err := tx.Where("match_id = ?", matchID).Find(&results).Error; err != nil {
			log.Printf("Failed getting match results")
			return err
		}

		results[0].Points = 0
		results[0].Win = 0
		results[0].Loss = 0
		results[0].Draw = 0

		results[1].Points = 0
		results[1].Win = 0
		results[1].Loss = 0
		results[1].Draw = 0

		if err := tx.Save(&results[0]).Error; err != nil {
			log.Printf("Failed saving match results")
			return err
		}
		if err := tx.Save(&results[1]).Error; err != nil {
			log.Printf("Failed saving match results 2")
			return err
		}

		var match database.Match
		if err := tx.Find(&match, matchID).Error; err != nil {
			log.Printf("Failed finding match")
			return err
		}

		var inProgressCount int
		if err := tx.Model(database.Match{}).Where(
			"tournament_id = ? AND state = ?",
			match.TournamentID,
			int(core.InProgress),
		).Count(&inProgressCount).Error; err != nil {
			log.Printf("Failed counting matches in progress")
			return err
		}

		if inProgressCount > 0 {

			var maxSequence []int
			if err := tx.Raw(
				"SELECT MAX(sequence) FROM matches WHERE tournament_id = ?",
				match.TournamentID,
			).Scan(&maxSequence).Error; err != nil {
				return err
			}

			match.Table = ""
			match.State = int(core.Scheduled)
			match.Sequence = maxSequence[0] + 1
		} else {
			match.State = int(core.InProgress)
		}

		if err := tx.Save(match).Error; err != nil {
			return err
		}

		return nil
	})

	return c.NoContent(http.StatusOK)
}
