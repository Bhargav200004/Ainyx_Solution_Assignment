package routes

import (
	"github.com/gofiber/fiber/v3"

	"go-backend/internal/handler"
)

// Setup registers all application routes on the Fiber app.
func Setup(app *fiber.App, userHandler *handler.UserHandler) {
	// User CRUD routes.
	app.Post("/users", userHandler.CreateUser)
	app.Get("/users", userHandler.ListUsers)
	app.Get("/users/:id", userHandler.GetUserByID)
	app.Put("/users/:id", userHandler.UpdateUser)
	app.Delete("/users/:id", userHandler.DeleteUser)
}
