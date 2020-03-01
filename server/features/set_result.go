package features

import (
	"errors"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/database"
)

type setResultRequest struct {
	MatchID  int  `json:"matchId"`
	IsDraw   bool `json:"isDraw"`
	WinnerID int  `json:"winnerId"`
}

func SetResult(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	tournamentID, err := ac.GetParamID()
	if err != nil {
		return err
	}

	// TODO: Assert tournament started ?

	req := new(setResultRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	ac.SetResultMux.Lock() // TODO: Am I using this everywhere it's needed?
	defer ac.SetResultMux.Unlock()

	setResult(ac, tournamentID, *req)

	return c.NoContent(http.StatusOK)
}

func setResult(c *core.FoosmanContext, tournamentID int, req setResultRequest) error {

	log.Printf("Register results for match %d", req.MatchID)

	// Load match
	var match database.Match
	if err := c.DB.First(&match, req.MatchID).Error; err != nil {
		return err
	}
	var results []database.MatchResult
	if err := c.DB.Where("match_id = ?", req.MatchID).Find(&results).Error; err != nil {
		return err
	}

	// Set result and save
	if len(results) != 2 {
		return errors.New("Match has != 2 results")
	}

	if req.IsDraw {
		results[0].Points = 1
		results[1].Points = 1
		results[0].Draw = 1
		results[1].Draw = 1
		results[0].Win = 0
		results[1].Win = 0
		results[0].Loss = 0
		results[1].Loss = 0
	} else {
		for i := 0; i < 2; i++ {
			if results[i].TeamID == req.WinnerID {
				results[i].Points = 2
				results[i].Draw = 0
				results[i].Win = 1
				results[i].Loss = 0
			} else {
				results[i].Points = 0
				results[i].Draw = 0
				results[i].Win = 0
				results[i].Loss = 1
			}
		}
	}

	matchWasAlreadyPlayed := match.State == int(core.Played)

	if err := c.DB.Transaction(func(tx *gorm.DB) error {
		if !matchWasAlreadyPlayed {
			match.State = int(core.Played)
			if err := tx.Save(match).Error; err != nil {
				return err
			}
		}
		for i := 0; i < len(results); i++ {
			if err := tx.Save(results[i]).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	if !matchWasAlreadyPlayed {
		// Trigger next match on same table
		done := database.StartNextMatch(c, tournamentID, match.Table)
		done.Wait()
	}

	return nil
}
