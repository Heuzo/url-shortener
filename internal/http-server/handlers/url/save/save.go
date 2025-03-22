package save

import (
	"errors"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	cfg "url-shortener/internal/config"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
}

// New is a function to create handler to saving URL
func New(log *slog.Logger, serviceDB storage.SQLService, config *cfg.Config) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.url.save.New"
		requestId := middleware.GetReqID(request.Context())
		log = log.With(
			slog.String("operation", op),
			slog.String("request_id", requestId),
		)
		var req Request
		err := render.DecodeJSON(request.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(writer, request, response.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))
		var validateErr validator.ValidationErrors
		// Валидируем поля в соответствии со struct-тегами
		if errValid := validator.New().Struct(req); errValid != nil {
			if errors.As(errValid, &validateErr) {
				log.Error("invalid request", sl.Err(validateErr))
				render.JSON(writer, request, response.ValidationError(validateErr))
				return
			}
			log.Error("error while validating: ", sl.Err(errValid))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = generateUniqueAlias(log, serviceDB, requestId, config.AliasLenght)
		}

		id, err := serviceDB.SaveURL(req.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("url already exists", slog.String("url", req.URL))
				render.JSON(writer, request, response.Error("url already exist"))
				return
			}
			log.Error("failed to save url", sl.Err(err))
			render.JSON(writer, request, response.Error("failed to save url"))
			return
		}
		render.JSON(writer, request, response.OK(alias))
		log.Info("url successfully saved", slog.Int64("id", id))
	}
}

func generateUniqueAlias(log *slog.Logger, serviceDB storage.SQLService, reqID string, aliasLength int) string {
	const op = "handlers.url.save.generateUniqueAlias"
	log = log.With(
		slog.String("operation", op),
		slog.String("request_id", reqID),
	)
	for {
		alias := random.NewRandomString(aliasLength)
		exists, err := serviceDB.AliasExists(alias)
		if err != nil {
			log.Error("failed to check is alias exists", sl.Err(err))
		}
		if exists {
			continue
		}
		return alias
	}
}
