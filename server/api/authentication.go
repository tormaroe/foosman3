package api

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/tormaroe/foosman3/server/core"
	"github.com/tormaroe/foosman3/server/database"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func userCount(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	var count int
	if err := ac.DB.Model(database.User{}).Count(&count).Error; err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]int{
		"count": count,
	})
}

func me(c echo.Context) error {
	user := c.Get("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"username":    claims["name"].(string),
		"scorekeeper": claims["scorekeeper"].(bool),
		"admin":       claims["admin"].(bool),
	})
}

func createFirstAdmin(c echo.Context) error {
	ac := c.(*core.FoosmanContext)
	var count int
	if err := ac.DB.Model(database.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return echo.ErrForbidden
	}

	req := new(loginRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	log.Printf("Creating initial admin user '%s'\n", req.Username)

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return err
	}

	user := database.User{
		Name:          req.Username,
		PasswordHash:  string(hash),
		IsAdmin:       true,
		IsScorekeeper: true,
	}

	if err := ac.DB.Create(&user).Error; err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func createUser(c echo.Context) error {
	return echo.ErrServiceUnavailable // TODO
}

func updateUser(c echo.Context) error {
	return echo.ErrServiceUnavailable // TODO
}

func getUsers(c echo.Context) error {
	return echo.ErrServiceUnavailable // TODO
}

func login(c echo.Context) error {
	ac := c.(*core.FoosmanContext)

	req := new(loginRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	var user database.User
	if err := ac.DB.Where("name = ?", req.Username).First(&user).Error; err != nil {
		log.Println("Error getting requested user")
		return echo.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		log.Println("Bad password")
		return echo.ErrUnauthorized
	}

	exp := time.Now().Add(72 * time.Hour)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Name
	claims["admin"] = user.IsAdmin
	claims["scorekeeper"] = user.IsScorekeeper
	claims["exp"] = exp

	// Generate encoded token
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = t
	cookie.Expires = exp
	c.SetCookie(cookie)

	return c.String(http.StatusOK, "logged in")
}
