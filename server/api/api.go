package api

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
)

func getTournaments(c echo.Context) error {
	ac := c.(*apiContext)
	lst, err := ac.db.GetTournaments()
	if err != nil {
		log.Print("Error getting tournaments", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSONPretty(http.StatusOK, lst, "  ")
}

func getTournamentTeams(c echo.Context) error {
	ac := c.(*apiContext)
	tournamentID, err := ac.getParamID()
	if err != nil {
		return err
	}
	teams, err := ac.db.GetTournamentTeams(tournamentID)
	if err != nil {
		log.Print("Error getting teams", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSONPretty(http.StatusOK, teams, "  ")
}

func addTournament(c echo.Context) error {
	ac := c.(*apiContext)
	tournament := new(core.Tournament)
	if err := c.Bind(tournament); err != nil {
		return err
	}
	// TODO: Validate input
	log.Printf("About to save tournament '%s'", tournament.Name)
	if err := ac.db.AddTournament(*tournament); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func updateTournament(c echo.Context) error {
	ac := c.(*apiContext)
	tournament := new(core.Tournament)
	if err := c.Bind(tournament); err != nil {
		return err
	}
	// TODO: Validate input
	log.Printf("About to save tournament '%s'", tournament.Name)
	if err := ac.db.UpdateTournament(*tournament); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func addTeam(c echo.Context) error {
	ac := c.(*apiContext)
	tournamentID, err := ac.getParamID()
	if err != nil {
		return err
	}
	team := new(core.Team)
	if err := c.Bind(team); err != nil {
		return err
	}
	// TODO: Validate input
	log.Printf("About to save team '%s'", team.Name)
	if err := ac.db.AddTeam(tournamentID, *team); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
