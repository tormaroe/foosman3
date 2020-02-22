package database

import (
	"log"
	"sync"

	"github.com/tormaroe/foosman3/server/core"
)

func NewScheduleChan() chan *core.ScheduleRequest {
	scheduleChan := make(chan *core.ScheduleRequest, 0)
	go func() {
		for {
			req := <-scheduleChan
			log.Println("SCHEDULE UPCOMING MATCH!!!!!!")

			doSchedule(req)

			// TODO: How to start scheduled matches? (done elsewhere)

			req.WG.Done()
		}
	}()
	return scheduleChan
}

func doSchedule(req *core.ScheduleRequest) {

	// TODO: Make a TEST

	// Get all teams
	type Team struct {
		ID      int
		GroupID int
		Weight  int
	}
	var teams []Team
	err := req.DB.Raw("SELECT id, group_id FROM teams WHERE tournament_id = ?", req.TournamentID).Scan(&teams).Error
	if err != nil {
		goto done
	}

	// Get all matches

	// If no more matches to schedule, end
	// Give teams weight based on recently played matches (sequence)
	// Give teams weight based on recently played groups
	// Pick any of the teams (up to count) with the lowest weight

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
