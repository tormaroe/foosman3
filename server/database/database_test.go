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

func TestInsertTeam(t *testing.T) {
	d, _ := Init(":memory:")
	defer d.Close()

	err := d.SaveTeam(core.Team{Name: "Bob & Janet", Player1: "Bob", Player2: "Janet"})
	if err != nil {
		t.Errorf("Error saving team: %s", err)
	}
}

func TestUpdateTeam(t *testing.T) {
	d, _ := Init(":memory:")
	defer d.Close()

	d.SaveTeam(core.Team{Name: "Bob & Janet", Player1: "Bob", Player2: "Janet"})

	teams, _ := d.AllTeams()
	team := teams[0]

	if team.ID <= 0 {
		t.Errorf("Team ID is %d", team.ID)
	}

	team.Name = "The Foosloose"
	team.Player1 = "Bob'er"
	team.Player2 = "Janetz"

	err := d.SaveTeam(team)
	if err != nil {
		t.Errorf("Error updating team: %s", err)
	}

	teams, _ = d.AllTeams()
	if len(teams) != 1 {
		t.Errorf("Expected one team, got %d", len(teams))
	}
	teamUpdated := teams[0]
	if !reflect.DeepEqual(team, teamUpdated) {
		t.Error("Team not properly updated")
	}
}

func TestGetAllTeams(t *testing.T) {
	d, _ := Init(":memory:")
	defer d.Close()

	d.SaveTeam(core.Team{Name: "Bob & Janet", Player1: "Bob", Player2: "Janet"})
	d.SaveTeam(core.Team{Name: "Jack & Jill", Player1: "Jack", Player2: "Jill"})

	teams, err := d.AllTeams()
	if err != nil {
		t.Errorf("Error getting teams: %s", err)
	}

	if len(teams) != 2 {
		t.Errorf("Expected 2 teams, but got %d", len(teams))
	}
}
