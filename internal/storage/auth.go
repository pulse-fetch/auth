package storage

import (
	"auth/internal/domain"
	"database/sql"
	"fmt"
)

func (s *Storage) Create(username, email, hashPass string) error {
	var op = "storage.auth.Create"
	if _, err := s.db.Exec("INSERT INTO users(username, email, password) values($1, $2, $3)", username, email, hashPass); err != nil {
		return fmt.Errorf("%s:%v", op, err)

	}
	return nil
}

func (s *Storage) Get(data interface{}, method string) (domain.User, error) {
	var op = "storage.auth.Get"
	var user domain.User
	var row *sql.Row
	switch method {
	case "name":
		row = s.db.QueryRow("SELECT * FROM users WHERE username=$1", data)
	case "id":
		row = s.db.QueryRow("SELECT * FROM users WHERE id=$1", data)
	default:
		return domain.User{}, fmt.Errorf("%s: Undefined finding method", op)
	}

	if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.HashPassword); err != nil {

		return domain.User{}, fmt.Errorf("%s:%v", op, err)
	}
	return user, nil
}

func (s *Storage) Del(id int64) error {
	var op = "storage.auth.Del"
	res, err := s.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {

		return fmt.Errorf("%s:%v", op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {

		return fmt.Errorf("%s:%v", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: Invalid username", op)
	}

	return nil
}

func (s *Storage) UpdateName(id int64, newName string) error {
	var op = "storage.auth.UpdateNameUser"
	res, err := s.db.Exec("UPDATE users SET username = $1 WHERE id = $2", newName, id)
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: Invalid id", op)
	}

	return nil
}
