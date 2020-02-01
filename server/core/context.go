package core

import (
	"log"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
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
