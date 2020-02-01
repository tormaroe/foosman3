package features

import (
	"log"
	"net/http"

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
	t, err := getTournament(ac, tournamentID)
	if err != nil {
		log.Print("Error getting tournament", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSONPretty(http.StatusOK, t, "  ")
}

func getTournament(d *core.FoosmanContext, ID int) (database.Tournament, error) {
	var t database.Tournament
	if err := d.DB.First(&t, ID).Error; err != nil {
		return t, err
	}

	var teams []database.Team
	if err := d.DB.Model(&t).Related(&teams).Error; err != nil {
		return t, err
	}
	t.Teams = teams

	var groups []database.Group
	if err := d.DB.Model(&t).Related(&groups).Error; err != nil {
		return t, err
	}
	t.Groups = groups

	return t, nil
}
