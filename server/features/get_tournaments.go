package features

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
)

type tournamentResponse struct {
	ID         int                  `json:"id"`
	Name       string               `json:"name"`
	TableCount int                  `json:"tableCount"`
	State      core.TournamentState `json:"state"`
}

// GetTournaments responds to a GET request for all tournaments
func GetTournaments(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	lst, err := getTournaments(ac)
	if err != nil {
		log.Print("Error getting tournaments", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSONPretty(http.StatusOK, lst, "  ")
}

func getTournaments(d *core.FoosmanContext) ([]tournamentResponse, error) {
	rows, err := d.DB.Query(`
		select id, name, table_count, state
		from tournament
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var result []tournamentResponse
	for rows.Next() {
		var t tournamentResponse
		err = rows.Scan(&t.ID, &t.Name, &t.TableCount, &t.State)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, rows.Err()
}
