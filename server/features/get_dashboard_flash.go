package features

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/database"
)

type dashboardFlash struct {
	Raw string
}

func GetDashboardFlash(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	ID, err := ac.GetParamID()
	if err != nil {
		return err
	}

	var tournament database.Tournament
	if err := ac.DB.Find(&tournament, ID).Error; err != nil {
		return err
	}

	res := dashboardFlash{
		Raw: tournament.Name,
	}

	if tournament.State == int(core.New) {
		option := rand.Intn(2)
		if option == 0 {
			res.Raw = "hello, world!"
		} else if option == 1 {
			res.Raw = tournament.Name
		}

	} else if tournament.State == int(core.GroupPlayStarted) {
		res.setFlashForGroupPlayStarted(ac.DB, ID)

	} else if tournament.State == int(core.GroupPlayDone) {
		option := rand.Intn(4)
		if option == 0 {
			res.Raw = "Group play is finished"
		} else if option == 1 {
			res.Raw = "Configuring playoff matches.."
		} else if option == 2 {
			res.Raw = "Elimination matches will begin soon"
		}

	} else if tournament.State == int(core.EliminationPlayStarted) {
		res.setFlashForEliminationPlay(ac.DB, ID)

	} else if tournament.State == int(core.Done) {
		res.setFlashForDone(ac.DB, ID)

	}

	return c.JSONPretty(http.StatusOK, res, "  ")
}

func (d *dashboardFlash) setFlashForEliminationPlay(db *gorm.DB, ID int) {
	var matches []database.Match
	if err := db.Preload("Team1").Preload("Team2").Preload("MatchResults").Where(
		"tournament_id = ? and playoff_tier > 0 and state = ?",
		ID,
		int(core.Played),
	).Find(&matches).Error; err != nil {
		return
	}

	if len(matches) == 0 {
		d.Raw = "Elimination play started!"
		return
	}

	tierN, tierMatches := filterTopTier(matches)

	if tierN == 2 && len(tierMatches) == 2 {
		w1, _, _, e1 := tierMatches[0].GetWinnerAndLooser()
		w2, _, _, e2 := tierMatches[1].GetWinnerAndLooser()
		if e1 != nil || e2 != nil {
			return
		}
		d.Raw = fmt.Sprintf(
			"Final match:<br>%s vs. %s",
			w1.Name,
			w2.Name,
		)
		return
	}

	focusMatch := tierMatches[rand.Intn(len(tierMatches))]
	if winner, looser, _, err := focusMatch.GetWinnerAndLooser(); err == nil {
		d.Raw = fmt.Sprintf(
			"%s has eliminated %s",
			winner.Name,
			looser.Name,
		)
	}
}

func filterTopTier(ms []database.Match) (int, []database.Match) {
	tierN := 10000 // a big number
	for i := 0; i < len(ms); i++ {
		if ms[i].PlayoffTier < tierN {
			tierN = ms[i].PlayoffTier
		}
	}
	var tierMatches []database.Match
	for i := 0; i < len(ms); i++ {
		if ms[i].PlayoffTier == tierN {
			tierMatches = append(tierMatches, ms[i])
		}
	}
	return tierN, tierMatches
}

func (d *dashboardFlash) setFlashForDone(db *gorm.DB, ID int) {
	option := rand.Intn(3)
	if option == 0 {
		d.Raw = "Thank you for playing!"
	} else {
		var finals []database.Match
		if err := db.Preload("Team1").Preload("Team2").Preload("MatchResults").Where(
			"tournament_id = ? and playoff_tier = 1",
			ID,
		).Find(&finals).Error; err != nil {
			return
		}

		winner, _, _, err := finals[0].GetWinnerAndLooser()
		if err != nil {
			return
		}
		d.Raw = fmt.Sprintf("Congratulations to the winner<br>%s", winner.Name)
	}
}

func (d *dashboardFlash) setFlashForGroupPlayStarted(db *gorm.DB, ID int) {
	matches, err := getPlayedTournamentMatches(db, ID)

	if err != nil {
		return
	}

	if len(matches) == 0 {
		d.Raw = "Waiting for first match result.."
		return
	}

	if len(matches) == 1 {
		if winner, _, isDraw, err := matches[0].GetWinnerAndLooser(); err == nil {
			if isDraw {
				d.Raw = "First match ended in a draw"
			} else {
				d.Raw = fmt.Sprintf("%s won the first match of the tournament", winner.Name)
			}
		}
		return
	}

	option := rand.Intn(3)
	if option == 0 {
		if scores, err := aggregateTournamentScores(db, ID); err == nil {
			best := getBestScores(scores)
			if team, err := getTeam(db, best.TeamID); err == nil {
				d.Raw = fmt.Sprintf(
					"%s has %d points after %d matches in %s",
					team.Name,
					best.Points,
					best.PlayedCount,
					team.Group.Name,
				)
			}
		}

	} else if option == 1 {
		d.Raw = fmt.Sprintf("%d matches played", len(matches))

	} else {
		var totalMatchCount int
		if err := db.Table("matches").Where("tournament_id = ?", ID).Count(&totalMatchCount).Error; err != nil {
			return
		}
		d.Raw = fmt.Sprintf(
			"Group play progress: %d%%",
			(len(matches)*100)/totalMatchCount,
		)
	}

	return
}

func getBestScores(ss []teamScores) teamScores {
	// Assumption: len(ss) > 1
	best := ss[0]
	for i := 1; i < len(ss); i++ {
		if ss[i].Points > best.Points {
			best = ss[i]
		}
	}
	return best
}
