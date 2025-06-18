package producer

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/Noviiich/kafka-notify/internal/config"
	"github.com/Noviiich/kafka-notify/pkg/models"
)

type Producer struct {
	saramaProducer sarama.SyncProducer
	topic          string
}

func New(cfg config.KafkaConfig) (*Producer, error) {
	const op = "kafka.producer.New"

	// Инициализирует новую конфигурацию Kafka по умолчанию
	config := sarama.NewConfig()
	// Обеспечивает получение подтверждения от отправителя после успешного сохранения сообщения
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3

	// Инициализация синхронного producer Kafka
	saramaProducer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("%s: create producer: %w", op, err)
	}

	return &Producer{
		saramaProducer: saramaProducer,
		topic:          cfg.Topic,
	}, nil
}

func (p *Producer) SendNotification(notification models.Notification) error {
	const op = "kafka.producer.SendNotification"
	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("%s: failed to marshal notification: %w", op, err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(strconv.Itoa(notification.To.ID)),
		Value: sarama.StringEncoder(notificationJSON),
	}

	_, _, err = p.saramaProducer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("%s: failed to send kafka message: %w", op, err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.saramaProducer.Close()
}
