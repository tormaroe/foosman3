package features

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/database"
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
	var team database.Team
	return d.DB.First(&team, t.ID).Updates(map[string]interface{}{
		"name":    t.Name,
		"player1": t.Player1,
		"player2": t.Player2,
		"player3": t.Player3,
	}).Error
}
