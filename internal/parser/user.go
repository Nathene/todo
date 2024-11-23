package parser

import (
	"github.com/golang-jwt/jwt"
)

type User struct {
	Username   string
	FirstName  string
	GroupID    int
	Email      string
	DarkMode   bool
	IsLoggedIn bool
	Lists      []TodoList
	Groups     []TodoGroup
}

var JwtSecret = []byte("your-secret-key")

type JWTClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type LoginRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}
