package features

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
)

type getTeamResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Player1 string `json:"player1"`
	Player2 string `json:"player2"`
	Player3 string `json:"player3"`
}

func GetTeams(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	tournamentID, err := ac.GetParamID()
	if err != nil {
		return err
	}
	teams, err := getTeams(ac, tournamentID)
	if err != nil {
		log.Print("Error getting teams", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSONPretty(http.StatusOK, teams, "  ")
}

// GetTournamentTeams gets all the teams for a Tournament from the database
func getTeams(d *core.FoosmanContext, ID int) ([]getTeamResponse, error) {
	rows, err := d.DB.Query(`
		select id, name, player_1, player_2, player_3
		from team 
		where tournament_id=?
	`, ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var result []getTeamResponse
	for rows.Next() {
		var t getTeamResponse
		err = rows.Scan(&t.ID, &t.Name, &t.Player1, &t.Player2, &t.Player3)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, rows.Err()
}
