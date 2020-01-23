package features

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
)

type updateTournamentRequest struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	TableCount int    `json:"tableCount"`
}

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
	stmt, err := d.DB.Prepare(`
		update tournament
		set name=?, table_count=?
		where id=?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(t.Name, t.TableCount, t.ID)
	return err
}
