package subtask

import (
	"net/http"
	"strconv"
	"todo/internal/db"

	"github.com/labstack/echo/v4"
)

// func ToggleSubtaskDone(db *db.Database) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		// Get the subtask ID from the URL parameter
// 		subtaskID, err := strconv.Atoi(c.Param("id"))
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusBadRequest, "Invalid subtask ID")
// 		}

// 		// Get the "done" value from the form
// 		done := c.FormValue("done") == "on"

// 		// Update the database
// 		err = db.UpdateSubtaskDone(subtaskID, done)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update subtask status")
// 		}

// 		// Redirect back to the current page
// 		return util.BackPage(c)
// 	}
// }

func ToggleSubtaskDone(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		subtaskID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid subtask ID")
		}

		// Toggle the done status of the subtask
		_, err = db.Exec("UPDATE subtasks SET done = NOT done WHERE id = ?", subtaskID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to toggle subtask status")
		}

		// Check if all subtasks are done for the parent todo
		var todoID int
		err = db.QueryRow("SELECT todo_id FROM subtasks WHERE id = ?", subtaskID).Scan(&todoID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch parent Todo ID")
		}

		var allDone bool
		err = db.QueryRow("SELECT NOT EXISTS (SELECT 1 FROM subtasks WHERE todo_id = ? AND done = false)", todoID).Scan(&allDone)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check subtask completion status")
		}

		// Update the parent todo status if all subtasks are done
		if allDone {
			_, err = db.Exec("UPDATE todo_lists SET status = 'Completed' WHERE id = ?", todoID)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update todo status")
			}
		} else {
			// Reset the status if subtasks are not fully done
			_, err = db.Exec("UPDATE todo_lists SET status = 'In Progress' WHERE id = ?", todoID)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to reset todo status")
			}
		}

		// Redirect back to the edit page or referer
		referer := c.Request().Header.Get("Referer")
		if referer == "" {
			referer = "/groups" // Fallback if no referer
		}
		return c.Redirect(http.StatusSeeOther, referer)
	}
}
