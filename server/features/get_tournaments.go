package features

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/database"
)

// GetTournaments responds to a GET request for all tournaments
func GetTournaments(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	lst, err := getTournaments(ac.DB)
	if err != nil {
		log.Print("Error getting tournaments", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSONPretty(http.StatusOK, lst, "  ")
}

func getTournaments(db *gorm.DB) ([]database.Tournament, error) {
	var result []database.Tournament
	err := db.Select("id, name, table_count, state").Find(&result).Error
	return result, err
}
