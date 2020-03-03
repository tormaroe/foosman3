package database

import (
	"log"
	"math"
	"math/rand"
	"sync"

	"github.com/tormaroe/foosman3/server/core"
)

// Refactor to actor pattern https://www.appdynamics.com/blog/engineering/three-productive-go-patterns-put-radar/

func NewScheduleChan() chan *core.ScheduleRequest {
	scheduleChan := make(chan *core.ScheduleRequest, 0)
	go func() {
		for {
			req := <-scheduleChan
			log.Println("Do schedule request")

			doSchedule(req)

			req.WG.Done()
		}
	}()
	return scheduleChan
}

func doSchedule(req *core.ScheduleRequest) {

	// Get all teams
	type Team struct {
		ID      int
		GroupID int
		Weight  int
	}
	var teams []Team
	var matches []Match
	unscheduledCount := 0
	maxSequence := 0

	teamMap := make(map[int]*Team)

	// Get all teams
	err := req.DB.Raw("SELECT id, group_id FROM teams WHERE tournament_id = ?", req.TournamentID).Scan(&teams).Error
	if err != nil {
		goto done
	}
	for i, team := range teams {
		teamMap[team.ID] = &teams[i]
	}

	// Get all matches
	err = req.DB.Where("tournament_id = ?", req.TournamentID).Find(&matches).Error
	if err != nil {
		goto done
	}

	// Give teams weight based on recently played matches (sequence)
	for _, m := range matches {
		if m.Sequence > 0 {
			// m is a match already scheduled.
			// Use the schedule to add weight to teams, making it less
			// likely that they bubble up to the top for next scheduling.
			teamMap[m.Team1ID].Weight += m.Sequence
			teamMap[m.Team2ID].Weight += m.Sequence
			if maxSequence < m.Sequence {
				maxSequence = m.Sequence
			}
		} else {
			unscheduledCount = unscheduledCount + 1
		}
	}

	// If no more matches to schedule, end
	if unscheduledCount == 0 {
		goto done
	}

	// Shuffle matches for some randomness
	// rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(matches), func(i, j int) { matches[i], matches[j] = matches[j], matches[i] })

	// Pick any of the teams (up to count) with the lowest weight
	for i := 0; i < req.Count; i++ {
		bestWeightSoFar := math.MaxInt32
		var bestMatch *Match
		for mi, m := range matches {
			if m.Sequence == 0 {
				combWeight := teamMap[m.Team1ID].Weight + teamMap[m.Team2ID].Weight
				if combWeight < bestWeightSoFar {
					bestWeightSoFar = combWeight
					bestMatch = &matches[mi]
				}
			}
		}
		if bestMatch != nil {
			maxSequence++
			bestMatch.Sequence = maxSequence
			bestMatch.State = int(core.Scheduled)
			log.Printf("Scheduling match %d schedule %d", bestMatch.ID, maxSequence)
			req.DB.Save(bestMatch) // TODO: Check error
			teamMap[bestMatch.Team1ID].Weight += maxSequence
			teamMap[bestMatch.Team2ID].Weight += maxSequence
		}
	}

done:
}

func ScheduleUpcoming(c *core.FoosmanContext, tournamentID int, count int) *sync.WaitGroup {
	req := core.ScheduleRequest{
		TournamentID: tournamentID,
		Count:        count,
		DB:           c.DB,
		WG:           sync.WaitGroup{},
	}
	req.WG.Add(1)
	c.ScheduleChan <- &req
	return &req.WG
}
