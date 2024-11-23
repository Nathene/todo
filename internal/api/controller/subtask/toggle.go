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
		// Get the subtask ID and new state
		subtaskID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid subtask ID")
		}
		done := c.FormValue("done") == "on"

		// Update the subtask
		err = db.UpdateSubtaskDone(subtaskID, done)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update subtask status")
		}

		// Redirect back to the subtasks section
		referer := c.Request().Header.Get("Referer")
		if referer == "" {
			referer = "/groups" // Fallback if referer is missing
		}
		return c.Redirect(http.StatusSeeOther, referer+"#subtasks")
	}
}
