package features

import (
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
	"testing"

	"github.com/tormaroe/foosman3/server/database"
)

func TestGetTournaments(t *testing.T) {
	db, _ := database.Init(":memory:")
	defer db.Close()
	var err error

	addTournament(db, addTournamentRequest{
		Name:       "T1",
		TableCount: 5,
	})
	addTournament(db, addTournamentRequest{
		Name:       "T2",
		TableCount: 5,
	})

	var tournaments []database.Tournament
	if tournaments, err = getTournaments(db); err != nil {
		t.Errorf("Error getting tournaments: %s", err)
	}
	assert.Assert(t, is.Len(tournaments, 2))
}
