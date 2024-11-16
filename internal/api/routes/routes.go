package routes

import (
	"todo/internal/api/controller"
	"todo/internal/db"

	"github.com/labstack/echo/v4"
)

func Handle(e *echo.Echo, db *db.Database) {
	// Public routes
	e.POST("/login", controller.Login())

	// Authenticated routes
	todoGroup := e.Group("/todo")
	todoGroup.Use(controller.AuthMiddleware)

	// Todo lists and groups
	todoGroup.POST("/list", controller.CreateTodoList(db))
	todoGroup.POST("/groups", controller.CreateTodoGroup(db))
	todoGroup.GET("/groups", controller.GetTodoGroups(db))
	todoGroup.GET("/groups/:name", controller.GetTodoListsByGroupName(db))
	todoGroup.POST("/groups/:name", controller.AddTodoListToGroup(db))

	// Subtasks
	todoGroup.POST("/groups/:group_name/:todo_id", controller.CreateSubtask(db))
	todoGroup.GET("/groups/:group_name/:todo_id", controller.GetSubtasks(db))
}
