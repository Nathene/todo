package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

type Database struct {
	*sql.DB
}

const (
	not_started = "not started"
	in_progress = "in progress"
	completed   = "completed"
)

// InitDB initializes the SQLite database and creates tables if they don't exist
func InitDB() (*Database, error) {
	db, err := sql.Open("sqlite", "./todo_app.db")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return nil, err
	}

	createTables(db)

	return &Database{db}, nil
}
func createTables(db *sql.DB) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS todo_groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		username TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS todo_lists (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		group_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		urgent BOOLEAN NOT NULL,
		priority INTEGER DEFAULT 0,
		done BOOLEAN NOT NULL,
		status TEXT NOT NULL,
		username TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(group_id) REFERENCES todo_groups(id)
	);

	CREATE TABLE IF NOT EXISTS subtasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		todo_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		done BOOLEAN NOT NULL DEFAULT false,
		username TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(todo_id) REFERENCES todo_lists(id)
	);

	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(50) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		firstname VARCHAR(100) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		darkmode BOOLEAN DEFAULT false,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		username VARCHAR(255) NOT NULL,
		description TEXT,
		event_date DATE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(255) NOT NULL,
		title VARCHAR(255) NOT NULL,
		content TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (username) REFERENCES users(username)
	);

	CREATE TRIGGER IF NOT EXISTS update_notes_updated_at
	AFTER UPDATE ON notes
	FOR EACH ROW
	BEGIN
		UPDATE notes
		SET updated_at = CURRENT_TIMESTAMP
		WHERE id = OLD.id;
	END;
	`)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
}
