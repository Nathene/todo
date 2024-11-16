package debug

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func DebugGroupsTable(db *sql.DB) {
	fmt.Println("\n=== Debug: todo_groups Table ===")
	rows, err := db.Query("SELECT id, name, username FROM todo_groups")
	if err != nil {
		log.Printf("Failed to query todo_groups: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Printf("%-5s %-20s %-20s\n", "ID", "Group Name", "Username")
	fmt.Println(strings.Repeat("-", 50))
	for rows.Next() {
		var id int
		var name, username string
		if err := rows.Scan(&id, &name, &username); err != nil {
			log.Printf("Failed to scan todo_groups row: %v\n", err)
			return
		}
		fmt.Printf("%-5d %-20s %-20s\n", id, name, username)
	}
	fmt.Println(strings.Repeat("=", 50))
}

func DebugListsTable(db *sql.DB) {
	fmt.Println("\n=== Debug: todo_lists Table ===")
	rows, err := db.Query("SELECT id, group_id, name, description, urgent, priority, done FROM todo_lists")
	if err != nil {
		log.Printf("Failed to query todo_lists: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Printf("%-5s %-10s %-20s %-30s %-7s %-9s %-5s\n", "ID", "Group ID", "Name", "Description", "Urgent", "Priority", "Done")
	fmt.Println(strings.Repeat("-", 100))
	for rows.Next() {
		var id, groupID, priority int
		var name, description string
		var urgent, done bool
		if err := rows.Scan(&id, &groupID, &name, &description, &urgent, &priority, &done); err != nil {
			log.Printf("Failed to scan todo_lists row: %v\n", err)
			return
		}
		fmt.Printf("%-5d %-10d %-20s %-30s %-7t %-9d %-5t\n", id, groupID, name, description, urgent, priority, done)
	}
	fmt.Println(strings.Repeat("=", 100))
}

func DebugSubtasksTable(db *sql.DB) {
	fmt.Println("\n=== Debug: subtasks Table ===")
	rows, err := db.Query("SELECT id, todo_id, name, done FROM subtasks")
	if err != nil {
		log.Printf("Failed to query subtasks: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Printf("%-5s %-10s %-20s %-5s\n", "ID", "Todo ID", "Name", "Done")
	fmt.Println(strings.Repeat("-", 50))
	for rows.Next() {
		var id, todoID int
		var name string
		var done bool
		if err := rows.Scan(&id, &todoID, &name, &done); err != nil {
			log.Printf("Failed to scan subtasks row: %v\n", err)
			return
		}
		fmt.Printf("%-5d %-10d %-20s %-5t\n", id, todoID, name, done)
	}
	fmt.Println(strings.Repeat("=", 50))
}

func DebugGroupTodoLists(db *sql.DB, groupName, username string) {
	fmt.Printf("\n=== Debug: Todo Lists for Group '%s' and User '%s' ===\n", groupName, username)
	var groupID int
	err := db.QueryRow("SELECT id FROM todo_groups WHERE name = ? AND username = ?", groupName, username).Scan(&groupID)
	if err != nil {
		log.Printf("Group not found: %s for user %s: %v\n", groupName, username, err)
		return
	}
	fmt.Printf("Group ID: %d\n", groupID)

	rows, err := db.Query("SELECT id, name, description, urgent, priority, done FROM todo_lists WHERE group_id = ?", groupID)
	if err != nil {
		log.Printf("Failed to query todo lists for group '%s': %v\n", groupName, err)
		return
	}
	defer rows.Close()

	fmt.Printf("%-5s %-20s %-30s %-7s %-9s %-5s\n", "ID", "Name", "Description", "Urgent", "Priority", "Done")
	fmt.Println(strings.Repeat("-", 100))
	for rows.Next() {
		var id, priority int
		var name, description string
		var urgent, done bool
		if err := rows.Scan(&id, &name, &description, &urgent, &priority, &done); err != nil {
			log.Printf("Failed to scan todo_lists row: %v\n", err)
			return
		}
		fmt.Printf("%-5d %-20s %-30s %-7t %-9d %-5t\n", id, name, description, urgent, priority, done)
	}
	fmt.Println(strings.Repeat("=", 100))
}
