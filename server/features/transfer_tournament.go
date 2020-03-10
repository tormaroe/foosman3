package features

import (
	"bytes"
	"net/http"

	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/transfer"
)

func ExportTournament(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	tournamentID, err := ac.GetParamID()
	if err != nil {
		return err
	}

	buf, err := transfer.ExportTournament(ac.DB, tournamentID)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(buf.Bytes())
	return c.Stream(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", reader)
}
