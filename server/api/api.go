package api

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/database"
)

type apiContext struct {
	echo.Context
	db *database.Database
}

func getTeams(c echo.Context) error {
	ac := c.(*apiContext)
	teams, err := ac.db.AllTeams()
	if err != nil {
		log.Print("Error getting all teams", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, teams)
}

func saveTeam(c echo.Context) error {
	ac := c.(*apiContext)
	team := new(core.Team)
	if err := c.Bind(team); err != nil {
		return err
	}
	log.Printf("About to save team '%s'", team.Name)
	if err := ac.db.SaveTeam(*team); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

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

	e.GET("/teams", getTeams)
	e.POST("/team", saveTeam)
}
