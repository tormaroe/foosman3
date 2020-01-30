package features

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
)

type getTournamentResponse struct {
	ID         int                  `json:"id"`
	Name       string               `json:"name"`
	TableCount int                  `json:"tableCount"`
	State      core.TournamentState `json:"state"`
	Teams      []teamDTO            `json:"teams"`
}

type teamDTO struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Player1 string `json:"player1"`
	Player2 string `json:"player2"`
	Player3 string `json:"player3"`
	GroupID *int   `json:"groupId"`
}

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

func getTournament(d *core.FoosmanContext, ID int) (getTournamentResponse, error) {
	row := d.DB.QueryRow(`
		select id, name, table_count, state
		from tournament
		where id=?
	`, ID)
	var t getTournamentResponse
	err := row.Scan(&t.ID, &t.Name, &t.TableCount, &t.State)
	if err == nil && t.ID > 0 {
		if teams, err := getTeams(d, ID); err == nil {
			t.Teams = teams
		}
	}
	return t, err
}

// getTeams gets all the teams for a Tournament from the database
func getTeams(d *core.FoosmanContext, ID int) ([]teamDTO, error) {
	rows, err := d.DB.Query(`
		select id, name, player_1, player_2, player_3, group_id
		from team 
		where tournament_id=?
	`, ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var result []teamDTO
	for rows.Next() {
		var t teamDTO
		err = rows.Scan(&t.ID, &t.Name, &t.Player1, &t.Player2, &t.Player3, &t.GroupID)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, rows.Err()
}
