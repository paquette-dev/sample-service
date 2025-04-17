package routes

import (
    "github.com/labstack/echo/v4"
    "sample-service/internal/controllers"
    "sample-service/internal/repository"
    "database/sql"
)

// RegisterUserRoutes registers the user routes
func RegisterUserRoutes(e *echo.Echo, db *sql.DB) {
    userRepo := repository.NewUserRepository(db)
    userController := controllers.NewUserController(userRepo)

    e.GET("/users", userController.GetAllUsers)
    e.GET("/users/:id", userController.GetUserByID)
    e.POST("/users", userController.CreateUser)
    e.PUT("/users/:id", userController.UpdateUser)
}