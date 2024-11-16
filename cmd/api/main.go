package main

import (
	"log"
	"todo/internal/api/routes"
	"todo/internal/db"

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
	defer db.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	routes.Handle(e, db)

	e.Logger.Fatal(e.Start(":" + port))
}
