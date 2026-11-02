package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/strahsvid17/shorturl/internal/db/postgres"
	"github.com/strahsvid17/shorturl/internal/lib"
)

type URLHandler struct {
	q *db.Queries
}

func NewURLHandler(q *db.Queries) *URLHandler {
	return &URLHandler{
		q: q,
	}
}

type aliasCreateBody struct {
	FullURL string `json:"full_url"`
}

type aliasCreateResponse struct {
	ShortURL string `json:"short_url"`
}

func (uh *URLHandler) FindURLAlias() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			http.Error(w, "empty alias not allowed", http.StatusBadRequest)
			return
		}

		res, err := uh.q.FindURLByAlias(context.Background(), alias)
		if err != nil {
			http.Error(w, "alias not found", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, res, http.StatusPermanentRedirect)
	}
}

func (uh *URLHandler) CreateURLAlias() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var aliasCreateBody aliasCreateBody

		if err := json.NewDecoder(r.Body).Decode(&aliasCreateBody); err != nil {
			http.Error(w, "invalid json", http.StatusNotAcceptable)
			return
		}

		if err := lib.ValidateURL(aliasCreateBody.FullURL); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		alias := lib.GenerateURLAlias(7)
		if _, err := uh.q.CreateURLAlias(
			context.Background(),
			db.CreateURLAliasParams{Url: aliasCreateBody.FullURL, Alias: alias},
		); err != nil {
			if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
				if pgErr.Code == "23505" {
					http.Error(w, "alias duplication", http.StatusConflict)
					return
				} else {
					http.Error(w, "some error", http.StatusInternalServerError)
					return
				}
			}
		}

		aliasCreateResponse := aliasCreateResponse{
			ShortURL: fmt.Sprintf("http://localhost:9090/api/%s", alias),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(aliasCreateResponse); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
