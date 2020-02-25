package features

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/tormaroe/foosman3/server/core"
)

type matchSchedule struct {
	ID        int    `json:"id"`
	Team1ID   int    `json:"team1Id"`
	Team1Name string `json:"team1Name"`
	Team2ID   int    `json:"team2Id"`
	Team2Name string `json:"team2Name"`
	GroupName string `json:"groupName"`
	Sequence  int    `json:"sequence"`
}

func GetMatchesScheduled(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	tournamentID, err := ac.GetParamID()
	if err != nil {
		return err
	}

	var result []matchSchedule
	if err := ac.DB.Table(
		"matches",
	).Select(
		"matches.id, t1.id as team1_id, t1.name as team1_name, t2.id as team2_id, t2.name as team2_name, g.name as group_name, matches.[sequence]",
	).Joins(
		"join groups g on g.id = matches.group_id",
	).Joins(
		"join teams t1 on t1.id = matches.team1_id",
	).Joins(
		"join teams t2 on t2.id = matches.team2_id",
	).Where(
		"matches.tournament_id = ? and matches.state = ?",
		tournamentID,
		int(core.Scheduled),
	).Order(
		"matches.[sequence]",
	).Scan(&result).Error; err != nil {
		log.Print(err)
	}

	return c.JSONPretty(http.StatusOK, result, "  ")
}
