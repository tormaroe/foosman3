package features

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/database"
)

// GetTournamentByID gets a Tournament by ID.
// Response contains a list of all teams.
func GetTournamentByID(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	tournamentID, err := ac.GetParamID()
	if err != nil {
		return err
	}
	t, err := getTournament(ac.DB, tournamentID)
	if err != nil {
		log.Print("Error getting tournament", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSONPretty(http.StatusOK, t, "  ")
}

func getTournament(db *gorm.DB, ID int) (database.Tournament, error) {
	var t database.Tournament
	if err := db.Preload("Teams").Preload("Groups").First(&t, ID).Error; err != nil {
		return t, err
	}

	return t, nil
}
