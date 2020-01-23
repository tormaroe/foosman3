package features

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
)

type addTournamentRequest struct {
	Name       string `json:"name"`
	TableCount int    `json:"tableCount"`
}

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
	stmt, err := d.DB.Prepare(`
		insert into tournament
		(name, table_count, state)
		values
		(?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(t.Name, t.TableCount, core.New)
	return err
}
