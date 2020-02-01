package features

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/database"
)

type groupDefinition struct {
	Name    string `json:"name"`
	TeamIDs []int  `json:"teams"`
}

// SetGroups defines the groups for a Tournament.
// It should be provided a list of groups. Each group
// has a name and a list of Team IDs belonging to that group.
// Any pre-existing group structure for the Tournament in question
// is completely replaced.
func SetGroups(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	tournamentID, err := ac.GetParamID()
	if err != nil {
		return err
	}

	if err := ac.AssertTournamentNotStarted(tournamentID); err != nil {
		return err
	}

	req := new([]groupDefinition)
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := ac.DB.Transaction(func(tx *gorm.DB) error {
		err = deleteAllGroups(tx, tournamentID)
		for _, g := range *req {
			gID, err := makeGroup(tx, tournamentID, g.Name)
			if err != nil {
				return err
			}
			if err = setTeamsGroup(tx, gID, g.TeamIDs); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func deleteAllGroups(tx *gorm.DB, tournamentID int) error {
	err := tx.Model(&database.Team{}).Where("tournament_id = ?", tournamentID).Update("group_id", nil).Error
	if err != nil {
		return err
	}

	return tx.Delete(&database.Group{}, "tournament_id", tournamentID).Error
}

func makeGroup(tx *gorm.DB, tournamentID int, name string) (int, error) {
	var tournament database.Tournament
	if err := tx.First(&tournament, tournamentID).Error; err != nil {
		return 0, err
	}
	group := database.Group{
		Name:       name,
		Tournament: tournament,
	}
	err := tx.Create(&group).Error
	return group.ID, err
}

func setTeamsGroup(tx *gorm.DB, groupID int, teamIDs []int) error {
	return tx.Model(&database.Team{}).Where("id IN (?)", teamIDs).Update("group_id", groupID).Error
}
