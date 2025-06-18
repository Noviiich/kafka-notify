package models

import (
	"time"

	resp "github.com/Noviiich/kafka-notify/internal/lib/api/response"
)

type Notification struct {
	ID        string    `json:"id,omitempty"`
	From      User      `json:"from"`
	To        User      `json:"to"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type Request struct {
	FromID  int    `json:"from_id" binding:"required"`
	ToID    int    `json:"to_id" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type Response struct {
	resp.Response
	NotificationID string `json:"notification_id,omitempty"`
}
