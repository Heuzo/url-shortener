package redirect

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log/slog"
	"net/http"
	"url-shortener/internal/storage"
)

func New(log *slog.Logger, serviceDB storage.SQLService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.url.redirect.New"
		requestId := middleware.GetReqID(request.Context())
		log = log.With(
			slog.String("operation", op),
			slog.String("request_id", requestId),
		)

		alias := chi.URLParam(request, "alias")

		url, errDB := serviceDB.GetURL(alias)
		log.Info("Getted alias", "alias", alias)
		if errDB != nil {
			slog.Error("failed to get url by given alias:", errDB)
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		writer.Header().Set("Location", url)
		writer.WriteHeader(http.StatusMovedPermanently)
		log.Info("successfully redirection", slog.String("url", url))
	}
}
