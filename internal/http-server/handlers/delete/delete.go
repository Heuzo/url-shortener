package delete

import (
	"errors"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"
)

type Request struct {
	Alias string `json:"alias" validate:"required"`
}

func New(log *slog.Logger, serviceDB storage.SQLService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.url.redirect.Delete"
		requestId := middleware.GetReqID(request.Context())
		log = log.With(
			slog.String("operation", op),
			slog.String("request_id", requestId),
		)
		var req Request
		err := render.DecodeJSON(request.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		// Валидируем поля в соответствии со struct-тегами
		var validateErr validator.ValidationErrors
		if errValid := validator.New().Struct(req); errValid != nil {
			if errors.As(errValid, &validateErr) {
				log.Error("invalid request", sl.Err(errValid))
				render.JSON(writer, request, response.ValidationError(validateErr))
				return
			}
			log.Error("error while validating: ", sl.Err(errValid))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		urlToDelete := req.Alias
		if errDel := serviceDB.DeleteURL(urlToDelete); errDel != nil {
			if errors.Is(errDel, storage.ErrAliasNotFound) {
				log.Error("failed to delete url", sl.Err(errDel), slog.String("alias", req.Alias))
				writer.WriteHeader(http.StatusNotFound)
				return
			}
			log.Error("failed to delete url", sl.Err(err), slog.String("alias", req.Alias))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
