package redirect

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
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

		alias := request.URL.Path[1:]

		url, errDB := serviceDB.GetURL(alias)
		log.Info("Getted url", "url", url)
		if errDB != nil {
			slog.Error("failed to get url:", errDB)
			render.JSON(writer, request, errDB)
			return
		}
		writer.Header().Set("Location", url)
		writer.WriteHeader(http.StatusMovedPermanently)
		log.Info("successfully redirection", slog.String("url", url))
	}
}
