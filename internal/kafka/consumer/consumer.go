package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/Noviiich/kafka-notify/internal/config"
	"github.com/Noviiich/kafka-notify/pkg/models"
)

type Consumer struct {
	consumerGroup sarama.ConsumerGroup
	topic         string
	NotificationRepository
}

type NotificationRepository interface {
	Add(notification models.Notification) error
}

func New(cfg config.KafkaConfig, notificationRepository NotificationRepository) (*Consumer, error) {
	const op = "kafka.consumer.New"

	config := sarama.NewConfig()

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.ConsumerGroup, config)
	if err != nil {
		return nil, fmt.Errorf("%s: create consumer group: %w", op, err)
	}

	return &Consumer{
		consumerGroup:          consumerGroup,
		topic:                  cfg.Topic,
		NotificationRepository: notificationRepository,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	handler := &consumerGroupHandler{NotificationRepository: c.NotificationRepository}

	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer context cancelled")
			return c.consumerGroup.Close()
		default:
			err := c.consumerGroup.Consume(ctx, []string{c.topic}, handler)
			if err != nil {
				log.Printf("Error from consumer: %v", err)
				return err
			}
		}
	}
}

func (c *Consumer) Close() error {
	return c.consumerGroup.Close()
}

type consumerGroupHandler struct {
	NotificationRepository
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		userID := string(message.Key)

		var notification models.Notification
		if err := json.Unmarshal(message.Value, &notification); err != nil {
			log.Printf("Failed to unmarshal notification: %v", err)
			session.MarkMessage(message, "")
			continue
		}

		// Валидация userID
		if _, err := strconv.Atoi(userID); err != nil {
			log.Printf("Invalid user ID: %s", userID)
			session.MarkMessage(message, "")
			continue
		}

		h.NotificationRepository.Add(notification)
		log.Printf("Notification added for user %s: %s", userID, notification.Message)

		session.MarkMessage(message, "")
	}
	return nil
}
