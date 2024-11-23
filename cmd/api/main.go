package main

import (
	"log"
	"os"
	"todo/internal/api/routes"
	"todo/internal/db"
	"todo/internal/pkg/renderer"
	"todo/internal/util"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	port = "8040"
)

func main() {
	db, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer util.Defer(db.Close)

	e := echo.New()

	// Initialize the custom renderer with templates directory
	e.Renderer = renderer.NewRenderer("templates")
	// Serve static files like styles.css
	e.Static("/static", "templates/static")

	e.Logger.SetOutput(os.Stdout) // Ensure logs are printed to stdout
	e.Use(middleware.Logger())    // Enable Echo's request logging middleware

	e.Use(middleware.Recover())

	routes.Handle(e, db)

	e.Logger.Fatal(e.Start(":" + port))
}
