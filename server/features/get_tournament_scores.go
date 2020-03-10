package features

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
)

type teamScores struct {
	TeamID      int
	PlayedCount int
	Points      int
	Wins        int
	Losses      int
	Draws       int
}

// GetTournamentScores aggregates match statistics
// for all teams in a tournament.
func GetTournamentScores(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	tournamentID, err := ac.GetParamID()
	if err != nil {
		return err
	}
	res, err := aggregateTournamentScores(ac.DB, tournamentID)
	if err != nil {
		return nil
	}
	return c.JSONPretty(http.StatusOK, res, "  ")
}

func aggregateTournamentScores(db *gorm.DB, ID int) ([]teamScores, error) {
	query := `
		SELECT
		  r.team_id,
		  sum(case when m.state = 3 then 1 else 0 end) as played_count,
		  sum(r.points) as points,
		  sum(r.win) as wins,
		  sum(r.loss) as losses,
		  sum(r.draw) as draws
		FROM matches m
		join match_results r ON r.match_id = m.id
		where m.tournament_id = ?
		  and m.group_id > 0
		group by r.team_id
	`
	var res []teamScores
	err := db.Raw(query, ID).Scan(&res).Error
	return res, err
}
