package db

import (
	"database/sql"
	"errors"
	"todo/internal/parser"
)

func (db *Database) GetUser(username string) (*parser.User, error) {
	var user parser.User

	// Fetch basic user details
	query := `
		SELECT username, firstname, email, darkmode
		FROM users
		WHERE username = ?
	`
	err := db.QueryRow(query, username).Scan(
		&user.Username,
		&user.FirstName,
		&user.Email,
		&user.DarkMode,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Fetch user's groups
	groupRows, err := db.Query(`
		SELECT id, name, created_at, updated_at
		FROM todo_groups
		WHERE username = ?
	`, username)
	if err != nil {
		return nil, err
	}
	defer groupRows.Close()

	for groupRows.Next() {
		var group parser.TodoGroup
		err := groupRows.Scan(&group.ID, &group.Name, &group.CreatedAt, &group.UpdatedAt)
		if err != nil {
			return nil, err
		}
		user.Groups = append(user.Groups, group)
	}

	// Fetch user's lists
	listRows, err := db.Query(`
		SELECT id, group_id, name, description, urgent, priority, done, status, created_at, updated_at
		FROM todo_lists
		WHERE username = ?
	`, username)
	if err != nil {
		return nil, err
	}
	defer listRows.Close()

	for listRows.Next() {
		var list parser.TodoList
		err := listRows.Scan(
			&list.ID,
			&list.GroupID,
			&list.Name,
			&list.Description,
			&list.Urgent,
			&list.Priority,
			&list.Done,
			&list.Status,
			&list.CreatedAt,
			&list.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		user.Lists = append(user.Lists, list)
	}

	// Set IsLoggedIn (use your own logic to check session or authentication)
	user.IsLoggedIn = true

	return &user, nil
}
