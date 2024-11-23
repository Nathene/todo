package controller

import (
	"log"
	"net/http"
	"todo/internal/db"
	"todo/internal/parser"

	"github.com/labstack/echo/v4"
)

type RequestBody struct {
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
	Urgent      bool   `json:"urgent" form:"urgent"`
	Priority    int    `json:"priority" form:"priority"`
	Status      string `json:"status" form:"status"`
}

func CreateTodoList(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println("Entered CreateTodoList Handler")

		user, ok := c.Get("user").(parser.User)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		groupID := c.Param("id")
		if groupID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Group ID is required"})
		}

		log.Println("Form Values: ", c.Request().PostForm)

		name := c.FormValue("name")
		if name == "" {
			log.Println("Missing 'name' field in form data")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "List name is required"})
		}

		description := c.FormValue("description")
		urgent := c.FormValue("urgent") == "on"
		priority := c.FormValue("priority")
		status := c.FormValue("status")

		log.Printf("Received Data: name=%s, description=%s, urgent=%t, priority=%s, status=%s\n",
			name, description, urgent, priority, status)

		_, err := db.Exec(
			"INSERT INTO todo_lists (group_id, name, description, urgent, priority, done, status, username, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
			groupID, name, description, urgent, priority, false, status, user.Username,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create todo list"})
		}

		// Redirect to the group page
		return c.Redirect(http.StatusSeeOther, "/groups/"+groupID)
	}
}

// func GetTodoListsByGroupName(db *db.Database) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		// Extract username from context
// 		user, ok := c.Get("user").(parser.User)
// 		if !ok {
// 			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: user not found in context"})
// 		}

// 		// Extract group name from URL parameters
// 		groupName := c.Param("name")
// 		if groupName == "" {
// 			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Group name is required"})
// 		}

// 		// Fetch group ID
// 		groupID, err := db.GetIDfromGroup(groupName, user.Username)
// 		if err != nil {
// 			return c.JSON(http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Group '%s' not found for user '%s'", groupName, user.Username)})
// 		}

// 		// Fetch todo lists for the group
// 		rows, err := db.Query(
// 			"SELECT id, name, username, description, urgent, priority, status, done, created_at, updated_at FROM todo_lists WHERE group_id = ? AND username = ?",
// 			groupID, user.Username,
// 		)
// 		if err != nil {
// 			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch todo lists"})
// 		}
// 		defer rows.Close()

// 		var todoLists []map[string]interface{}
// 		for rows.Next() {
// 			var todoList struct {
// 				ID          int    `json:"id"`
// 				Name        string `json:"name"`
// 				Username    string `json:"username"`
// 				Description string `json:"description"`
// 				Urgent      bool   `json:"urgent"`
// 				Priority    int    `json:"priority"`
// 				Status      string `json:"status"`
// 				Done        bool   `json:"done"`
// 				CreatedAt   string `json:"created_at"`
// 				UpdatedAt   string `json:"updated_at"`
// 			}
// 			if err := rows.Scan(
// 				&todoList.ID, &todoList.Name, &todoList.Username, &todoList.Description,
// 				&todoList.Urgent, &todoList.Priority, &todoList.Status, &todoList.Done,
// 				&todoList.CreatedAt, &todoList.UpdatedAt,
// 			); err != nil {
// 				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse todo lists"})
// 			}

// 			// Fetch subtasks for the current todo list
// 			subtaskRows, err := db.Query(
// 				"SELECT id, todo_id, name, description, done, created_at, updated_at FROM subtasks WHERE todo_id = ?", todoList.ID,
// 			)
// 			if err != nil {
// 				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch subtasks"})
// 			}
// 			defer subtaskRows.Close()

// 			var subtasks []map[string]any
// 			for subtaskRows.Next() {
// 				var subtask parser.Subtask
// 				if err := subtaskRows.Scan(&subtask.ID, &subtask.TodoID, &subtask.Name, &subtask.Description, &subtask.Done, &subtask.CreatedAt, &subtask.UpdatedAt); err != nil {
// 					return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse subtasks"})
// 				}
// 				subtasks = append(subtasks, map[string]any{
// 					"id":          subtask.ID,
// 					"todo_id":     subtask.TodoID,
// 					"name":        subtask.Name,
// 					"done":        subtask.Done,
// 					"description": subtask.Description,
// 					"created_at":  subtask.CreatedAt,
// 					"updated_at":  subtask.UpdatedAt,
// 				})
// 			}

// 			// Append todo list with subtasks
// 			todoLists = append(todoLists, map[string]any{
// 				"id":          todoList.ID,
// 				"name":        todoList.Name,
// 				"username":    todoList.Username,
// 				"description": todoList.Description,
// 				"urgent":      todoList.Urgent,
// 				"priority":    todoList.Priority,
// 				"status":      todoList.Status,
// 				"done":        todoList.Done,
// 				"created_at":  todoList.CreatedAt,
// 				"updated_at":  todoList.UpdatedAt,
// 				"subtasks":    subtasks,
// 			})
// 		}

// 		return c.JSON(http.StatusOK, todoLists)
// 	}
// }

func GetTodoListsByGroupID(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract user information from context
		user, ok := c.Get("user").(parser.User)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: user not found in context"})
		}

		// Extract group ID from URL parameters
		groupID := c.Param("id")
		if groupID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Group ID is required"})
		}

		// Fetch the group name
		var groupName string
		err := db.DB.QueryRow("SELECT name FROM todo_groups WHERE id = ? AND username = ?", groupID, user.Username).Scan(&groupName)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Group not found"})
		}

		// Fetch todo lists for the group ID
		rows, err := db.Query(
			"SELECT id, name, username, description, urgent, priority, status, done, created_at, updated_at FROM todo_lists WHERE group_id = ? AND username = ?",
			groupID, user.Username,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch todo lists"})
		}
		defer rows.Close()

		var todoLists []parser.TodoList
		for rows.Next() {
			var todoList parser.TodoList
			if err := rows.Scan(
				&todoList.ID, &todoList.Name, &todoList.Username, &todoList.Description,
				&todoList.Urgent, &todoList.Priority, &todoList.Status, &todoList.Done,
				&todoList.Created_at, &todoList.Updated_at,
			); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse todo lists"})
			}

			// Fetch subtasks for each todo list
			subtaskRows, err := db.Query(
				"SELECT id, todo_id, name, description, done, created_at, updated_at FROM subtasks WHERE todo_id = ?",
				todoList.ID,
			)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch subtasks"})
			}
			defer subtaskRows.Close()

			for subtaskRows.Next() {
				var subtask parser.Subtask
				if err := subtaskRows.Scan(
					&subtask.ID, &subtask.TodoID, &subtask.Name, &subtask.Description,
					&subtask.Done, &subtask.CreatedAt, &subtask.UpdatedAt,
				); err != nil {
					return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse subtasks"})
				}
				todoList.Subtasks = append(todoList.Subtasks, subtask)
			}

			// Append todo list to the slice
			todoLists = append(todoLists, todoList)
		}

		// Render group details page
		return c.Render(http.StatusOK, "todos/todos.gohtml", map[string]interface{}{
			"GroupID":    groupID,
			"GroupName":  groupName,
			"Username":   user.Username,
			"TodoLists":  todoLists,
			"darkMode":   user.DarkMode,
			"isLoggedIn": user.IsLoggedIn,
		})
	}
}

func UpdateTodoStatus(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(parser.User)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		todoID := c.Param("id")
		if todoID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Todo ID is required"})
		}

		status := c.FormValue("status")
		if status == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Status is required"})
		}

		_, err := db.Exec(
			"UPDATE todo_lists SET status = ?, done = ? WHERE id = ? AND username = ?",
			status, status == "Completed", todoID, user.Username,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update status"})
		}

		referer := c.Request().Header.Get("Referer")
		if referer == "" {
			referer = "/" // Default to homepage
		}
		return c.Redirect(http.StatusSeeOther, referer)
	}
}

func UpdateTodoPriority(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(parser.User)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		todoID := c.Param("id")
		if todoID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Todo ID is required"})
		}

		priority := c.FormValue("priority")
		if priority == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Priority is required"})
		}

		_, err := db.Exec(
			"UPDATE todo_lists SET priority = ? WHERE id = ? AND username = ?",
			priority, todoID, user.Username,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update priority"})
		}

		referer := c.Request().Header.Get("Referer")
		if referer == "" {
			referer = "/" // Default to homepage
		}
		return c.Redirect(http.StatusSeeOther, referer)
	}
}

func DeleteTodo(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(parser.User)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		todoID := c.Param("id")
		if todoID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Todo ID is required"})
		}

		_, err := db.Exec(
			"DELETE FROM todo_lists WHERE id = ? AND username = ?",
			todoID, user.Username,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete todo"})
		}

		referer := c.Request().Header.Get("Referer")
		if referer == "" {
			referer = "/" // Default to homepage
		}
		return c.Redirect(http.StatusSeeOther, referer)
	}
}
