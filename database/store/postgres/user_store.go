package postgres

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/lib/pq"
	"github.com/ranabd36/project-qa/database/store"
	"github.com/ranabd36/project-qa/pb"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (s *Store) FindByEmail(email string) (*pb.User, error) {
	const statement = `SELECT id, first_name, last_name, username, email,password,is_active, is_admin, created_at, update_at FROM users where email = $1;`
	return s.selectUser(statement, email)
}

func (s *Store) FindByUsername(username string) (*pb.User, error) {
	const statement = `SELECT id, first_name, last_name, username, email,password,is_active, is_admin, created_at, update_at FROM users where username = $1;`
	return s.selectUser(statement, username)
}

func (s *Store) ToggleActive(id int32) error {
	const updateStatement = `Update users set is_active = not is_active where id = $1;`
	return s.executeStatement(updateStatement, id)
}

func (s *Store) ToggleAdmin(id int32) error {
	const updateStatement = `Update users set is_admin = not is_admin where id = $1;`
	return s.executeStatement(updateStatement, id)
}

func (s *Store) UpdatePassword(id int32, newPassword string) error {
	hasPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	const updateStatement = `Update users set password = $2 where id = $1;`
	return s.executeStatement(updateStatement, id, hasPassword)
}

func (s *Store) Delete(id int32) error {
	const deleteStatement = `DELETE from users where id = $1;`
	return s.executeStatement(deleteStatement, id)
}

func (s *Store) Update(user *pb.User) error {
	const updateStatement = `Update users set first_name = $2, last_name = $3 where id = $1;`
	return s.executeStatement(updateStatement, user.GetId(), user.GetFirstName(), user.GetLastName())
}

func (s *Store) Find(id int32) (*pb.User, error) {
	const statement = `SELECT id, first_name, last_name, username, email, password, is_active, is_admin, created_at, update_at FROM users where id = $1;`
	return s.selectUser(statement, id)
}

func (s *Store) Save(user *pb.User) error {
	hasPassword, err := bcrypt.GenerateFromPassword([]byte(user.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	const insertStatement = `INSERT INTO users (first_name, last_name, username, email, password, is_active, is_admin) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	
	err = s.db.QueryRow(insertStatement,
		user.FirstName,
		user.FirstName,
		user.Username,
		user.Email,
		hasPassword,
		user.IsActive,
		user.IsAdmin,
	).Scan(&user.Id)
	
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == "23505" {
			return store.ErrAlreadyExists
		}
		return fmt.Errorf("failed to save row: %w", err)
	}
	return nil
}

func (s *Store) selectUser(statement string, args ...interface{}) (*pb.User, error) {
	user := &pb.User{}
	var createdAt time.Time
	var updatedAt time.Time
	if err := s.db.QueryRow(statement, args...).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.IsActive,
		&user.IsAdmin,
		&createdAt,
		&updatedAt,
	); err != nil {
		return nil, err
	}
	c, _ := ptypes.TimestampProto(createdAt)
	u, _ := ptypes.TimestampProto(updatedAt)
	user.CreatedAt = c
	user.UpdatedAt = u
	return user, nil
}

func (s *Store) executeStatement(statement string, args ...interface{}) error {
	rows, err := s.db.Exec(statement, args...)
	if err != nil {
		return err
	}
	
	updateCount, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	
	if updateCount > 0 {
		return nil
	}
	return errors.New("unknown error")
}
