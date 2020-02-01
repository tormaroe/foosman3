package features

import (
	"net/http"

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
	req := new([]groupDefinition)
	if err := c.Bind(req); err != nil {
		return err
	}
	// TODO: Transaction
	err = deleteAllGroups(ac, tournamentID)
	for _, g := range *req {
		gID, err := makeGroup(ac, tournamentID, g.Name)
		if err != nil {
			return err
		}
		if err = setTeamsGroup(ac, gID, g.TeamIDs); err != nil {
			return err
		}
	}
	return c.NoContent(http.StatusOK)
}

func deleteAllGroups(d *core.FoosmanContext, tournamentID int) error {
	err := d.DB.Model(&database.Team{}).Where("tournament_id = ?", tournamentID).Update("group_id", nil).Error
	if err != nil {
		return err
	}

	return d.DB.Delete(&database.Group{}, "tournament_id", tournamentID).Error
}

func makeGroup(d *core.FoosmanContext, tournamentID int, name string) (int, error) {
	var tournament database.Tournament
	if err := d.DB.First(&tournament, tournamentID).Error; err != nil {
		return 0, err
	}
	group := database.Group{
		Name:       name,
		Tournament: tournament,
	}
	err := d.DB.Create(&group).Error
	return group.ID, err
}

func setTeamsGroup(d *core.FoosmanContext, groupID int, teamIDs []int) error {
	return d.DB.Model(&database.Team{}).Where("id IN (?)", teamIDs).Update("group_id", groupID).Error
}
