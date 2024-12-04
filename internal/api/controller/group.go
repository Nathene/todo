package controller

import (
	"fmt"
	"log"
	"net/http"
	"todo/internal/db"
	"todo/internal/parser"
	"todo/internal/util"

	"github.com/labstack/echo/v4"
)

func CreateTodoGroup(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(parser.User)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: user not found"})
		}

		groupName := c.FormValue("name")
		if groupName == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Group name is required"})
		}

		// Check if the group already exists for the user
		err := db.GroupExists(groupName, user.Username)
		if err == nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Group already exists"})
		}

		// Insert the group into the database
		_, err = db.Exec("INSERT INTO todo_groups (name, username, created_at, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)", groupName, user.Username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create group"})
		}

		return c.Redirect(http.StatusSeeOther, "/groups") // Redirect back to the groups page
	}
}

func GetTodoGroups(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := util.RequireLogin(c); err != nil {
			return err
		}

		user, ok := util.GetUserFromContext(c)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: user not found"})
		}

		rows, err := db.Query(`SELECT id, name FROM todo_groups WHERE username = ?`, user.Username)
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
				"id":   id,
				"name": name,
			})
		}

		return c.Render(http.StatusOK, "groups/groups.gohtml", map[string]interface{}{
			"Groups":     groups,
			"user":       user,
			"isLoggedIn": user.IsLoggedIn,
			"darkMode":   user.DarkMode,
		})
	}
}

func AddTodoListToGroup(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println("WHY AM I HERE?")
		user, ok := c.Get("user").(parser.User)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: user not found in context"})
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

		groupID, err := db.GetIDfromGroup(groupName, user.Username)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Group '%s' not found for user '%s'", groupName, user.Username)})
		}

		err = db.AddListToGroup(groupID, user.Username, req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add list to group"})
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"status":  http.StatusCreated,
			"message": "Todo list added successfully",
			"user":    user,
		})
	}
}
