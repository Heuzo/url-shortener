package redirect

import (
	"github.com/go-chi/chi/middleware"
	"log/slog"
	"net/http"
	cfg "url-shortener/internal/config"
	"url-shortener/internal/storage"
)

func New(log *slog.Logger, serviceDB storage.SQLService, config *cfg.Config) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.url.redirect.New"
		requestId := middleware.GetReqID(request.Context())
		log = log.With(
			slog.String("operation", op),
			slog.String("request_id", requestId),
		)
	}
}
