package server

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	customMiddleware "github.com/poomipat-k/running-fund/pkg/middleware"
	"github.com/poomipat-k/running-fund/pkg/projects"
	"github.com/poomipat-k/running-fund/pkg/review"
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
	userHandler := users.NewUserHandler(userStore)

	reviewStore := review.NewStore(db)
	reviewHandler := review.NewProjectHandler(reviewStore, userStore)

	projectStore := projects.NewStore(db)
	projectHandler := projects.NewProjectHandler(projectStore, userStore)

	mux.Route("/api/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("API landing page"))
		})

		r.Post("/project/reviewer", customMiddleware.MyFirstMiddleWare(projectHandler.GetReviewerDashboard))
		r.Get("/project/review-period", projectHandler.GetReviewPeriod)
		r.Get("/project/review/{projectCode}", customMiddleware.MyFirstMiddleWare(projectHandler.GetReviewerProjectDetails))
		r.Get("/review/criteria/{criteriaVersion}", projectHandler.GetProjectCriteria)

		r.Post("/project/review", reviewHandler.AddReview)

		r.Get("/user/reviewers", userHandler.GetReviewers)
		r.Get("/user/reviewer", userHandler.GetReviewerById)

		r.Post("/auth/register", userHandler.SignUp)
		r.Post("/auth/login", userHandler.SignIn)
	})

	return mux
}
