package features

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
)

type gTMResult struct {
	Team1     string
	Team2     string
	GroupName string
}

// TODO: Test needed

// ! Unused

func GetTournamentMatches(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	tournamentID, err := ac.GetParamID()
	if err != nil {
		return err
	}
	// var matches []database.Match
	// if err := ac.DB.Where("tournament_id = ?", tournamentID).Find(&matches).Error; err != nil {
	// 	return err
	// }

	var result []gTMResult
	ac.DB.
		Table("matches").
		Select("t1.name as team1, t2.name as team2, groups.name as group_name").
		Joins("join groups on groups.id = matches.group_id and groups.tournament_id = ?", tournamentID).
		Joins("left join teams t1 on t1.id = matches.team1_id").
		Joins("left join teams t2 on t2.id = matches.team2_id").
		Scan(&result)

	return c.JSONPretty(http.StatusOK, result, "  ")
}
