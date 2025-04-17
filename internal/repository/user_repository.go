package repository

import (
	"database/sql"
	"sample-service/internal/model"
	"fmt"
)

type UserRepository interface {
	GetAllUsers() ([]model.User, error)
	GetUserByID(id int) (*model.User, error)
	CheckIfUsernameExists(username string) (bool, error)
	CreateUser(user model.User) (*model.User, error)
	UpdateUser(user model.User) (*model.User, error)
}

type userRepo struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

// GetAllUsers retrieves all users from the database
func (r *userRepo) GetAllUsers() ([]model.User, error) {
	rows, err := r.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.Department, &user.UserStatus)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUserByID retrieves a user by their ID from the database
func (r *userRepo) GetUserByID(id int) (*model.User, error) {
	row := r.db.QueryRow("SELECT * FROM users WHERE user_id = ?", id)

	var user model.User
	err := row.Scan(&user.ID, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.Department, &user.UserStatus)
	if err != nil {
		return nil, err
	}	

	return &user, nil
}

// CheckIfUsernameExists checks if a username exists in the database
func (r *userRepo) CheckIfUsernameExists(username string) (bool, error) {
    var exists bool
    err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE user_name = ?", username).Scan(&exists)
    if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

        return false, fmt.Errorf("failed to check username existence: %v", err)
    }
    return exists, nil
}

// CreateUser creates a new user in the database
func (r *userRepo) CreateUser(user model.User) (*model.User, error) {
	exists, err := r.CheckIfUsernameExists(user.UserName)
    if err != nil {
        return nil, fmt.Errorf("error checking username: %w", err)
    }
    
    if exists {
        return nil, fmt.Errorf("username '%s' already exists", user.UserName)
    }
	
	result, err := r.db.Exec("INSERT INTO users (user_name, first_name, last_name, email, department, user_status) VALUES (?, ?, ?, ?, ?, ?)",
		user.UserName, user.FirstName, user.LastName, user.Email, user.Department, user.UserStatus)
	if err != nil {
		return nil, err	
	}

	userID, err := result.LastInsertId()
	user.ID = userID

    return &user, nil
}

// UpdateUser updates a user in the database
func (r *userRepo) UpdateUser(user model.User) (*model.User, error) {
	// Check if user exists
	_, err := r.GetUserByID(int(user.ID))
	if err != nil {
		return nil, fmt.Errorf("user with ID %d not found: %w", user.ID, err)
	}
	
	// Update the user
	_, err = r.db.Exec(
		"UPDATE users SET user_name = ?, first_name = ?, last_name = ?, email = ?, department = ?, user_status = ? WHERE user_id = ?",
		user.UserName, user.FirstName, user.LastName, user.Email, user.Department, user.UserStatus, user.ID)
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}
