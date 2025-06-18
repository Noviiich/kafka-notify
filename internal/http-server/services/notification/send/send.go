package send

import (
	"time"

	"github.com/Noviiich/kafka-notify/internal/storage"
	"github.com/Noviiich/kafka-notify/pkg/models"
	"github.com/google/uuid"
)

type Producer interface {
	SendNotification(notification models.Notification) error
	Close() error
}

type UserRepository interface {
	User(id int) (models.User, error)
}

type NotificationService struct {
	UserRepository
	Producer
}

func New(userRepo UserRepository, producer Producer) *NotificationService {
	return &NotificationService{
		UserRepository: userRepo,
		Producer:       producer,
	}
}

func (s *NotificationService) SendNotification(fromID int, toID int, message string) (string, error) {
	// Валидация пользователей
	fromUser, err := s.UserRepository.User(fromID)
	if err != nil {
		return "", storage.ErrUserNotFound
	}

	toUser, err := s.UserRepository.User(toID)
	if err != nil {
		return "", storage.ErrUserNotFound
	}

	// Создание уведомления
	notification := models.Notification{
		ID:        uuid.New().String(),
		From:      fromUser,
		To:        toUser,
		Message:   message,
		Timestamp: time.Now(),
	}

	if err := s.Producer.SendNotification(notification); err != nil {
		return "", err
	}

	return notification.ID, nil
}
