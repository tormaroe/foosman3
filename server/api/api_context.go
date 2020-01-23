package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/database"
)

type apiContext struct {
	echo.Context
	db *database.Database
}

func (ac *apiContext) getParamID() (int, error) {
	p := ac.Param("id")
	tID, err := strconv.Atoi(p)
	if err != nil {
		log.Print("Unable to parse id route param", err)
		return 0, ac.NoContent(http.StatusBadRequest)
	}
	return tID, nil
}
