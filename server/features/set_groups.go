package features

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
)

type groupDefinition struct {
	Name    string `json:"name"`
	TeamIDs []int  `json:"teams"`
}

type setGroupsRequest struct {
	Groups []groupDefinition `json:"groups"`
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
	req := new(setGroupsRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	// TODO: Transaction
	err = deleteAllGroups(ac, tournamentID)
	for _, g := range req.Groups {
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
	stmt, err := d.DB.Prepare(`
		update team
		set group_id = null
		where tournament_id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(tournamentID)
	return err
}

func makeGroup(d *core.FoosmanContext, tournamentID int, name string) (int, error) {
	//TODO
	return 1, nil
}

func setTeamsGroup(d *core.FoosmanContext, groupID int, teamIDs []int) error {
	stmt, err := d.DB.Prepare(`
		update team
		set group_id = ?
		where id in (?` + strings.Repeat(",?", len(teamIDs)-1) + ")")
	if err != nil {
		return err
	}
	defer stmt.Close()
	args := []interface{}{groupID}
	teamIDsSlice := make([]interface{}, len(teamIDs))
	args = append(args, teamIDsSlice...)

	_, err = stmt.Exec(args...)
	return err
}
