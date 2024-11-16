package controller

import (
	"net/http"
	"time"
	"todo/internal/service"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var jwtSecret = []byte("your-secret-key") // Replace with a secure secret key

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JWTClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Figure out how i can use this type, instead of having a million functions everywhere.
type User struct {
	userName string
	lists    []service.TodoList
	groups   []service.TodoGroup
}

func Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req LoginRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		// Replace this with database validation
		if req.Username != "admin" || req.Password != "password" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
		}

		// Generate JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
			Username: req.Username,
			StandardClaims: jwt.StandardClaims{
				// FIXME change this, adding infinite expiry for testing
				ExpiresAt: time.Now().Add(time.Hour * 100000).Unix(), // Token expires in 1 hour
			},
		})
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		}

		return c.JSON(http.StatusOK, map[string]string{"token": tokenString})
	}
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing or invalid token"})
		}

		token, err := jwt.ParseWithClaims(authHeader, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid claims"})
		}

		c.Set("username", claims.Username)
		return next(c)
	}
}
