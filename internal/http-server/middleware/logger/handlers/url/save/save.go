package save

import (
	"errors"
	"net/http"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	resp "url-shortener/internal/lib/logger/api/response"
	"url-shortener/internal/lib/logger/random"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 6

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		entry := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			entry.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		entry.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			entry.Error("invalid request", sl.Err(err))
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			entry.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, resp.Error("url already exists"))
			return
		}

		if err != nil {
			entry.Error("failed to add url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}

		entry.Info("url added", slog.Int64("id", id))
		responseOk(w, r, alias)
	}
}

func responseOk(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
