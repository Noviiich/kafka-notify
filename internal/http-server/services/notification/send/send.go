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

func (s *NotificationService) SendNotification(FromID int, ToID int, Message string) (string, error) {
	// Валидация пользователей
	fromUser, err := s.UserRepository.User(FromID)
	if err != nil {
		return "", storage.ErrUserNotFound
	}

	toUser, err := s.UserRepository.User(ToID)
	if err != nil {
		return "", storage.ErrUserNotFound
	}

	// Создание уведомления
	notification := models.Notification{
		ID:        uuid.New().String(),
		From:      fromUser,
		To:        toUser,
		Message:   Message,
		Timestamp: time.Now(),
	}

	if err := s.Producer.SendNotification(notification); err != nil {
		return "", err
	}

	return notification.ID, nil
}
