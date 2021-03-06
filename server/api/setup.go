package api

import (
	"net/http"
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/features"
)

// Init adds middlewares and all API routes to Echo
func Init(
	e *echo.Echo,
	d *gorm.DB,
	scheduleChan chan *core.ScheduleRequest,
	startNextMatchChan chan *core.StartNextMatchRequest,
	setResultMux *sync.Mutex,
) {

	// Middleware

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &core.FoosmanContext{
				Context:            c,
				DB:                 d,
				ScheduleChan:       scheduleChan,
				StartNextMatchChan: startNextMatchChan,
				SetResultMux:       setResultMux,
			}
			return next(cc)
		}
	})

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339_nano} ${method} ${uri} status: ${status} ${latency_human} ${error}\n",
	}))
	//e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	e.GET("/tournaments", features.GetTournaments)
	e.POST("/tournaments", features.AddTournament)
	e.PATCH("/tournaments", features.UpdateTournament)

	e.GET("/tournaments/:id", features.GetTournamentByID)
	e.GET("/tournaments/:id/export", features.ExportTournament)
	e.GET("/tournaments/:id/dashboard", features.GetDashboardFlash)
	e.DELETE("/tournaments/:id", features.DeleteTournament)
	e.POST("/tournaments/:id/reset", features.ResetTournament)
	e.POST("/tournaments/:id/teams", features.AddTeam)
	e.PATCH("/tournaments/teams", features.UpdateTeam)
	e.DELETE("/tournaments/teams/:id", features.DeleteTeam)
	e.POST("/tournaments/:id/groups", features.SetGroups)

	e.POST("/tournaments/:id/start", features.StartTournament)
	e.POST("/tournaments/:id/start-elimination", features.StartElimination)
	e.GET("/tournaments/:id/matches", features.GetTournamentMatches)
	e.GET("/tournaments/:id/elimination-matches", features.GetEliminationMatches)
	e.GET("/tournaments/:id/scores", features.GetTournamentScores)
	e.POST("/matches/:id/reset", features.ResetMatch)

	e.GET("/teams/:id", features.GetTeam)

	e.GET("/tournaments/:id/matches/in-progress", features.GetMatchesInProgress)
	e.GET("/tournaments/:id/matches/scheduled", features.GetMatchesScheduled)
	e.POST("/tournaments/:id/match/set-result", features.SetResult)
}
