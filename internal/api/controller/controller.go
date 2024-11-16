package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"todo/internal/db"
	"todo/internal/service"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var jwtSecret = []byte("your-secret-key") // Replace with a secure secret key

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JWTClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Figure out how i can use this type, instead of having a million functions everywhere.
type User struct {
	userName string
	lists    []service.TodoList
	groups   []service.TodoGroup
}

func Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req LoginRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		// Replace this with database validation
		if req.Username != "admin" || req.Password != "password" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
		}

		// Generate JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
			Username: req.Username,
			StandardClaims: jwt.StandardClaims{
				// FIXME change this, adding infinite expiry for testing
				ExpiresAt: time.Now().Add(time.Hour * 100000).Unix(), // Token expires in 1 hour
			},
		})
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		}

		return c.JSON(http.StatusOK, map[string]string{"token": tokenString})
	}
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing or invalid token"})
		}

		token, err := jwt.ParseWithClaims(authHeader, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid claims"})
		}

		c.Set("username", claims.Username)
		return next(c)
	}
}

func CreateSubtask(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract username from context
		username, ok := c.Get("username").(string)
		if !ok || username == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: username not found in context"})
		}

		// Extract parameters from URL
		groupName := c.Param("group_name") // Group name from URL
		todoIDParam := c.Param("todo_id")  // Todo ID from URL

		// Convert `todo_id` to integer
		todoID, err := strconv.Atoi(todoIDParam)
		if err != nil {
			c.Logger().Errorf("Invalid todo_id: %s", todoIDParam)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid todo_id"})
		}

		// Verify the group exists for the user
		groupID, err := db.GetIDfromGroup(groupName, username)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Group '%s' not found for user '%s'", groupName, username)})
		}

		// Verify the todo list exists within the group for the user
		var todoExists bool
		err = db.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM todo_lists WHERE id = ? AND group_id = ? AND username = ?)", todoID, groupID, username,
		).Scan(&todoExists)
		if err != nil || !todoExists {
			return c.JSON(http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Todo list '%d' not found in group '%s' for user '%s'", todoID, groupName, username)})
		}

		// Parse the subtask details
		type RequestBody struct {
			Name string `json:"name"`
		}
		var req RequestBody
		if err := c.Bind(&req); err != nil {
			c.Logger().Errorf("Failed to parse request body: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		if req.Name == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Subtask name is required"})
		}

		// Insert the subtask into the database
		result, err := db.Exec(
			"INSERT INTO subtasks (todo_id, name, done, username, created_at, updated_at) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
			todoID, req.Name, false, username,
		)
		if err != nil {
			c.Logger().Errorf("Failed to create subtask for todo_id '%d': %v", todoID, err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create subtask"})
		}

		// Fetch the ID of the newly created subtask
		subtaskID, err := result.LastInsertId()
		if err != nil {
			c.Logger().Errorf("Failed to fetch subtask ID for todo_id '%d': %v", todoID, err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve subtask ID"})
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message":    "Subtask created successfully",
			"subtask_id": subtaskID,
		})
	}
}

func GetSubtasks(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract username from context
		userName, ok := c.Get("username").(string)
		if !ok || userName == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: username not found in context"})
		}

		// Extract parameters from URL
		groupName := c.Param("group_name") // Group name from URL
		todoIDParam := c.Param("todo_id")  // Todo ID from URL

		// Convert `todo_id` to integer
		todoID, err := strconv.Atoi(todoIDParam)
		if err != nil {
			c.Logger().Errorf("Invalid todo_id: %s", todoIDParam)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid todo_id"})
		}

		// Verify the group exists for the user
		groupID, err := db.GetIDfromGroup(groupName, userName)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Group '%s' not found for user '%s'", groupName, userName)})
		}

		// Verify the todo list exists within the group for the user
		var todoExists bool
		err = db.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM todo_lists WHERE id = ? AND group_id = ? AND username = ?)", todoID, groupID, userName,
		).Scan(&todoExists)
		if err != nil || !todoExists {
			return c.JSON(http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Todo list '%d' not found in group '%s' for user '%s'", todoID, groupName, userName)})
		}

		// Fetch the subtasks for the todo list
		rows, err := db.Query("SELECT id, todo_id, name, done FROM subtasks WHERE todo_id = ?", todoID)
		if err != nil {
			c.Logger().Errorf("Failed to fetch subtasks for todo_id '%d': %v", todoID, err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve subtasks"})
		}
		defer rows.Close()

		var subtasks []service.Subtask
		for rows.Next() {
			subtask := service.Subtask{}
			if err := rows.Scan(&subtask.ID, &subtask.TodoID, &subtask.Name, &subtask.Done); err != nil {
				c.Logger().Errorf("Failed to parse subtasks for todo_id '%d': %v", todoID, err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse subtasks"})
			}
			subtasks = append(subtasks, subtask)
		}

		return c.JSON(http.StatusOK, subtasks)
	}
}

func CreateTodoGroup(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		username, ok := c.Get("username").(string)
		if !ok || username == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: username not found in context"})
		}

		type RequestBody struct {
			Name string `json:"name"`
		}

		var req RequestBody
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		if req.Name == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Group name is required"})
		}

		// Check if the group already exists for the user
		err := db.GroupExists(req.Name, username)
		if err == nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": fmt.Sprintf("Group '%s' already exists for user '%s'", req.Name, username)})
		}

		// Insert the new group into the database
		res, err := db.DB.Exec("INSERT INTO todo_groups (name, username, created_at, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)", req.Name, username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create group"})
		}

		groupID, _ := res.LastInsertId()
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"id":       groupID,
			"name":     req.Name,
			"username": username,
		})
	}
}

func GetTodoGroups(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		username, ok := c.Get("username").(string)
		if !ok || username == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: username not found in context"})
		}

		// Fetch all groups for the user
		rows, err := db.Query(`SELECT id, name FROM todo_groups WHERE username = ?`, username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve groups"})
		}
		defer rows.Close()

		var groups []map[string]interface{}
		for rows.Next() {
			var id int
			var name string
			if err := rows.Scan(&id, &name); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse groups"})
			}
			groups = append(groups, map[string]interface{}{
				"id":         id,
				"group_name": name,
			})
		}

		return c.JSON(http.StatusOK, groups)
	}
}

func CreateTodoList(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract username from context
		username, ok := c.Get("username").(string)
		if !ok || username == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: username not found in context"})
		}

		// Parse the request body
		type RequestBody struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Urgent      bool   `json:"urgent"`
			Priority    int    `json:"priority"`
			Status      string `json:"status"`
		}
		var req RequestBody
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		// Ensure required fields are provided
		if req.Name == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Todo list name is required"})
		}

		// Insert the todo list into the database
		_, err := db.Exec(
			"INSERT INTO todo_lists (group_id, name, description, urgent, priority, done, status, username, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
			1, req.Name, req.Description, req.Urgent, req.Priority, false, req.Status, username,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create todo list"})
		}

		return c.JSON(http.StatusCreated, map[string]string{"message": "Todo list created successfully"})
	}
}

func GetTodoListsByGroupName(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract username from context
		username, ok := c.Get("username").(string)
		if !ok || username == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: username not found in context"})
		}

		// Extract group name from URL parameters
		groupName := c.Param("name")
		if groupName == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Group name is required"})
		}

		// Fetch group ID
		groupID, err := db.GetIDfromGroup(groupName, username)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Group '%s' not found for user '%s'", groupName, username)})
		}

		// Fetch todo lists for the group
		rows, err := db.Query(
			"SELECT id, name, username, description, urgent, priority, status, done, created_at, updated_at FROM todo_lists WHERE group_id = ? AND username = ?",
			groupID, username,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch todo lists"})
		}
		defer rows.Close()

		var todoLists []map[string]interface{}
		for rows.Next() {
			var todoList struct {
				ID          int    `json:"id"`
				Name        string `json:"name"`
				Username    string `json:"username"`
				Description string `json:"description"`
				Urgent      bool   `json:"urgent"`
				Priority    int    `json:"priority"`
				Status      string `json:"status"`
				Done        bool   `json:"done"`
				CreatedAt   string `json:"created_at"`
				UpdatedAt   string `json:"updated_at"`
			}
			if err := rows.Scan(
				&todoList.ID, &todoList.Name, &todoList.Username, &todoList.Description,
				&todoList.Urgent, &todoList.Priority, &todoList.Status, &todoList.Done,
				&todoList.CreatedAt, &todoList.UpdatedAt,
			); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse todo lists"})
			}

			// Fetch subtasks for the current todo list
			subtaskRows, err := db.Query(
				"SELECT id, todo_id, name, done, created_at, updated_at FROM subtasks WHERE todo_id = ?", todoList.ID,
			)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch subtasks"})
			}
			defer subtaskRows.Close()

			var subtasks []map[string]interface{}
			for subtaskRows.Next() {
				var subtask struct {
					ID        int    `json:"id"`
					TodoID    int    `json:"todo_id"`
					Name      string `json:"name"`
					Done      bool   `json:"done"`
					CreatedAt string `json:"created_at"`
					UpdatedAt string `json:"updated_at"`
				}
				if err := subtaskRows.Scan(&subtask.ID, &subtask.TodoID, &subtask.Name, &subtask.Done, &subtask.CreatedAt, &subtask.UpdatedAt); err != nil {
					return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse subtasks"})
				}
				subtasks = append(subtasks, map[string]interface{}{
					"id":         subtask.ID,
					"todo_id":    subtask.TodoID,
					"name":       subtask.Name,
					"done":       subtask.Done,
					"created_at": subtask.CreatedAt,
					"updated_at": subtask.UpdatedAt,
				})
			}

			// Append todo list with subtasks
			todoLists = append(todoLists, map[string]interface{}{
				"id":          todoList.ID,
				"name":        todoList.Name,
				"username":    todoList.Username,
				"description": todoList.Description,
				"urgent":      todoList.Urgent,
				"priority":    todoList.Priority,
				"status":      todoList.Status,
				"done":        todoList.Done,
				"created_at":  todoList.CreatedAt,
				"updated_at":  todoList.UpdatedAt,
				"subtasks":    subtasks,
			})
		}

		return c.JSON(http.StatusOK, todoLists)
	}
}

func AddTodoListToGroup(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		username, ok := c.Get("username").(string)
		if !ok || username == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: username not found in context"})
		}

		groupName := c.Param("name")
		if groupName == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Group name is required"})
		}

		type RequestBody struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Urgent      bool   `json:"urgent"`
		}

		var req RequestBody
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		if req.Name == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "List name is required"})
		}

		groupID, err := db.GetIDfromGroup(groupName, username)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Group '%s' not found for user '%s'", groupName, username)})
		}

		err = db.AddListToGroup(groupID, username, req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add list to group"})
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"status":  http.StatusCreated,
			"message": "Todo list added successfully",
		})
	}
}
