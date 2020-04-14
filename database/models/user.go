package models

import (
	"github.com/ranabd36/project-qa/database"
	"time"
)

type User struct {
	Id        int
	FirstName string
	LastName  string
	Username  string
	Email     string
	Password  string
	IsActive  bool
	IsAdmin   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (user *User) GetUserById(id int) *User {
	err := database.Connection.QueryRow("SELECT * FROM users where id = $1", id).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.IsActive,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		return nil
	}
	
	return user
}

func (user *User) Create() (int, error) {
	statement := `INSERT INTO users (first_name, last_name, username, email, password, is_active, is_admin) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	id := 0
	err := database.Connection.QueryRow(statement,
		user.FirstName,
		user.FirstName,
		user.Username,
		user.Email,
		user.Password,
		user.IsActive,
		user.IsAdmin,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	
	return id, nil
}
