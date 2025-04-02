package redirect

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/logger/sl"
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
		if errors.Is(errDB, storage.ErrURLNotFound) {
			log.Error("url not found by gien alias", "alias", alias)
			writer.WriteHeader(http.StatusNotFound)
		}
		if errDB != nil {
			log.Error("failed to get url by given alias:", sl.Err(errDB))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Location", url)
		writer.WriteHeader(http.StatusMovedPermanently)
		log.Info("successfully redirection", slog.String("url", url))
	}
}
