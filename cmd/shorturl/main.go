package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	db "github.com/strahsvid17/shorturl/internal/db/postgres"
	"github.com/strahsvid17/shorturl/internal/http-server/handlers"
)

func init() {
	_ = godotenv.Load(".env")
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	pool, err := pgxpool.New(
		context.Background(),
		os.Getenv("DB_URL"),
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	queries := db.New(pool)
	urlHandler := handlers.NewURLHandler(queries)

	r.Route("/api", func(r chi.Router) {
		r.Get("/{alias}", urlHandler.FindURLAlias())
		r.Group(func(r chi.Router) {
			r.Use(middleware.BasicAuth(
				"urlshortener",
				map[string]string{
					os.Getenv("AUTH_USERNAME"): os.Getenv("AUTH_PASSWORD"),
				}),
			)
			r.Post("/", urlHandler.CreateURLAlias())
		})
	})

	server := http.Server{
		Addr:        ":9090",
		Handler:     r,
		IdleTimeout: time.Second * 60,
	}

	log.Fatal(server.ListenAndServe())
}
