package controllers

import (
	"sample-service/internal/repository"
	"github.com/labstack/echo/v4"
	"strconv"
	"sample-service/internal/response"	
	"fmt"
	"sample-service/internal/model"
)

type UserController struct {
	repo repository.UserRepository
}

// NewUserController creates a new UserController
func NewUserController(repo repository.UserRepository) *UserController {
	return &UserController{
		repo: repo,
	}
}

// @Summary Get all users
// @Description Retrieve all users from the database
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users [get]
func (uc *UserController) GetAllUsers(ctx echo.Context) error {
	users, err := uc.repo.GetAllUsers()
	if err != nil {
		return response.JSONErrorResponse(ctx, "Failed to retrieve users", err.Error())
	}
	return response.JSONSuccessResponse(ctx, "Users retrieved successfully", users)
}

// @Summary Get user by ID
// @Description Retrieve a user by their ID
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{id} [get]
func (uc *UserController) GetUserByID(ctx echo.Context) error {
	id := ctx.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return response.JSONErrorResponse(ctx, "Failed to retrieve user", "Invalid user ID")
	}

	user, err := uc.repo.GetUserByID(userID)
	if err != nil {
		return response.JSONErrorResponse(ctx, "User not found", err.Error())
	}
	return response.JSONSuccessResponse(ctx, "User retrieved successfully", user)
}

// @Summary Create a new user
// @Description Create a new user in the database
// @Accept json
// @Produce json
// @Param user body model.User true "User details"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users [post]
func (uc *UserController) CreateUser(ctx echo.Context) error {
	var user model.User
	if err := ctx.Bind(&user); err != nil {
		return response.JSONErrorResponse(ctx, "Invalid request body", err.Error())
	}

	newUser, err := uc.repo.CreateUser(user)
	if err != nil {
        if err.Error() == fmt.Sprintf("username '%s' already exists", user.UserName) {
            return response.JSONErrorResponse(ctx, "Username already exists", err.Error())
        }
        return response.JSONErrorResponse(ctx, "Failed to create user", err.Error())
    }

	return response.JSONSuccessResponse(ctx, "User created successfully", newUser)
}
