package calendar

import (
	"net/http"
	"time"
	"todo/internal/db"
	"todo/internal/parser"
	"todo/internal/util"

	"github.com/labstack/echo/v4"
)

// func AddPage(db *db.Database) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		rows, err := db.Query("SELECT id, name, description, event_date FROM events ORDER BY event_date ASC")
// 		if err != nil {
// 			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch events"})
// 		}
// 		defer rows.Close()

// 		var events []map[string]interface{}
// 		today := time.Now()
// 		for rows.Next() {
// 			var event parser.Event
// 			if err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.EventDate); err != nil {
// 				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse events"})
// 			}

// 			eventDate, _ := time.Parse("2006-01-02", event.EventDate.String())
// 			daysLeft := int(eventDate.Sub(today).Hours() / 24)

// 			events = append(events, map[string]interface{}{
// 				"id":          event.ID,
// 				"name":        event.Name,
// 				"description": event.Description,
// 				"event_date":  event.EventDate,
// 				"days_left":   daysLeft,
// 			})
// 		}

// 		return c.Render(http.StatusOK, "calendar/calendar_create.gohtml", map[string]interface{}{
// 			"Events": events,
// 		})
// 	}
// }

func AddEvent(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Require user to be logged in
		if err := util.RequireLogin(c); err != nil {
			return err
		}

		user, _ := util.GetUserFromContext(c)

		name := c.FormValue("name")
		description := c.FormValue("description")
		eventDate := c.FormValue("event_date")

		if name == "" || eventDate == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name and Event Date are required"})
		}

		_, err := db.Exec(
			"INSERT INTO events (name, description, event_date, username, created_at) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)",
			name, description, eventDate, user.Username,
		)
		if err != nil {
			c.Logger().Errorf("Failed to add event: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add event"})
		}

		return c.Redirect(http.StatusSeeOther, "/calendar")
	}
}

func GetAll(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Require user to be logged in
		if err := util.RequireLogin(c); err != nil {
			return err
		}

		user, _ := util.GetUserFromContext(c)

		rows, err := db.Query(
			"SELECT id, name, description, event_date, created_at FROM events WHERE username = ? ORDER BY event_date ASC",
			user.Username,
		)
		if err != nil {
			c.Logger().Errorf("Failed to fetch events: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch events"})
		}
		defer rows.Close()

		var events []parser.Event
		for rows.Next() {
			var event parser.Event
			if err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.EventDate, &event.CreatedAt); err != nil {
				c.Logger().Errorf("Failed to parse events: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse events"})
			}

			// Calculate days left
			event.DaysLeft = int(time.Until(event.EventDate).Hours() / 24)
			events = append(events, event)
		}

		// Render events
		return c.Render(http.StatusOK, "calendar/calendar.gohtml", map[string]interface{}{
			"Events":     events,
			"darkMode":   user.DarkMode,
			"isLoggedIn": user.IsLoggedIn,
		})
	}
}
