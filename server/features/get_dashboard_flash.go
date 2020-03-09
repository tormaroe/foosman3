package features

import (
	"math/rand"
	"net/http"

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

		// if no match results yet

		// if only one match result

		// XXX is in the lead in group Y

		// XXX have won group Y with Z wins

		// X matches has been played so far

		// Group play progress: X%

	} else if tournament.State == int(core.GroupPlayDone) {

	} else if tournament.State == int(core.EliminationPlayStarted) {

	} else if tournament.State == int(core.Done) {

	}

	return c.JSONPretty(http.StatusOK, res, "  ")
}
