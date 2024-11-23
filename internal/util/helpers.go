package util

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Defer wraps a resource close function with error handling.
func Defer(f func() error) {
	err := f()
	if err != nil {
		log.Printf("Failed to defer resource: %v", err)
	}
}

func Capitalize(name string) string {
	// Create a caser for the title case transformation
	caser := cases.Title(language.English)
	return caser.String(name)
}

func StatusColor(status string) string {
	switch status {
	case "Not Started":
		return "bg-secondary text-white" // Grey
	case "In Progress":
		return "bg-primary text-white" // Blue
	case "Completed":
		return "bg-success text-white" // Green
	default:
		return "bg-light text-dark" // Default
	}
}

func BackPage(c echo.Context) error {
	// Redirect to the previous page

	referer := c.Request().Header.Get("Referer")
	if referer == "" {
		referer = "/" // Default to the homepage if no referer
	}
	return c.Redirect(http.StatusSeeOther, referer)
}
