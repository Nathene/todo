package db

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

type jsonErr struct {
	StatusCode int               `json:"code"`
	Errors     map[string]string `json:"errors"`
}

func (j *jsonErr) Error() string {
	for _, error := range j.Errors {
		return error
	}
	return ""
}

// GetIDfromGroup fetches the ID of a todo group based on group name and username
func (db *Database) GetIDfromGroup(groupName, userName string) (int, error) {
	var id string

	// Query the database for the group ID by name and username
	err := db.QueryRow("SELECT id FROM todo_groups WHERE name = ? AND username = ?", groupName, userName).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found for the given group name and username
			return -1, &jsonErr{http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("no rows for given group name '%s' or username '%s'", groupName, userName)}}
		}
		// Other errors
		return -1, &jsonErr{http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("something went wrong for group '%s'", groupName)}}
	}
	groupID, err := strconv.Atoi(id)
	if err != nil {
		return -1, err
	}
	return groupID, nil
}

// GroupExists checks if a group exists for a given user
func (db *Database) GroupExists(groupName, userName string) error {
	var groupID int
	err := db.QueryRow(
		"SELECT id FROM todo_groups WHERE name = ? AND username = ?", groupName, userName,
	).Scan(&groupID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &jsonErr{http.StatusNotFound, map[string]string{"error": fmt.Sprintf("group '%s' not found for user '%s'", groupName, userName)}}
		}
		return &jsonErr{http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("something went wrong while checking group '%s'", groupName)}}
	}
	return nil
}

// AddListToGroup adds a todo list to a group, ensuring ownership with username
func (db *Database) AddListToGroup(groupID int, username string, req struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Urgent      bool   `json:"urgent"`
}) error {
	_, err := db.Exec(
		"INSERT INTO todo_lists (group_id, name, description, urgent, status, done, username, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
		groupID, req.Name, req.Description, req.Urgent, "not_started", false, username,
	)
	if err != nil {
		return &jsonErr{http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to add list '%s' to group '%d'", req.Name, groupID)}}
	}
	return nil
}
