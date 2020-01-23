package database

import (
	"github.com/tormaroe/foosman3/server/core"
	"reflect"
	"testing"
)

func TestInitDatabase(t *testing.T) {
	d, err := Init(":memory:")
	if err != nil {
		t.Errorf("Error in Init(): %s", err)
	}
	defer d.Close()
}

func TestAddTeam(t *testing.T) {
	d, _ := Init(":memory:")
	defer d.Close()

	err := d.AddTeam(1, core.Team{
		Name:    "Bob & Janet",
		Player1: "Bob",
		Player2: "Janet",
		Player3: "",
	})
	if err != nil {
		t.Errorf("Error saving team: %s", err)
	}
}

func TestUpdateTeam(t *testing.T) {
	d, _ := Init(":memory:")
	defer d.Close()

	const tournamentID = 1

	d.AddTeam(tournamentID, core.Team{Name: "Bob & Janet", Player1: "Bob", Player2: "Janet"})

	teams, _ := d.GetTournamentTeams(tournamentID)
	team := teams[0]

	if team.ID <= 0 {
		t.Errorf("Team ID is %d", team.ID)
	}

	team.Name = "The Foosloose"
	team.Player1 = "Bob'er"
	team.Player2 = "Janetz"

	err := d.UpdateTeam(team)
	if err != nil {
		t.Errorf("Error updating team: %s", err)
	}

	teams, _ = d.GetTournamentTeams(tournamentID)
	if len(teams) != 1 {
		t.Errorf("Expected one team, got %d", len(teams))
	}
	teamUpdated := teams[0]
	if !reflect.DeepEqual(team, teamUpdated) {
		t.Error("Team not properly updated")
	}
}

func TestGetTournamentTeams(t *testing.T) {
	d, _ := Init(":memory:")
	defer d.Close()

	d.AddTeam(1, core.Team{Name: "Bob & Janet", Player1: "Bob", Player2: "Janet"})
	d.AddTeam(1, core.Team{Name: "Jack & Jill", Player1: "Jack", Player2: "Jill"})

	teams, err := d.GetTournamentTeams(1)
	if err != nil {
		t.Errorf("Error getting teams: %s", err)
	}

	if len(teams) != 2 {
		t.Errorf("Expected 2 teams, but got %d", len(teams))
	}
}
