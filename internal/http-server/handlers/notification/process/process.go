package process

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	resp "github.com/Noviiich/kafka-notify/internal/lib/api/response"
	"github.com/Noviiich/kafka-notify/internal/lib/logger/sl"
	"github.com/Noviiich/kafka-notify/internal/storage"
	"github.com/Noviiich/kafka-notify/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type NotificationGetter interface {
	GetByUserID(userID int) ([]models.Notification, error)
}

func New(log *slog.Logger, notificationGetter NotificationGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.notification.process.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userStr := chi.URLParam(r, "userID")
		if userStr == "" {
			log.Info("userID is empty")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		userID, err := strconv.Atoi(userStr)
		if err != nil {
			log.Error("failed to parse userID", sl.Err(err))
			render.JSON(w, r, resp.Error("invalid userID"))
			return
		}

		// Process the notification
		result, err := notificationGetter.GetByUserID(userID)
		if err != nil {
			if errors.Is(err, storage.ErrNoNotificationsFound) {
				log.Info("no notifications found", sl.Err(err))
				render.JSON(w, r, resp.Error("no notifications found"))
				return
			}
			log.Error("failed to process notification", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to process notification"))
			return
		}
		log.Info("notifications processed successfully", slog.Int("userID", userID))
		render.JSON(w, r, resp.Success(result))
	}
}
