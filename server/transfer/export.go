package transfer

// Check out https://xuri.me/excelize/en/ for xlsx generation docs

import (
	"bytes"
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/jinzhu/gorm"
	"github.com/tormaroe/foosman3/server/database"
)

const teamsSheet = "Teams"
const matchesSheet = "Matches"
const resultTeam1Win = "Team1"
const resultTeam2Win = "Team2"
const resultDraw = "Draw"

var matchStates = []string{"Planned", "Scheduled", "In progress", "Played"}

func ExportTournament(db *gorm.DB, ID int) (*bytes.Buffer, error) {

	var t database.Tournament
	if err := db.Preload("Teams").Preload("Groups").First(&t, ID).Error; err != nil {
		return nil, err
	}

	var matches []database.Match
	if err := db.Preload("MatchResults").Preload("Team1").Preload("Team2").Preload("Group").Where(
		"tournament_id = ?",
		ID,
	).Order("sequence desc").Find(&matches).Error; err != nil {
		return nil, err
	}

	f := excelize.NewFile()

	f.SetSheetName("Sheet1", teamsSheet)

	f.SetCellValue(teamsSheet, "A1", "ID")
	f.SetCellValue(teamsSheet, "B1", "Name")
	f.SetCellValue(teamsSheet, "C1", "Player1")
	f.SetCellValue(teamsSheet, "D1", "Player2")
	f.SetCellValue(teamsSheet, "E1", "Player3")
	f.SetCellValue(teamsSheet, "F1", "GroupID")

	for i, team := range t.Teams {
		f.SetCellValue(teamsSheet, fmt.Sprintf("A%d", i+2), team.ID)
		f.SetCellValue(teamsSheet, fmt.Sprintf("B%d", i+2), team.Name)
		f.SetCellValue(teamsSheet, fmt.Sprintf("C%d", i+2), team.Player1)
		f.SetCellValue(teamsSheet, fmt.Sprintf("D%d", i+2), team.Player2)
		f.SetCellValue(teamsSheet, fmt.Sprintf("E%d", i+2), team.Player3)
		f.SetCellValue(teamsSheet, fmt.Sprintf("F%d", i+2), team.GroupID)
	}

	// TODO: Load and join scores on teamsSheet

	f.NewSheet(matchesSheet)

	f.SetCellValue(matchesSheet, "A1", "ID")
	f.SetCellValue(matchesSheet, "B1", "Sequence")
	f.SetCellValue(matchesSheet, "C1", "State")
	f.SetCellValue(matchesSheet, "D1", "State desc")
	f.SetCellValue(matchesSheet, "E1", "Table")
	f.SetCellValue(matchesSheet, "F1", "PlayoffTier")
	f.SetCellValue(matchesSheet, "G1", "PlayoffMatchNumber")
	f.SetCellValue(matchesSheet, "H1", "GroupID")
	f.SetCellValue(matchesSheet, "I1", "GroupName")
	f.SetCellValue(matchesSheet, "J1", "Team1ID")
	f.SetCellValue(matchesSheet, "K1", "Team1Name")
	f.SetCellValue(matchesSheet, "L1", "Team2ID")
	f.SetCellValue(matchesSheet, "M1", "Team2Name")
	f.SetCellValue(matchesSheet, "N1", "Result")

	for i, m := range matches {
		f.SetCellValue(matchesSheet, fmt.Sprintf("A%d", i+2), m.ID)
		f.SetCellValue(matchesSheet, fmt.Sprintf("B%d", i+2), m.Sequence)
		f.SetCellValue(matchesSheet, fmt.Sprintf("C%d", i+2), m.State)
		f.SetCellValue(matchesSheet, fmt.Sprintf("D%d", i+2), matchStates[m.State])
		f.SetCellValue(matchesSheet, fmt.Sprintf("E%d", i+2), m.Table)
		f.SetCellValue(matchesSheet, fmt.Sprintf("F%d", i+2), m.PlayoffTier)
		f.SetCellValue(matchesSheet, fmt.Sprintf("G%d", i+2), m.PlayoffMatchNumber)
		f.SetCellValue(matchesSheet, fmt.Sprintf("H%d", i+2), m.Group.ID)
		f.SetCellValue(matchesSheet, fmt.Sprintf("I%d", i+2), m.Group.Name)
		f.SetCellValue(matchesSheet, fmt.Sprintf("J%d", i+2), m.Team1.ID)
		f.SetCellValue(matchesSheet, fmt.Sprintf("K%d", i+2), m.Team1.Name)
		f.SetCellValue(matchesSheet, fmt.Sprintf("L%d", i+2), m.Team2.ID)
		f.SetCellValue(matchesSheet, fmt.Sprintf("M%d", i+2), m.Team2.Name)
		f.SetCellValue(matchesSheet, fmt.Sprintf("N%d", i+2), getResult(m))
	}

	buf, err := f.WriteToBuffer()
	return buf, err
}

func getResult(m database.Match) string {
	if m.MatchResults[0].Draw == 1 {
		return resultDraw
	}
	if m.Team1ID == m.MatchResults[0].TeamID {
		if m.MatchResults[0].Win == 1 {
			return resultTeam1Win
		}
		return resultTeam2Win
	}
	if m.MatchResults[0].Win == 1 {
		return resultTeam2Win
	}
	return resultTeam1Win
}
