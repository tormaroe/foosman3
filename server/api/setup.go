package api

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/features"
)

// Init adds middlewares and all API routes to Echo
func Init(e *echo.Echo, d *gorm.DB) {

	// Middleware

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &core.FoosmanContext{c, d}
			return next(cc)
		}
	})

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	e.GET("/tournaments", features.GetTournaments)
	e.POST("/tournaments", features.AddTournament)
	e.PATCH("/tournaments", features.UpdateTournament)

	e.GET("/tournaments/:id", features.GetTournamentByID)
	e.POST("/tournaments/:id/teams", features.AddTeam)
	e.PATCH("/tournaments/teams", features.UpdateTeam)
	e.DELETE("/tournaments/teams/:id", features.DeleteTeam)
	e.POST("/tournaments/:id/groups", features.SetGroups)

	e.POST("/tournaments/:id/generate-matches", features.GenerateMatches)
	e.GET("/tournaments/:id/matches", features.GetTournamentMatches)

	e.GET("/teams/:id", features.GetTeam)
}
