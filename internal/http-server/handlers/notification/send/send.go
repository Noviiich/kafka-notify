package send

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	resp "github.com/Noviiich/kafka-notify/internal/lib/api/response"
	"github.com/Noviiich/kafka-notify/internal/lib/logger/sl"
	"github.com/Noviiich/kafka-notify/internal/storage"
	"github.com/Noviiich/kafka-notify/pkg/models"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type NotificationSender interface {
	SendNotification(fromID, toID int, message string) (string, error)
}

func New(log *slog.Logger, sender NotificationSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.notifiction.send.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req models.Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			// обработка ошибки при пустом теле запроса
			if errors.Is(err, io.EOF) {
				log.Error("request body is empty")
				render.JSON(w, r, resp.Error("empty request"))
				return
			}

			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// Проверка полей и отправка уведомления
		notificationID, err := sender.SendNotification(req.FromID, req.ToID, req.Message)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				log.Error("user not found", sl.Err(err))
				render.JSON(w, r, resp.Error("user not found"))
				return
			}
			log.Error("failed to send notification", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to send notification"))
			return
		}

		log.Info("notification sent successfully", slog.String("notification_id", notificationID))

		responseOK(w, r, notificationID)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, notificationID string) {
	render.JSON(w, r, models.Response{
		Response:       resp.OK(),
		NotificationID: notificationID,
	})
}
