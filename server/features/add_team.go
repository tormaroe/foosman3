package features

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
)

type addTeamRequest struct {
	Name    string `json:"name"`
	Player1 string `json:"player1"`
	Player2 string `json:"player2"`
	Player3 string `json:"player3"`
}

// AddTeam creates a new Team and adds it to an existing Tournament.
func AddTeam(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	tournamentID, err := ac.GetParamID()
	if err != nil {
		return err
	}
	team := new(addTeamRequest)
	if err := c.Bind(team); err != nil {
		return err
	}
	// TODO: Validate input
	log.Printf("About to save team '%s'", team.Name)
	if err := addTeam(ac, tournamentID, *team); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

// AddTeam saves a new Team entity
func addTeam(d *core.FoosmanContext, tournamentID int, t addTeamRequest) error {
	// TODO: Don't add if tournament has started
	stmt, err := d.DB.Prepare(`
		insert into team
		(name, tournament_id, player_1, player_2, player_3) 
		values
		(?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(t.Name, tournamentID, t.Player1, t.Player2, t.Player3)
	return err
}
