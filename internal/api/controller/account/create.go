package account

import (
	"net/http"
	"time"
	"todo/internal/db"
	"todo/internal/parser"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Display the "Create Account" page
func CreatePage() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "create_account/create_account.gohtml", nil)
	}
}

// Handle form submission to create an account
func CreateHandler(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Handle form submission
		var req struct {
			Username  string `form:"username"`
			Password  string `form:"password"`
			FirstName string `form:"firstname"`
			Email     string `form:"email"`
		}
		if err := c.Bind(&req); err != nil {
			return c.Render(http.StatusBadRequest, "create_account/create_account.gohtml", map[string]interface{}{
				"Error": "Invalid form data",
			})
		}

		// Check if username or email already exists
		var exists int
		err := db.QueryRow(
			"SELECT COUNT(*) FROM users WHERE username = ? OR email = ?",
			req.Username, req.Email,
		).Scan(&exists)
		if err != nil || exists > 0 {
			return c.Render(http.StatusBadRequest, "create_account/create_account.gohtml", map[string]interface{}{
				"Error": "Username or email already exists",
			})
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Render(http.StatusInternalServerError, "create_account/create_account.gohtml", map[string]interface{}{
				"Error": "Failed to create account",
			})
		}

		// Insert the user into the database
		_, err = db.Exec(
			"INSERT INTO users (username, password, firstname, email) VALUES (?, ?, ?, ?)",
			req.Username, string(hashedPassword), req.FirstName, req.Email,
		)
		if err != nil {
			return c.Render(http.StatusInternalServerError, "create_account/create_account.gohtml", map[string]interface{}{
				"Error": "Failed to create account",
			})
		}

		// Automatically log in the user
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, parser.JWTClaims{
			Username: req.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			},
		})
		tokenString, err := token.SignedString(parser.JwtSecret)
		if err != nil {
			return c.Render(http.StatusInternalServerError, "create_account/create_account.gohtml", map[string]interface{}{
				"Error": "Failed to log in",
			})
		}

		// Set the token cookie
		c.SetCookie(&http.Cookie{
			Name:     "token",
			Value:    tokenString,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Path:     "/",
		})

		// Redirect to homepage
		return c.Redirect(http.StatusSeeOther, "/")
	}
}
