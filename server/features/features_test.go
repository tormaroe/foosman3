package features

import (
	"errors"
	"testing"

	"gotest.tools/assert"

	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/database"
)

func TestFeatures(t *testing.T) {
	// 0. Setup
	db, err := database.Init(":memory:")
	//db, err := database.Init("test.sqlite")
	assert.NilError(t, err)
	defer db.Close()
	scheduleChan := database.NewScheduleChan()
	startNextMatchChan := database.NewStartMatchChan()
	cnx := &core.FoosmanContext{
		DB:                 db,
		ScheduleChan:       scheduleChan,
		StartNextMatchChan: startNextMatchChan,
	}

	// 1. Create tournament
	assert.NilError(t, addTournament(db, addTournamentRequest{
		Name:       "TournamentOne",
		TableCount: 2,
	}))

	// Get and validate tournament
	tournaments, err := getTournaments(db)
	assert.NilError(t, err)
	assert.Assert(t, len(tournaments) == 1)

	tournament := tournaments[0]
	assert.Equal(t, tournament.ID, 1)
	assert.Equal(t, tournament.Name, "TournamentOne")
	assert.Equal(t, tournament.TableCount, 2)
	assert.Equal(t, tournament.State, int(core.Planned))

	// 2. Create teams
	assert.NilError(t, addTeam(db, tournament.ID, addTeamRequest{Name: "Team 1"}))
	assert.NilError(t, addTeam(db, tournament.ID, addTeamRequest{Name: "Team 2"}))
	assert.NilError(t, addTeam(db, tournament.ID, addTeamRequest{Name: "Team 3"}))
	assert.NilError(t, addTeam(db, tournament.ID, addTeamRequest{Name: "Team A"}))
	assert.NilError(t, addTeam(db, tournament.ID, addTeamRequest{Name: "Team B"}))
	assert.NilError(t, addTeam(db, tournament.ID, addTeamRequest{Name: "Team C"}))
	assert.NilError(t, addTeam(db, tournament.ID, addTeamRequest{Name: "Team D"}))

	// Get and validate teams
	tournament, err = getTournament(db, tournament.ID)
	assert.NilError(t, err)
	assert.Assert(t, len(tournament.Teams) == 7)

	team1, _ := ByName(tournament.Teams, "Team 1")
	team2, _ := ByName(tournament.Teams, "Team 2")
	team3, _ := ByName(tournament.Teams, "Team 3")
	teamA, _ := ByName(tournament.Teams, "Team A")
	teamB, _ := ByName(tournament.Teams, "Team B")
	teamC, _ := ByName(tournament.Teams, "Team C")
	teamD, _ := ByName(tournament.Teams, "Team D")

	// 3. Create groups
	assert.NilError(t, setGroups(db, tournament.ID, &[]groupDefinition{
		groupDefinition{Name: "Group I", TeamIDs: []int{team1.ID, team2.ID, team3.ID}},
		groupDefinition{Name: "Group II", TeamIDs: []int{teamA.ID, teamB.ID, teamC.ID, teamD.ID}},
	}))

	// Get and validate groups
	tournament, err = getTournament(db, tournament.ID)
	assert.NilError(t, err)
	assert.Assert(t, len(tournament.Groups) == 2)

	// 4. Generate matches
	assert.NilError(t, generateMatches(tournament)(db))

	// TODO: Get and validate matches

	// 5. Schedule matches
	done := database.ScheduleUpcoming(cnx, tournament.ID, 3)
	done.Wait()

	// 6. Start first games (2 tables)
	done = database.StartNextMatch(cnx, tournament.ID, "Table 1")
	done.Wait()
	done = database.StartNextMatch(cnx, tournament.ID, "Table 2")
	done.Wait()

	// TODO: Register result. Should start a game and schedule one more

	// TODO: Register result for all group games.
	// TODO: Group results.
	// TODO: Start CUP

}

func ByName(teams []database.Team, name string) (database.Team, error) {
	for _, t := range teams {
		if t.Name == name {
			return t, nil
		}
	}
	return database.Team{}, errors.New("Team missing")
}
