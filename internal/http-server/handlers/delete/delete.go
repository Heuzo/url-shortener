package delete

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
		const op = "handlers.url.redirect.Delete"
		requestId := middleware.GetReqID(request.Context())
		log = log.With(
			slog.String("operation", op),
			slog.String("request_id", requestId),
		)

		aliaToDelete := chi.URLParam(request, "alias")
		log.Info("Getted alias", "alias", aliaToDelete)
		if errDel := serviceDB.DeleteURL(aliaToDelete); errDel != nil {
			if errors.Is(errDel, storage.ErrAliasNotFound) {
				log.Error("failed to delete url", sl.Err(errDel), slog.String("alias", aliaToDelete))
				writer.WriteHeader(http.StatusNotFound)
				return
			}
			log.Error("failed to delete url", sl.Err(errDel), slog.String("alias", aliaToDelete))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
