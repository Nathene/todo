package db

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
	"todo/internal/parser"
)

func (db *Database) InsertList(groupID int, username string, req struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Urgent      bool   `json:"urgent"`
	Priority    int    `json:"priority"`
	Status      string `json:"status"`
}) error {
	// Insert the todo list into the database
	log.Printf("Inserting todo list with group ID: %d for username: %s", groupID, username)

	_, err := db.Exec(
		"INSERT INTO todo_lists (group_id, name, description, urgent, priority, done, status, username, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		groupID, req.Name, req.Description, req.Urgent, req.Priority, false, req.Status, username, time.Now(), time.Now(),
	)
	if err != nil {
		return errors.New("Failed to create todo list")
	}
	return nil
}

func (db *Database) GetTodoLists(groupID int, username string) (*[]parser.TodoList, error) {
	log.Printf("Fetching todo lists for groupID: %d, username: %s", groupID, username)

	// Execute the query
	rows, err := db.Query(
		"SELECT id, name, description, CAST(urgent AS BOOLEAN), priority, status, done, created_at, updated_at FROM todo_lists WHERE group_id = ? AND username = ?",
		groupID, username,
	)
	if err != nil {
		log.Printf("Query error for groupID '%d' and username '%s': %v", groupID, username, err)
		return nil, &jsonErr{
			StatusCode: http.StatusNotFound,
			Errors:     map[string]string{"error": fmt.Sprintf("No todo lists found for groupID '%d'", groupID)},
		}
	}
	defer rows.Close()

	var todoLists []parser.TodoList
	for rows.Next() {
		var todoList parser.TodoList
		if err := rows.Scan(&todoList.ID, &todoList.Name, &todoList.Description, &todoList.Urgent, &todoList.Priority, &todoList.Status, &todoList.Done, &todoList.CreatedAt, &todoList.UpdatedAt); err != nil {
			log.Printf("Row scan error for groupID '%d': %v", groupID, err)
			return nil, &jsonErr{
				StatusCode: http.StatusInternalServerError,
				Errors:     map[string]string{"error": "Failed to parse todo list data"},
			}
		}
		todoLists = append(todoLists, todoList)
	}

	// If no rows are found
	if len(todoLists) == 0 {
		log.Printf("No todo lists found for groupID '%d' and username '%s'", groupID, username)
		return &todoLists, nil
	}

	log.Printf("Successfully fetched %d todo lists for groupID '%d'", len(todoLists), groupID)
	return &todoLists, nil
}

func (db *Database) ListExists(listName, username string) (int, error) {
	var listID int
	err := db.QueryRow(
		"SELECT id FROM todo_lists WHERE name = ? AND username = ?", listName, username,
	).Scan(&listID)
	if err != nil {
		return -1, &jsonErr{http.StatusNotFound, map[string]string{"error": fmt.Sprintf("list '%s' not found", listName)}}
	}
	return listID, nil
}
