package features

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/database"
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

	if err := database.AssertTournamentNotStarted(ac, tournamentID); err != nil {
		return err
	}

	team := new(addTeamRequest)
	if err := c.Bind(team); err != nil {
		return err
	}

	// TODO: Validate input

	log.Printf("About to save team '%s'", team.Name)
	if err := addTeam(ac.DB, tournamentID, *team); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

// AddTeam saves a new Team entity
func addTeam(db *gorm.DB, tournamentID int, t addTeamRequest) error {
	var tournament database.Tournament
	if err := db.First(&tournament, tournamentID).Error; err != nil {
		return err
	}

	return db.Create(&database.Team{
		Name:       t.Name,
		Tournament: tournament,
		Player1:    t.Player1,
		Player2:    t.Player2,
		Player3:    t.Player3,
	}).Error
}
