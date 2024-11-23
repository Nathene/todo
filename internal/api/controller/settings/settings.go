package settings

import (
	"net/http"
	"todo/internal/db"
	"todo/internal/parser"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Page() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the user from the context
		user, ok := c.Get("user").(parser.User)
		if !ok {
			return c.Redirect(http.StatusUnauthorized, "/login")
		}

		// Render the settings page with user details
		return c.Render(http.StatusOK, "settings/settings.gohtml", map[string]interface{}{
			"title":      "Account Settings",
			"isLoggedIn": user.IsLoggedIn,
			"username":   user.Username,
			"darkMode":   user.DarkMode,
			"firstname":  user.FirstName,
			"email":      user.Email,
		})
	}
}

func Update(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Ensure the user is logged in
		user, ok := c.Get("user").(parser.User)
		if !ok {
			return c.Redirect(http.StatusUnauthorized, "/login")
		}

		// Parse form data
		firstName := c.FormValue("firstname")
		email := c.FormValue("email")
		currentPassword := c.FormValue("currentPassword")
		newPassword := c.FormValue("newPassword")
		darkMode := c.FormValue("darkMode") == "on"

		// Update profile and dark mode preference
		_, err := db.Exec("UPDATE users SET firstname = ?, email = ?, darkmode = ? WHERE username = ?", firstName, email, darkMode, user.Username)
		if err != nil {
			c.Logger().Errorf("Failed to update user info: %v", err)
			return c.Render(http.StatusInternalServerError, "settings/settings.gohtml", map[string]interface{}{
				"Error": "Failed to update account.",
			})
		}

		// Update password if provided
		if newPassword != "" {
			// Verify current password
			var storedPassword string
			err := db.QueryRow("SELECT password FROM users WHERE username = ?", user.Username).Scan(&storedPassword)
			if err != nil || bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(currentPassword)) != nil {
				return c.Render(http.StatusUnauthorized, "settings/settings.gohtml", map[string]interface{}{
					"Error": "Current password is incorrect.",
				})
			}

			// Hash and update new password
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
			_, err = db.Exec("UPDATE users SET password = ? WHERE username = ?", string(hashedPassword), user.Username)
			if err != nil {
				c.Logger().Errorf("Failed to update password: %v", err)
				return c.Render(http.StatusInternalServerError, "settings/settings.gohtml", map[string]interface{}{
					"Error": "Failed to update password.",
				})
			}
		}

		return c.Redirect(http.StatusSeeOther, "/settings")
	}
}
