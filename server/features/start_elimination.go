package features

import (
	"errors"
	"fmt"
	"net/http"
	"sort"

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

type groupTeamPoints struct {
	GroupID int
	TeamID  int
	Points  int
}

type groupTeamPointsCollection struct {
	gtpMap      map[int][]groupTeamPoints
	teamsFilled int
}

func newGroupTeamPointsCollection(gtps []groupTeamPoints) groupTeamPointsCollection {
	m := make(map[int][]groupTeamPoints)
	currGroupID := -1
	gIdx := -1
	for _, pd := range gtps {
		if pd.GroupID != currGroupID {
			gIdx++
			m[gIdx] = []groupTeamPoints{}
			currGroupID = pd.GroupID
		}
		m[gIdx] = append(m[gIdx], pd)
	}
	return groupTeamPointsCollection{gtpMap: m}
}

func (gtpColl *groupTeamPointsCollection) groupCount() int {
	return len(gtpColl.gtpMap)
}

func (gtpColl *groupTeamPointsCollection) fillMatches(ms *[]database.Match) error {
	mCount := len(*ms)
	grCount := gtpColl.groupCount()
	teamsGoal := mCount * 2
	currentPlaceIdx := 0

	for teamsGoal > gtpColl.teamsFilled { // ! DANGER: May loop forever if too few teams
		tempGtps := gtpColl.getTeamsByPlace(currentPlaceIdx)
		missingCount := teamsGoal - gtpColl.teamsFilled
		if missingCount >= grCount {
			start := 0
			if currentPlaceIdx%2 == 1 {
				start = mCount / 2
			}
			if err := addTeamsToMatches(ms, tempGtps, start); err != nil {
				return err
			}
			currentPlaceIdx++
			gtpColl.teamsFilled += grCount
		} else {
			sort.Slice(tempGtps, func(i, j int) bool {
				return tempGtps[i].Points > tempGtps[j].Points
			})
			tempGtps = tempGtps[:missingCount]
			if err := addTeamsToAnyMatches(ms, tempGtps); err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

func (gtpColl *groupTeamPointsCollection) getTeamsByPlace(idx int) []groupTeamPoints {
	grCount := gtpColl.groupCount()
	tempGtps := make([]groupTeamPoints, grCount)
	for gIdx := 0; gIdx < grCount; gIdx++ {
		tempGtps[gIdx] = gtpColl.gtpMap[gIdx][idx]
	}
	return tempGtps
}

func addTeamsToMatches(ms *[]database.Match, gtps []groupTeamPoints, start int) error {
	mIdx := start
	msLen := len(*ms)
	for gtpIdx := 0; gtpIdx < len(gtps); gtpIdx++ {
		m := &(*ms)[mIdx]
		if m.Team1ID == 0 {
			m.Team1ID = gtps[gtpIdx].TeamID
		} else if m.Team2ID == 0 {
			m.Team2ID = gtps[gtpIdx].TeamID
		} else {
			return errors.New("Too few matches compared to groups")
		}
		mIdx++
		if mIdx == msLen {
			mIdx = 0
		}
	}
	return nil
}

func addTeamsToAnyMatches(ms *[]database.Match, gtps []groupTeamPoints) error {
	gtpIdx := 0
	for mIdx := 0; mIdx < len(*ms); mIdx++ {
		m := &(*ms)[mIdx]
		if m.Team1ID == 0 {
			m.Team1ID = gtps[gtpIdx].TeamID
			gtpIdx++
		} else if m.Team2ID == 0 {
			m.Team2ID = gtps[gtpIdx].TeamID
			gtpIdx++
		}
		if gtpIdx >= len(gtps) {
			return nil // DONE
		}
	}
	return errors.New("Unable to find matches for all playoff teams")
}

func loadGroupTeamPoints(db *gorm.DB, tournamentID int) ([]groupTeamPoints, error) {
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
	var gtps []groupTeamPoints
	if err := db.Raw(q, tournamentID).Scan(&gtps).Error; err != nil {
		return nil, err
	}
	return gtps, nil
}

func createFirstTierElimMatches(db *gorm.DB, tournamentID int, inclTeamCnt int) error {
	tierMatchCnt := inclTeamCnt / 2

	if tierMatchCnt*2 != inclTeamCnt {
		return fmt.Errorf("Invalid elimination team count: %d", inclTeamCnt)
	}

	gtos, err := loadGroupTeamPoints(db, tournamentID)
	if err != nil {
		return err
	}
	matches := make([]database.Match, tierMatchCnt)
	gtoColl := newGroupTeamPointsCollection(gtos)
	if err := gtoColl.fillMatches(&matches); err != nil {
		return err
	}

	// Save matches w/ match results
	for i := 0; i < len(matches); i++ {
		match := matches[i]
		match.TournamentID = tournamentID
		match.PlayoffTier = tierMatchCnt
		match.PlayoffMatchNumber = i + 1
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
