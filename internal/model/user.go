package model

// User represents a user in the system
type User struct {
	ID       	 int64  `json:"user_id"`
	UserName     string `json:"user_name"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	UserStatus   string `json:"user_status"`
	Department   string `json:"department"`
}


