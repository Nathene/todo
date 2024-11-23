package account

import (
	"net/http"
	"todo/internal/db"
	"todo/internal/parser"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Update First Name
func UpdateFirstName(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		firstName := c.FormValue("firstname")
		username := c.Get("username").(string)

		if firstName == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "First name cannot be empty"})
		}

		// Capitalize the first name
		firstName = cases.Title(language.English).String(firstName)

		_, err := db.Exec("UPDATE users SET firstname = ? WHERE username = ?", firstName, username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update first name"})
		}

		return c.Redirect(http.StatusSeeOther, "/settings")
	}
}

// Update Username
func UpdateUsername(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		newUsername := c.FormValue("username")
		username := c.Get("username").(string)

		if newUsername == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username cannot be empty"})
		}

		// Check if the new username already exists
		var existingUsername string
		err := db.QueryRow("SELECT username FROM users WHERE username = ?", newUsername).Scan(&existingUsername)
		if err == nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Username already exists"})
		}

		_, err = db.Exec("UPDATE users SET username = ? WHERE username = ?", newUsername, username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update username"})
		}

		return c.Redirect(http.StatusSeeOther, "/settings")
	}
}

// Update Email
func UpdateEmail(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		newEmail := c.FormValue("email")
		username := c.Get("username").(string)

		if newEmail == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email cannot be empty"})
		}

		// Check if the new email already exists
		var existingEmail string
		err := db.QueryRow("SELECT email FROM users WHERE email = ?", newEmail).Scan(&existingEmail)
		if err == nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Email already exists"})
		}

		_, err = db.Exec("UPDATE users SET email = ? WHERE username = ?", newEmail, username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update email"})
		}

		return c.Redirect(http.StatusSeeOther, "/settings")
	}
}

// Update Password
func UpdatePassword(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		currentPassword := c.FormValue("currentPassword")
		newPassword := c.FormValue("newPassword")
		username := c.Get("username").(string)

		if currentPassword == "" || newPassword == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Both current and new passwords are required"})
		}

		// Retrieve the current password hash
		var hashedPassword string
		err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve current password"})
		}

		// Check if the current password is correct
		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currentPassword)); err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Current password is incorrect"})
		}

		// Hash the new password
		newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash new password"})
		}

		_, err = db.Exec("UPDATE users SET password = ? WHERE username = ?", string(newHashedPassword), username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update password"})
		}

		return c.Redirect(http.StatusSeeOther, "/settings")
	}
}

// Update Dark Mode
func UpdateDarkMode(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Safely retrieve the username from the context
		user, ok := c.Get("user").(parser.User)
		if !ok {
			c.Logger().Error("Username not found in context")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized access"})
		}

		// Get the dark mode toggle value from the form
		darkMode := c.FormValue("darkMode") == "on"

		// Update the dark_mode setting in the database
		_, err := db.Exec("UPDATE users SET darkmode = ? WHERE username = ?", darkMode, user.Username)
		if err != nil {
			c.Logger().Errorf("Failed to update dark mode for user %s: %v", user.Username, err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update dark mode setting"})
		}

		c.Logger().Infof("Dark mode updated for user %s: %v", user.Username, darkMode)

		// Redirect back to the settings page
		return c.Redirect(http.StatusSeeOther, "/settings")
	}
}
