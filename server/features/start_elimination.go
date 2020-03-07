package features

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/database"
)

type startEliminationRequest struct {
	TeamCount int `json:"teamCount"` // NuMber of teams to include in play-off
}

func StartElimination(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	tournamentID, err := ac.GetParamID()
	if err != nil {
		return err
	}

	var t database.Tournament
	if err := ac.DB.First(&t, tournamentID).Error; err != nil {
		return err
	}

	// TODO: Assert tournament.state == core.GroupPlayDone

	req := new(startEliminationRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	ac.SetResultMux.Lock()
	defer ac.SetResultMux.Unlock()

	if err := ac.DB.Transaction(func(tx *gorm.DB) error {
		if err := createFirstTierElimMatches(tx, tournamentID, req.TeamCount); err != nil {
			return err
		}
		t.State = int(core.EliminationPlayStarted)
		tx.Save(&t)

		return nil
	}); err != nil {
		return err
	}

	done := database.ScheduleUpcoming(ac, t.ID, t.TableCount)
	done.Wait() // Block until initial scheduling done

	// Start group play (Starts TableCount matches, which again schedules new matches)
	for i := 0; i < t.TableCount; i++ {
		table := fmt.Sprintf("Table %d", i+1)
		done = database.StartNextMatch(ac, tournamentID, table)
		done.Wait()
	}

	return c.NoContent(http.StatusOK)
}

func getGroupedTeamOrder(db *gorm.DB, tournamentID int) (map[int][]int, error) {
	q := `
		select 
			m.group_id, 
			r.team_id, 
			sum(r.points) as points,
			random() as rnd
		from matches m
		join match_results r ON r.match_id = m.id
		where m.tournament_id = ?
		group by r.team_id
		order by m.group_id, points desc, rnd
	`
	var pointData []struct {
		GroupID int
		TeamID  int
		Points  int
	}
	if err := db.Raw(q, tournamentID).Scan(&pointData).Error; err != nil {
		return nil, err
	}

	groupedOrderedTeams := make(map[int][]int)
	currGroupID := -1
	gIdx := -1
	for _, pd := range pointData {
		if pd.GroupID != currGroupID {
			gIdx++
			groupedOrderedTeams[gIdx] = []int{}
			log.Printf("Setting key %d", gIdx)
			currGroupID = pd.GroupID
		}
		groupedOrderedTeams[gIdx] = append(groupedOrderedTeams[gIdx], pd.TeamID)
		log.Printf("gIdx: %d teamId: %d", gIdx, pd.TeamID)
	}
	return groupedOrderedTeams, nil
}

func createFirstTierElimMatches(db *gorm.DB, tournamentID int, inclTeamCnt int) error {
	tierMatchCnt := inclTeamCnt / 2

	if tierMatchCnt*2 != inclTeamCnt {
		return fmt.Errorf("Invalid elimination team count: %d", inclTeamCnt)
	}

	gto, err := getGroupedTeamOrder(db, tournamentID)
	if err != nil {
		return err
	}

	teamsFoundCnt := 0
	matches := make([]database.Match, tierMatchCnt)
	matchIdx := 0
	groupTier := 0 // winner, 1 = second place, and so on..
	gIdx := 0
	goingLeft := true

	log.Printf("len(gto): %d", len(gto))

	for teamsFoundCnt < inclTeamCnt {
		log.Printf(
			"gIdx:%d matchIdx:%d groupTier:%d",
			gIdx,
			matchIdx,
			groupTier,
		)

		teamID := gto[gIdx][groupTier]
		match := matches[matchIdx]

		if goingLeft {
			match.Team1ID = teamID
		} else {
			match.Team2ID = teamID
		}
		match.TournamentID = tournamentID
		match.PlayoffTier = tierMatchCnt
		match.PlayoffMatchNumber = matchIdx + 1
		teamsFoundCnt++

		gIdx++
		if gIdx >= len(gto) {
			gIdx = 0
			groupTier++
		}

		matches[matchIdx] = match

		if goingLeft {
			matchIdx++
			if matchIdx >= len(matches) {
				matchIdx--
				goingLeft = false
			}
		} else {
			matchIdx--
			if matchIdx < 0 && teamsFoundCnt != inclTeamCnt {
				return fmt.Errorf("Loop issue")
			}
		}
	}

	// Save matches w/ match results
	for i := 0; i < len(matches); i++ {
		match := matches[i]
		if err := db.Create(&match).Error; err != nil {
			return err
		}

		if err := db.Create(&database.MatchResult{
			TeamID:  match.Team1ID,
			MatchID: match.ID,
			Points:  0,
			Win:     0,
			Loss:    0,
			Draw:    0,
		}).Error; err != nil {
			return err
		}

		if err := db.Create(&database.MatchResult{
			TeamID:  match.Team2ID,
			MatchID: match.ID,
			Points:  0,
			Win:     0,
			Loss:    0,
			Draw:    0,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}
