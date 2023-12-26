package server

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	appEmail "github.com/poomipat-k/running-fund/pkg/email"
	mw "github.com/poomipat-k/running-fund/pkg/middleware"
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

	emailService := appEmail.NewEmailService()
	userStore := users.NewStore(db, emailService)
	userHandler := users.NewUserHandler(userStore)

	reviewStore := review.NewStore(db)
	reviewHandler := review.NewProjectHandler(reviewStore, userStore)

	projectStore := projects.NewStore(db)
	projectHandler := projects.NewProjectHandler(projectStore, userStore)

	mux.Route("/api/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("API landing page"))
		})

		r.Post("/project/reviewer", mw.IsReviewer(projectHandler.GetReviewerDashboard))
		r.Get("/project/review-period", mw.IsReviewer(projectHandler.GetReviewPeriod))
		r.Get("/project/review/{projectCode}", mw.IsReviewer(projectHandler.GetReviewerProjectDetails))
		r.Get("/review/criteria/{criteriaVersion}", mw.IsReviewer(projectHandler.GetProjectCriteria))

		r.Post("/project/review", mw.IsReviewer(reviewHandler.AddReview))

		r.Get("/user/activate-email", userHandler.ActivateUser)
		r.Post("/user/password/forgot", userHandler.ForgotPassword)
		r.Post("/user/password/reset", userHandler.ResetPassword)

		r.Get("/auth/current", mw.IsLoggedIn(userHandler.GetCurrentUser))
		r.Post("/auth/register", userHandler.SignUp)
		r.Post("/auth/login", userHandler.SignIn)
		r.Post("/auth/logout", userHandler.SignOut)
		r.Post("/auth/refresh-token", userHandler.RefreshAccessToken)
	})

	return mux
}
