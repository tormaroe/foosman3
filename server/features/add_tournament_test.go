package features

import (
	"testing"

	"github.com/tormaroe/foosman3/server/database"
)

func TestAddTournament(t *testing.T) {
	db, _ := database.Init(":memory:")
	defer db.Close()
	var err error

	if err = addTournament(db, addTournamentRequest{
		Name:       "T1",
		TableCount: 5,
	}); err != nil {
		t.Errorf("Error adding tournament: %s", err)
	}
}
