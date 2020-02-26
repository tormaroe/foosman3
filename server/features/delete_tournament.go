package features

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/database"
)

func DeleteTournament(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	ID, err := ac.GetParamID()
	if err != nil {
		return err
	}

	ac.DB.Transaction(func(tx *gorm.DB) error {

		var matches []database.Match
		if err := tx.Where("tournament_id = ?", ID).Find(&matches).Error; err != nil {
			return err
		}

		for _, m := range matches {
			if err := tx.Delete(database.MatchResult{}, "match_id = ?", m.ID).Error; err != nil {
				return err
			}
		}

		if err := tx.Delete(database.Match{}, "tournament_id = ?", ID).Error; err != nil {
			return err
		}
		if err := tx.Delete(database.Team{}, "tournament_id = ?", ID).Error; err != nil {
			return err
		}
		if err := tx.Delete(database.Group{}, "tournament_id = ?", ID).Error; err != nil {
			return err
		}
		if err := tx.Delete(database.Tournament{ID: ID}).Error; err != nil {
			return err
		}

		return nil
	})

	return c.NoContent(http.StatusOK)
}
