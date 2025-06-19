package sqlite

import (
	"fmt"

	"github.com/Noviiich/kafka-notify/pkg/models"
)

func (s *Storage) Add(notification models.Notification) error {
	query := `
        INSERT INTO notifications (id, from_user_id, to_user_id, message, created_at)
        VALUES (?, ?, ?, ?, ?)
    `

	_, err := s.db.Exec(query,
		notification.ID,
		notification.From.ID,
		notification.To.ID,
		notification.Message,
		notification.Timestamp,
	)

	if err != nil {
		return fmt.Errorf("failed to add notification: %w", err)
	}

	return nil
}

func (s *Storage) GetByUserID(userID int) ([]models.Notification, error) {
	query := `
        SELECT 
            n.id, n.message, n.created_at,
            fu.id, fu.name,
            tu.id, tu.name
        FROM notifications n
        JOIN users fu ON n.from_user_id = fu.id
        JOIN users tu ON n.to_user_id = tu.id
        WHERE n.to_user_id = ?
        ORDER BY n.created_at DESC
    `

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		var fromUser, toUser models.User

		err := rows.Scan(
			&n.ID, &n.Message, &n.Timestamp,
			&fromUser.ID, &fromUser.Name,
			&toUser.ID, &toUser.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}

		n.From = fromUser
		n.To = toUser
		notifications = append(notifications, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return notifications, nil
}

// func (s *Storage) GetAll() ([]models.Notification, error) {
// 	// Аналогично GetByUserID, но без WHERE условия
// 	query := `
//         SELECT
//             n.id, n.message, n.created_at,
//             fu.id, fu.name,
//             tu.id, tu.name
//         FROM notifications n
//         JOIN users fu ON n.from_user_id = fu.id
//         JOIN users tu ON n.to_user_id = tu.id
//         ORDER BY n.created_at DESC
//     `

// 	// Реализация аналогична GetByUserID
// 	// ...
// 	return nil, nil
// }

func (s *Storage) GetByUserIDWithPagination(userID int, limit, offset int) ([]models.Notification, error) {
	query := `
        SELECT 
            n.id, n.message, n.created_at,
            fu.id, fu.name,
            tu.id, tu.name
        FROM notifications n
        JOIN users fu ON n.from_user_id = fu.id
        JOIN users tu ON n.to_user_id = tu.id
        WHERE n.to_user_id = ?
        ORDER BY n.created_at DESC
        LIMIT ? OFFSET ?
    `

	rows, err := s.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications with pagination: %w", err)
	}
	defer rows.Close()

	// Реализация сканирования аналогична GetByUserID
	// ...
	return nil, nil
}
