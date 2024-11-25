package controller

import (
	"log"
	"net/http"
	"time"
	"todo/internal/db"
	"todo/internal/parser"
	"todo/internal/util"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Figure out how i can use this type, instead of having a million functions everywhere.

func Login(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method == http.MethodGet {
			return c.Render(http.StatusOK, "login/login.gohtml", nil)
		}

		// Handle POST (login form submission)
		var req parser.LoginRequest
		if err := c.Bind(&req); err != nil {
			return c.Render(http.StatusBadRequest, "login/login.gohtml", map[string]interface{}{
				"Error": "Invalid request",
			})
		}

		// Validate credentials
		var user parser.User
		var hashedPassword string
		err := db.QueryRow(
			"SELECT username, firstname, email, password FROM users WHERE username = ?",
			req.Username,
		).Scan(&user.Username, &user.FirstName, &user.Email, &hashedPassword)
		if err != nil {
			// Username not found
			return c.Render(http.StatusUnauthorized, "login/login.gohtml", map[string]interface{}{
				"Error": "Invalid username or password",
			})
		}

		// Compare the provided password with the hashed password
		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
			// Password mismatch
			return c.Render(http.StatusUnauthorized, "login/login.gohtml", map[string]interface{}{
				"Error": "Invalid username or password",
			})
		}

		// Generate JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, parser.JWTClaims{
			Username: user.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			},
		})
		tokenString, err := token.SignedString(parser.JwtSecret)
		if err != nil {
			return c.Render(http.StatusInternalServerError, "login/login.gohtml", map[string]interface{}{
				"Error": "Failed to generate token",
			})
		}

		// Set cookie and context
		c.SetCookie(&http.Cookie{
			Name:     "token",
			Value:    tokenString,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Path:     "/",
		})
		c.Set("user", user)

		return c.Redirect(http.StatusSeeOther, "/")
	}
}

func Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Clear the token by setting an expired cookie
		c.SetCookie(&http.Cookie{
			Name:     "token",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour), // Set expiry to past
			HttpOnly: true,
			Path:     "/",
		})

		// Redirect to the landing page
		return c.Redirect(http.StatusSeeOther, "/")
	}
}

func AuthMiddleware(db *db.Database) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log.Println("AuthMiddleware executed")
			// Default user context
			c.Set("user", parser.User{IsLoggedIn: false})

			// Check for the token cookie
			cookie, err := c.Cookie("token")
			if err != nil || cookie.Value == "" {
				c.Logger().Info("No token found or invalid token")
				return next(c) // Proceed without authentication
			}

			// Parse and validate the token
			token, err := jwt.ParseWithClaims(cookie.Value, &parser.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				return parser.JwtSecret, nil
			})
			if err != nil || !token.Valid {
				c.Logger().Info("Invalid token")
				return next(c) // Proceed without authentication
			}

			// Extract claims
			claims, ok := token.Claims.(*parser.JWTClaims)
			if !ok || claims.Username == "" {
				c.Logger().Info("Invalid token claims")
				return next(c)
			}

			// Fetch user details from the database
			var user parser.User
			err = db.QueryRow("SELECT firstname, username, email, darkmode FROM users WHERE username = ?", claims.Username).
				Scan(&user.FirstName, &user.Username, &user.Email, &user.DarkMode)
			if err != nil {
				c.Logger().Errorf("Failed to fetch user details: %v", err)
				return next(c) // Proceed without authentication
			}
			user.IsLoggedIn = true

			// Set the user in the context
			c.Set("user", user)
			c.Set("username", user.Username)
			c.Set("darkMode", user.DarkMode)

			c.Logger().Infof("Authenticated user: %s", user.Username)
			return next(c)
		}
	}
}

// ToggleDarkMode toggles the user's dark mode preference.
func ToggleDarkMode(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Require user to be logged in
		if err := util.RequireLogin(c); err != nil {
			return err
		}

		user, _ := util.GetUserFromContext(c) // Safe to assume user exists after RequireLogin
		user.DarkMode = !user.DarkMode

		_, err := db.Exec("UPDATE users SET darkmode = ? WHERE username = ?", user.DarkMode, user.Username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update dark mode setting"})
		}

		referer := c.Request().Header.Get("Referer")
		if referer == "" {
			referer = "/" // Default to homepage
		}
		return c.Redirect(http.StatusSeeOther, referer)
	}
}
