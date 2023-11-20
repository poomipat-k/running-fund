package server

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/poomipat-k/running-fund/pkg/projects"
	server "github.com/poomipat-k/running-fund/pkg/server/handlers"
	"github.com/poomipat-k/running-fund/pkg/users"
)

type Server struct{}

func (app *Server) Routes(db *sql.DB) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	userStore := users.NewStore(db)
	userHandler := server.NewUserHandler(userStore)

	projectStore := projects.NewStore(db)
	projectHandler := server.NewProjectHandler(projectStore, userStore)

	mux.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("API landing page"))
		})

		r.Post("/project/reviewer", projectHandler.GetReviewerDashboard)
		r.Get("/project/review-period", projectHandler.GetReviewPeriod)
		r.Get("/project/review/{projectCode}", projectHandler.GetReviewerProjectDetails)
		r.Get("/review/criteria/{criteriaVersion}", projectHandler.GetProjectCriteria)

		r.Post("/project/review", projectHandler.AddReview)

		r.Get("/user/reviewers", userHandler.GetReviewers)
		r.Get("/user/reviewer", userHandler.GetReviewerById)
	})

	return mux
}
