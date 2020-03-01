package database

import "github.com/tormaroe/foosman3/server/core"
import "log"
import "sync"

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
		"state = ?", int(core.Scheduled),
	).Order(
		"sequence asc",
	).First(&match)

	if queryResult.RecordNotFound() {
		log.Println("No more scheduled matches")
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
