package core

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/database"
)

type FoosmanContext struct {
	echo.Context
	DB *gorm.DB
}

func (ac *FoosmanContext) GetParamID() (int, error) {
	p := ac.Param("id")
	tID, err := strconv.Atoi(p)
	if err != nil {
		log.Print("Unable to parse id route param", err)
		return 0, ac.NoContent(http.StatusBadRequest)
	}
	return tID, nil
}

// AssertTournamentNotStarted will return an error if tournament has
// started.
func (ac *FoosmanContext) AssertTournamentNotStarted(ID int) error {
	var tournament database.Tournament
	if err := ac.DB.First(&tournament, ID).Error; err != nil {
		return err
	}

	if tournament.State != int(New) {
		return errors.New("Can't delete team from a tournament that has started")
	}

	return nil
}
