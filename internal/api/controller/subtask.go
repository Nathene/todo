package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"todo/internal/db"
	"todo/internal/util"

	"github.com/labstack/echo/v4"
)

type subtaskReq struct {
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
}

func CreateSubtask(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		username, ok := c.Get("username").(string)
		if !ok || username == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		todoIDParam := c.Param("todo_id")
		todoID, err := strconv.Atoi(todoIDParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid todo_id"})
		}

		// Verify the todo list exists
		var todoExists bool
		err = db.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM todo_lists WHERE id = ? AND username = ?)", todoID, username,
		).Scan(&todoExists)
		if err != nil || !todoExists {
			return c.JSON(http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Todo list '%d' not found for user '%s'", todoID, username)})
		}

		// Parse the subtask details
		var req subtaskReq
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		if req.Name == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Subtask name is required"})
		}

		// Insert the subtask
		_, err = db.Exec(
			"INSERT INTO subtasks (todo_id, name, description, done, username, created_at, updated_at) VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
			todoID, req.Name, req.Description, false, username,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create subtask"})
		}

		// Debug log for confirmation
		fmt.Printf("Subtask created: todoID=%d, name=%s, username=%s\n", todoID, req.Name, username)

		// Redirect back to the referring page or re-render subtasks
		return util.BackPage(c)
	}
}

func GetSubtasks(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		username, ok := c.Get("username").(string)
		if !ok || username == "" {
			return c.Redirect(http.StatusUnauthorized, "/login")
		}

		todoIDParam := c.Param("todo_id")
		todoID, err := strconv.Atoi(todoIDParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid todo_id"})
		}

		rows, err := db.Query("SELECT id, name, description, done FROM subtasks WHERE todo_id = ?", todoID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve subtasks"})
		}
		defer rows.Close()

		var subtasks []map[string]interface{}
		for rows.Next() {
			var id int
			var name, description string
			var done bool
			if err := rows.Scan(&id, &name, &description, &done); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse subtasks"})
			}
			subtasks = append(subtasks, map[string]interface{}{
				"id":          id,
				"name":        name,
				"description": description,
				"done":        done,
			})
		}

		return c.Render(http.StatusOK, "subtasks/subtasks.gohtml", map[string]interface{}{
			"Todo": map[string]interface{}{
				"ID":          todoID,
				"Name":        "Sample Todo",
				"Description": "Sample Description",
			},
			"Subtasks": subtasks,
		})
	}
}
