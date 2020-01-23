package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/tormaroe/foosman3/server/database"
)

// Init adds middlewares and all API routes to Echo
func Init(e *echo.Echo, d *database.Database) {

	// Middleware

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &apiContext{c, d}
			return next(cc)
		}
	})

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	e.GET("/tournaments", getTournaments)
	e.POST("/tournaments", addTournament)
	e.PATCH("/tournaments", updateTournament)

	e.GET("/tournaments/:id/teams", getTournamentTeams)
	e.POST("/tournaments/:id/teams", addTeam)
}
