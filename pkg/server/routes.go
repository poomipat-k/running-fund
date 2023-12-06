package server

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	appMiddleware "github.com/poomipat-k/running-fund/pkg/middleware"
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
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "withCredentials"},
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

		r.Post("/project/reviewer", appMiddleware.IsReviewer(projectHandler.GetReviewerDashboard))
		r.Get("/project/review-period", appMiddleware.IsReviewer(projectHandler.GetReviewPeriod))
		r.Get("/project/review/{projectCode}", appMiddleware.IsReviewer(projectHandler.GetReviewerProjectDetails))
		r.Get("/review/criteria/{criteriaVersion}", appMiddleware.IsReviewer(projectHandler.GetProjectCriteria))

		r.Post("/project/review", appMiddleware.IsReviewer(reviewHandler.AddReview))

		r.Get("/user/current", appMiddleware.IsLoggedIn(userHandler.GetCurrentUser))

		r.Post("/auth/register", userHandler.SignUp)
		r.Post("/auth/login", userHandler.SignIn)
	})

	return mux
}
