package features

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/database"
)

type updateTournamentRequest struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	TableCount int    `json:"tableCount"`
}

// UpdateTournament will update a tournament.
func UpdateTournament(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	tournament := new(updateTournamentRequest)
	if err := c.Bind(tournament); err != nil {
		return err
	}
	// TODO: Validate input
	log.Printf("About to save tournament '%s'", tournament.Name)
	if err := updateTournament(ac, *tournament); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

// UpdateTournament saves changes to a Tournament entity
func updateTournament(d *core.FoosmanContext, t updateTournamentRequest) error {
	var tournament database.Tournament
	if err := d.DB.First(&tournament, t.ID).Error; err != nil {
		return err
	}
	return d.DB.Model(&tournament).Updates(map[string]interface{}{
		"name":        t.Name,
		"table_count": t.TableCount,
	}).Error
}
