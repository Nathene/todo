package account

import (
	"net/http"
	"time"
	"todo/internal/db"
	"todo/internal/parser"

	"github.com/labstack/echo/v4"
)

func DeletePage() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Ensure the user is logged in
		user, ok := c.Get("user").(parser.User)
		if !ok {
			return c.Redirect(http.StatusUnauthorized, "/login")
		}

		return c.Render(http.StatusOK, "delete_account/delete_account.gohtml", map[string]interface{}{
			"username": user.Username,
		})
	}
}

func DeleteHandler(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Retrieve the user details from context
		user, ok := c.Get("user").(parser.User)
		if !ok {
			return c.Redirect(http.StatusUnauthorized, "/login")
		}

		// Delete user from the database
		_, err := db.Exec("DELETE FROM users WHERE username = ?", user.Username)
		if err != nil {
			c.Logger().Errorf("Failed to delete account for user '%s': %v", user.Username, err)
			return c.Render(http.StatusInternalServerError, "delete_account/delete_account.gohtml", map[string]interface{}{
				"Error": "Failed to delete account.",
			})
		}

		// Clear the token cookie
		c.SetCookie(&http.Cookie{
			Name:     "token",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour), // Set expiry to past
			HttpOnly: true,
			Path:     "/",
		})

		// Redirect to the home page
		return c.Redirect(http.StatusSeeOther, "/")
	}
}
