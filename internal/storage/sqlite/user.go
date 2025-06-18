package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Noviiich/kafka-notify/internal/storage"
	"github.com/Noviiich/kafka-notify/pkg/models"
	_ "github.com/mattn/go-sqlite3"
)

func (s *Storage) User(id int) (models.User, error) {
	const op = "storage.sqlite.User"

	stmt, err := s.db.Prepare("SELECT id, name FROM users WHERE id = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var user models.User

	err = stmt.QueryRow(id).Scan(&user.ID, &user.Name)
	if err != nil {
		// обработка пустого результата
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return user, nil
}

func (s *Storage) Users() ([]models.User, error) {
	const op = "storage.sqlite.Users"

	stmt, err := s.db.Prepare("SELECT id, name FROM users")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: execute query: %w", op, err)
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: scan row: %w", op, err)
		}
		users = append(users, user)
	}

	// Проверяем на ошибки, которые могли возникнуть во время итерации
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows iteration: %w", op, err)
	}

	return users, nil
}
