package notes

import (
	"net/http"
	"todo/internal/db"
	"todo/internal/parser"

	"github.com/labstack/echo/v4"
)

// GetNotes fetches all notes for the logged-in user
func GetNotes(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(parser.User)
		if !ok || !user.IsLoggedIn {
			return c.Redirect(http.StatusFound, "/login")
		}

		rows, err := db.Query("SELECT id, title, content, updated_at FROM notes WHERE username = ? ORDER BY updated_at DESC", user.Username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch notes"})
		}
		defer rows.Close()

		var notes []parser.Note
		for rows.Next() {
			var note parser.Note
			if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.UpdatedAt); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse notes"})
			}
			notes = append(notes, note)
		}

		return c.Render(http.StatusOK, "notes/notes.gohtml", map[string]interface{}{
			"Notes":      notes,
			"darkMode":   c.Get("darkMode"),
			"isLoggedIn": user.IsLoggedIn,
		})
	}
}

// EditNotePage renders the edit note page
func EditNotePage(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(parser.User)
		if !ok || !user.IsLoggedIn {
			return c.Redirect(http.StatusFound, "/login")
		}

		noteID := c.Param("id")
		var note parser.Note
		err := db.QueryRow("SELECT id, title, content FROM notes WHERE id = ? AND username = ?", noteID, user.Username).
			Scan(&note.ID, &note.Title, &note.Content)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Note not found"})
		}

		return c.Render(http.StatusOK, "notes/edit.gohtml", map[string]interface{}{
			"Note":       note,
			"darkMode":   c.Get("darkMode"),
			"isLoggedIn": user.IsLoggedIn,
		})
	}
}

// UpdateNote updates an existing note
func UpdateNote(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(parser.User)
		if !ok || !user.IsLoggedIn {
			return c.Redirect(http.StatusFound, "/login")
		}

		noteID := c.Param("id")
		title := c.FormValue("title")
		content := c.FormValue("content")

		if title == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title cannot be empty"})
		}

		_, err := db.Exec("UPDATE notes SET title = ?, content = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND username = ?", title, content, noteID, user.Username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update note"})
		}

		return c.Redirect(http.StatusSeeOther, "/notes")
	}
}

// DeleteNote deletes a note
func DeleteNote(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(parser.User)
		if !ok || !user.IsLoggedIn {
			return c.Redirect(http.StatusFound, "/login")
		}

		noteID := c.Param("id")

		_, err := db.Exec("DELETE FROM notes WHERE id = ? AND username = ?", noteID, user.Username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete note"})
		}

		return c.Redirect(http.StatusSeeOther, "/notes")
	}
}

// AddNote adds a new note
func AddNote(db *db.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(parser.User)
		if !ok || !user.IsLoggedIn {
			return c.Redirect(http.StatusFound, "/login")
		}

		title := c.FormValue("title")
		content := c.FormValue("content")

		if title == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title cannot be empty"})
		}

		_, err := db.Exec("INSERT INTO notes (username, title, content) VALUES (?, ?, ?)", user.Username, title, content)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add note"})
		}

		return c.Redirect(http.StatusSeeOther, "/notes")
	}
}
