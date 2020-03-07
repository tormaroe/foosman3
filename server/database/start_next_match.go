package database

import (
	"fmt"
	"log"
	"sync"

	"github.com/tormaroe/foosman3/server/core"
)

func NewStartMatchChan() chan *core.StartNextMatchRequest {
	startMatchChan := make(chan *core.StartNextMatchRequest, 0)
	go func() {
		for {
			req := <-startMatchChan
			log.Println("Start next match request")

			doNextMatch(req)

			req.WG.Done()
		}
	}()
	return startMatchChan
}

func doNextMatch(req *core.StartNextMatchRequest) {
	var match Match
	queryResult := req.FoosmanContext.DB.Where(
		"tournament_id = ? and state = ?",
		req.TournamentID,
		int(core.Scheduled),
	).Order(
		"sequence asc",
	).First(&match)

	if queryResult.RecordNotFound() {
		log.Println("No more scheduled matches")

		var tournament Tournament
		if err := req.FoosmanContext.DB.First(&tournament, req.TournamentID).Error; err != nil {
			log.Printf("Error getting tournament")
			return
		}

		if tournament.State == int(core.GroupPlayStarted) {
			tournament.State = int(core.GroupPlayDone)
			if err := req.FoosmanContext.DB.Save(&tournament).Error; err != nil {
				log.Printf("ERROR setting tournament to GroupPlayDone")
			}
		} else if tournament.State == int(core.EliminationPlayStarted) {
			advanceElimination(req, tournament)
		}

		return
	}

	if queryResult.Error != nil {
		log.Fatal(queryResult.Error)
	}

	// TODO: I believe the following test will never be true
	if match.ID < 1 {
		log.Println("No more scheduled matches")
		return
	}

	match.State = int(core.InProgress)
	match.Table = req.Table
	log.Printf("Starting match %d on table %s", match.ID, match.Table)
	if err := req.FoosmanContext.DB.Save(match).Error; err != nil {
		log.Fatal(err)
	}

	done := ScheduleUpcoming(req.FoosmanContext, req.TournamentID, 1)
	done.Wait()
}

func advanceElimination(req *core.StartNextMatchRequest, tournament Tournament) {
	log.Println("ADVANCE ELIMINATION")

	var ongoingCnt int
	if err := req.FoosmanContext.DB.Table("matches").Where(
		"tournament_id = ? and state = ?",
		req.TournamentID,
		int(core.InProgress),
	).Count(&ongoingCnt).Error; err != nil {
		return
	}

	if ongoingCnt > 0 {
		log.Printf("%d matches ongoing, no new matches to generate yet", ongoingCnt)
		return
	}

	var prevPlayoffTier []struct{ Value int }
	if err := req.FoosmanContext.DB.Raw(
		`
		SELECT min(playoff_tier) as value FROM matches
		WHERE tournament_id = ? and group_id = 0
		`,
		req.TournamentID,
	).Scan(&prevPlayoffTier).Error; err != nil {
		return
	}
	log.Printf("prevPlayoffTier == %d", prevPlayoffTier[0].Value)

	if prevPlayoffTier[0].Value == 0 {
		panic("prevPlayoffTier == 0")
	}

	if prevPlayoffTier[0].Value == 1 {
		// Final played
		tournament.State = int(core.Done)
		if err := req.FoosmanContext.DB.Save(&tournament).Error; err != nil {
			log.Printf("ERROR setting tournament to Done")
		}
		return
	}

	var matches []Match
	if err := req.FoosmanContext.DB.Preload("MatchResults").Where(
		"tournament_id = ? and group_id = 0 and playoff_tier = ?",
		req.TournamentID,
		prevPlayoffTier[0].Value,
	).Order("playoff_match_number").Find(&matches).Error; err != nil {
		return
	}

	var maxSequence []struct{ Value int }
	if err := req.FoosmanContext.DB.Raw(
		"SELECT MAX(sequence) as value FROM matches WHERE tournament_id = ?",
		req.TournamentID,
	).Scan(&maxSequence).Error; err != nil {
		return
	}

	nextPlayoffTier := prevPlayoffTier[0].Value / 2

	cnt := 0
	for i := 0; i < len(matches); i = i + 2 {
		cnt++
		m1 := matches[i]
		m2 := matches[i+1]
		w1 := m1.Team1ID
		w2 := m2.Team1ID
		if m1.MatchResults[1].Win > 0 {
			w1 = m1.Team2ID
		}
		if m2.MatchResults[1].Win > 0 {
			w2 = m2.Team2ID
		}

		nextMatch := Match{
			Team1ID:            w1,
			Team2ID:            w2,
			TournamentID:       req.TournamentID,
			PlayoffTier:        nextPlayoffTier,
			PlayoffMatchNumber: cnt,
			State:              int(core.Scheduled),
			Sequence:           maxSequence[0].Value + cnt,
		}

		if cnt <= tournament.TableCount {
			nextMatch.State = int(core.InProgress)
			nextMatch.Table = fmt.Sprintf("Table %d", cnt)
		}

		if err := req.FoosmanContext.DB.Create(&nextMatch).Error; err != nil {
			return
		}

		if err := req.FoosmanContext.DB.Create(&MatchResult{
			TeamID:  nextMatch.Team1ID,
			MatchID: nextMatch.ID,
			Points:  0,
			Win:     0,
			Loss:    0,
			Draw:    0,
		}).Error; err != nil {
			return
		}

		if err := req.FoosmanContext.DB.Create(&MatchResult{
			TeamID:  nextMatch.Team2ID,
			MatchID: nextMatch.ID,
			Points:  0,
			Win:     0,
			Loss:    0,
			Draw:    0,
		}).Error; err != nil {
			return
		}
	}
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func StartNextMatch(c *core.FoosmanContext, tournamentID int, table string) *sync.WaitGroup {
	req := core.StartNextMatchRequest{
		TournamentID:   tournamentID,
		Table:          table,
		FoosmanContext: c,
		WG:             sync.WaitGroup{},
	}
	req.WG.Add(1)
	c.StartNextMatchChan <- &req
	return &req.WG
}
