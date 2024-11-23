package db

import (
	"fmt"
	"todo/internal/parser"
)

func (db *Database) UpdateSubtaskDone(subtaskID int, done bool) error {
	_, err := db.Exec("UPDATE subtasks SET done = ? WHERE id = ?", done, subtaskID)
	if err != nil {
		return fmt.Errorf("failed to update subtask done status: %v", err)
	}
	return nil
}

func (db *Database) GetTodoIDBySubtask(subtaskID int) (int, error) {
	var todoID int
	err := db.QueryRow("SELECT todo_id FROM subtasks WHERE id = ?", subtaskID).Scan(&todoID)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve todo ID for subtask: %v", err)
	}
	return todoID, nil
}

func (db *Database) GetSubtasksByTodoID(todoID int) ([]parser.Subtask, error) {
	rows, err := db.Query("SELECT id, todo_id, name, done FROM subtasks WHERE todo_id = ?", todoID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve subtasks: %v", err)
	}
	defer rows.Close()

	var subtasks []parser.Subtask
	for rows.Next() {
		var subtask parser.Subtask
		if err := rows.Scan(&subtask.ID, &subtask.TodoID, &subtask.Name, &subtask.Done); err != nil {
			return nil, fmt.Errorf("failed to parse subtask: %v", err)
		}
		subtasks = append(subtasks, subtask)
	}
	return subtasks, nil
}

func (db *Database) AreAllSubtasksDone(todoID int) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM subtasks WHERE todo_id = ? AND done = 0", todoID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check subtasks for todo: %v", err)
	}
	return count == 0, nil
}
