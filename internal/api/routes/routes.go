package routes

import (
	"todo/internal/api/controller"
	"todo/internal/api/controller/account"
	"todo/internal/api/controller/settings"
	"todo/internal/api/controller/subtask"
	"todo/internal/db"

	"github.com/labstack/echo/v4"
)

func Handle(e *echo.Echo, db *db.Database) {
	// Public routes
	e.GET("/", controller.LandingPage(), controller.AuthMiddleware(db))

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

	// protected.GET("/groups/:name", controller.GetTodoListsByGroupName(db))
	// protected.POST("/groups/:name", controller.AddTodoListToGroup(db))
	protected.POST("/todos/:todo_id/subtasks", controller.CreateSubtask(db))

	protected.POST("/subtasks/:id/toggle", subtask.ToggleSubtaskDone(db))
	// Subtasks
	protected.POST("/groups/:group_name/:todo_id", controller.CreateSubtask(db))
	protected.GET("/groups/:group_name/:todo_id", controller.GetSubtasks(db))
}
