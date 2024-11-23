package controller

import (
	"net/http"
	"todo/internal/parser"

	"github.com/labstack/echo/v4"
)

func LandingPage() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, _ := c.Get("user").(parser.User)

		// Render the landing page with the authentication and darkMode state
		return c.Render(http.StatusOK, "landing/landing.gohtml", map[string]interface{}{
			"title":      "Dashboard",
			"isLoggedIn": user.IsLoggedIn,
			"username":   user.Username,
			"darkMode":   user.DarkMode,
			"user":       user,
		})
	}
}
