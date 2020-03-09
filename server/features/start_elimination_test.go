package features

import "testing"

import "gotest.tools/assert"

import "github.com/tormaroe/foosman3/server/database"

func TestXXX(t *testing.T) {
	tierMatchCount := 4

	gtps := []groupTeamPoints{
		groupTeamPoints{GroupID: 1, TeamID: 4, Points: 8},
		groupTeamPoints{GroupID: 1, TeamID: 3, Points: 8},
		groupTeamPoints{GroupID: 1, TeamID: 2, Points: 4}, // second best 3rd place
		groupTeamPoints{GroupID: 1, TeamID: 1, Points: 1},

		groupTeamPoints{GroupID: 2, TeamID: 5, Points: 6},
		groupTeamPoints{GroupID: 2, TeamID: 6, Points: 5},
		groupTeamPoints{GroupID: 2, TeamID: 7, Points: 3}, // worst 3rd place
		groupTeamPoints{GroupID: 2, TeamID: 8, Points: 3},

		groupTeamPoints{GroupID: 3, TeamID: 9, Points: 6},
		groupTeamPoints{GroupID: 3, TeamID: 10, Points: 5},
		groupTeamPoints{GroupID: 3, TeamID: 11, Points: 5}, // best 3rd place
		groupTeamPoints{GroupID: 3, TeamID: 12, Points: 4},
	}

	gtpColl := newGroupTeamPointsCollection(gtps)

	assert.Equal(t, 3, gtpColl.groupCount())

	matches := make([]database.Match, tierMatchCount)
	err := gtpColl.fillMatches(&matches)

	assert.NilError(t, err)

	assert.Equal(t, matches[0].Team1ID, 4)
	assert.Equal(t, matches[0].Team2ID, 10)

	assert.Equal(t, matches[1].Team1ID, 5)
	assert.Equal(t, matches[1].Team2ID, 11)

	assert.Equal(t, matches[2].Team1ID, 9)
	assert.Equal(t, matches[2].Team2ID, 3)

	assert.Equal(t, matches[3].Team1ID, 6)
	assert.Equal(t, matches[3].Team1ID, 2)
}

// Must work when group count <= finals width
// TODO: Restrict group count > finals width (UI and server)
