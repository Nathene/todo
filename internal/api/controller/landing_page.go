package controller

import (
	"net/http"
	"time"
	"todo/internal/db"
	"todo/internal/parser"

	"github.com/labstack/echo/v4"
)

func Dashboard(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get logged-in user details from context
		user, ok := c.Get("user").(parser.User)
		if !ok || !user.IsLoggedIn {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		// Fetch upcoming events within a week
		rows, err := db.Query(`
			SELECT id, name, description, event_date 
			FROM events 
			WHERE username = ? AND event_date BETWEEN CURRENT_DATE AND DATE(CURRENT_DATE, '+7 days')
			ORDER BY event_date ASC
		`, user.Username)
		if err != nil {
			c.Logger().Errorf("Failed to fetch upcoming events: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load events"})
		}
		defer rows.Close()

		var upcomingEvents []parser.Event
		for rows.Next() {
			var event parser.Event
			if err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.EventDate); err != nil {
				c.Logger().Errorf("Failed to scan event: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process event data"})
			}
			event.DaysLeft = int(time.Until(event.EventDate).Hours() / 24)
			upcomingEvents = append(upcomingEvents, event)
		}

		// Fetch groups with urgent tickets
		rows, err = db.Query(`
			SELECT g.id, g.name, l.id as list_id, l.name as list_name, l.description
			FROM todo_groups g
			JOIN todo_lists l ON g.id = l.group_id
			WHERE l.username = ? AND l.urgent = true AND l.done = false
			ORDER BY g.name, l.priority DESC
		`, user.Username)
		if err != nil {
			c.Logger().Errorf("Failed to fetch groups with urgent tickets: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load groups with urgent tickets"})
		}
		defer rows.Close()

		var groupsWithUrgentTickets []map[string]interface{}
		for rows.Next() {
			var groupID int
			var groupName string
			var listID int
			var listName string
			var description string

			if err := rows.Scan(&groupID, &groupName, &listID, &listName, &description); err != nil {
				c.Logger().Errorf("Failed to scan group and urgent ticket: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process group and urgent ticket data"})
			}

			groupsWithUrgentTickets = append(groupsWithUrgentTickets, map[string]interface{}{
				"groupID":     groupID,
				"groupName":   groupName,
				"listID":      listID,
				"listName":    listName,
				"description": description,
			})
		}

		// Render the landing page template
		return c.Render(http.StatusOK, "landing/landing.gohtml", map[string]interface{}{
			"darkMode":                user.DarkMode,
			"isLoggedIn":              user.IsLoggedIn,
			"username":                user.Username,
			"UpcomingEvents":          upcomingEvents,
			"GroupsWithUrgentTickets": groupsWithUrgentTickets,
		})
	}
}
