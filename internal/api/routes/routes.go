package routes

import (
	"net/http"
	"todo/internal/api/controller"
	"todo/internal/api/controller/account"
	"todo/internal/api/controller/calendar"
	"todo/internal/api/controller/settings"
	"todo/internal/api/controller/subtask"
	"todo/internal/api/notes"
	"todo/internal/db"

	"github.com/labstack/echo/v4"
)

func Handle(e *echo.Echo, db *db.Database) {
	// Public routes
	e.GET("/", controller.Dashboard(db), controller.AuthMiddleware(db))

	e.GET("/favicon.ico", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	e.GET("/login", controller.Login(db))
	e.POST("/login", controller.Login(db))
	e.GET("/create-account", account.CreatePage())
	e.POST("/create-account", account.CreateHandler(db))
	e.GET("/logout", controller.Logout())

	// Authenticated routes
	protected := e.Group("")
	protected.Use(controller.AuthMiddleware(db))

	// Settings and account management
	protected.GET("/settings", settings.Page())
	protected.POST("/settings/update/firstname", account.UpdateFirstName(db))
	protected.POST("/settings/update/username", account.UpdateUsername(db))
	protected.POST("/settings/update/email", account.UpdateEmail(db))
	protected.POST("/settings/update/password", account.UpdatePassword(db))
	protected.POST("/settings/update/darkmode", account.UpdateDarkMode(db))
	protected.POST("/darkmode", controller.ToggleDarkMode(db))

	// Todo lists and groups
	protected.POST("/groups", controller.CreateTodoGroup(db))
	protected.GET("/groups", controller.GetTodoGroups(db))
	protected.GET("/groups/:id", controller.GetTodoListsByGroupID(db))
	protected.POST("/groups/:id", controller.CreateTodoList(db))
	protected.POST("/todos/:id/status", controller.UpdateTodoStatus(db))

	protected.POST("/todos/:id/status", controller.UpdateTodoStatus(db))
	protected.POST("/todos/:id/priority", controller.UpdateTodoPriority(db))
	protected.POST("/todos/:id/delete", controller.DeleteTodo(db))
	protected.GET("/todos/:id/edit", controller.GetEditTodoPage(db))
	protected.POST("/todos/:id/update", controller.UpdateTodoDetails(db))

	// Subtasks
	protected.POST("/todos/:todo_id/subtasks", controller.CreateSubtask(db))
	protected.POST("/subtasks/:id/toggle", subtask.ToggleSubtaskDone(db))
	protected.POST("/groups/:group_name/:todo_id", controller.CreateSubtask(db))
	protected.GET("/groups/:group_name/:todo_id", controller.GetSubtasks(db))

	// Calendar
	protected.GET("/calendar", calendar.GetAll(db))
	protected.POST("/calendar/add", calendar.AddEvent(db))

	// Notes
	protected.GET("/notes", notes.GetNotes(db))
	protected.GET("/notes/:id/edit", notes.EditNotePage(db))
	protected.POST("/notes/:id/edit", notes.UpdateNote(db))
	protected.POST("/notes/:id/delete", notes.DeleteNote(db))
	protected.POST("/notes/add", notes.AddNote(db))

}
