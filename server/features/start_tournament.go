package features

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/database"
)

// StartTournament ..
func StartTournament(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	tournamentID, err := ac.GetParamID()
	if err != nil {
		return err
	}

	if err := database.AssertTournamentNotStarted(ac, tournamentID); err != nil {
		return err
	}

	// TODO: Assert no un-grouped teams

	var t database.Tournament
	if err := ac.DB.Preload("Groups").First(&t, tournamentID).Error; err != nil {
		return err
	}

	ac.SetResultMux.Lock()
	defer ac.SetResultMux.Unlock()

	if err := ac.DB.Transaction(generateMatches(t)); err != nil {
		return err
	}

	done := database.ScheduleUpcoming(ac, t.ID, t.TableCount)
	done.Wait() // Block until initial scheduling done

	// TODO: Start group play (Starts TableCount matches, which again schedules new matches)
	for i := 0; i < t.TableCount; i++ {
		table := fmt.Sprintf("Table %d", i+1)
		done = database.StartNextMatch(ac, tournamentID, table)
		done.Wait()
	}

	return c.NoContent(http.StatusOK)
}

func generateMatches(t database.Tournament) func(*gorm.DB) error {
	return func(tx *gorm.DB) error {

		for _, g := range t.Groups {
			if err := generateGroupMatches(tx, g); err != nil {
				return err
			}
		}

		t.State = int(core.GroupPlayStarted)
		tx.Save(&t)

		return nil
	}
}

func generateGroupMatches(tx *gorm.DB, g database.Group) error {
	// Get Teams in group
	var teams []database.Team
	if err := tx.Where("group_id = ?", g.ID).Find(&teams).Error; err != nil {
		return err
	}
	// Create matches for all permutations
	for i := 0; i < len(teams)-1; i++ {
		for j := i + 1; j < len(teams); j++ {
			match := database.Match{
				GroupID:      g.ID,
				Team1ID:      teams[i].ID,
				Team2ID:      teams[j].ID,
				TournamentID: g.TournamentID,
				State:        int(core.Planned),
				Sequence:     0,
			}
			if err := tx.Create(&match).Error; err != nil {
				return err
			}

			if err := tx.Create(&database.MatchResult{
				TeamID:  teams[i].ID,
				MatchID: match.ID,
				Points:  0,
				Win:     0,
				Loss:    0,
				Draw:    0,
			}).Error; err != nil {
				return err
			}

			if err := tx.Create(&database.MatchResult{
				TeamID:  teams[j].ID,
				MatchID: match.ID,
				Points:  0,
				Win:     0,
				Loss:    0,
				Draw:    0,
			}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
