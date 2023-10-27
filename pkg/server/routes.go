package server

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/poomipat-k/running-fund/pkg/projects"
	server "github.com/poomipat-k/running-fund/pkg/server/handlers"
)

type Server struct {
	DB *sql.DB
}

func (app *Server) Routes() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// store := projects.NewMemStore(app.DB)
	store := projects.NewStore(app.DB)
	projectHandler := server.NewProjectHandler(store)

	mux.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("API landing page"))
		})

		r.Get("/projects", projectHandler.GetAll)
	})

	return mux
}
