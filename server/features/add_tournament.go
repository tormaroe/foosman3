package features

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/database"
)

type addTournamentRequest struct {
	Name       string `json:"name"`
	TableCount int    `json:"tableCount"`
}

// AddTournament creates a new Tournament
func AddTournament(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	tournament := new(addTournamentRequest)
	if err := c.Bind(tournament); err != nil {
		return err
	}
	// TODO: Validate input
	log.Printf("About to save tournament '%s'", tournament.Name)
	if err := addTournament(ac, *tournament); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

// AddTournament saves a new Tournament entity
func addTournament(d *core.FoosmanContext, t addTournamentRequest) error {
	return d.DB.Create(&database.Tournament{
		Name:       t.Name,
		TableCount: t.TableCount,
		State:      int(core.New),
	}).Error
}
