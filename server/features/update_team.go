package features

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
)

type updateTeamRequest struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Player1 string `json:"player1"`
	Player2 string `json:"player2"`
	Player3 string `json:"player3"`
}

// UpdateTeam will update an existing Team (name and players).
func UpdateTeam(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	team := new(updateTeamRequest)
	if err := c.Bind(team); err != nil {
		return err
	}
	// TODO: Validate input
	log.Printf("About to save team '%s'", team.Name)
	if err := updateTeam(ac, *team); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

// UpdateTeam saves changes to a Team entity
func updateTeam(d *core.FoosmanContext, t updateTeamRequest) error {
	stmt, err := d.DB.Prepare(`
		update team 
		set name=?, player_1=?, player_2=?, player_3=? 
		where id=?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(t.Name, t.Player1, t.Player2, t.Player3, t.ID)
	return err
}
