package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/patrickmn/go-cache"

	"github.com/poomipat-k/running-fund/pkg/address"
	"github.com/poomipat-k/running-fund/pkg/captcha"
	appEmail "github.com/poomipat-k/running-fund/pkg/email"
	mw "github.com/poomipat-k/running-fund/pkg/middleware"
	"github.com/poomipat-k/running-fund/pkg/projects"
	"github.com/poomipat-k/running-fund/pkg/review"
	s3Service "github.com/poomipat-k/running-fund/pkg/s3-service"

	"github.com/poomipat-k/running-fund/pkg/users"
)

type Server struct{}

type BucketBasics struct {
	S3Client *s3.Client
}

func (app *Server) Routes(db *sql.DB) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "withCredentials", "Content-Disposition"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// config.WithRegion("ap-southeast-1")
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		log.Fatal()
	}
	s3Client := s3.NewFromConfig(sdkConfig)
	serverS3Service := s3Service.S3Service{S3Client: s3Client}

	presignClient := s3.NewPresignClient(s3Client)
	presigner := s3Service.Presigner{
		PresignClient: presignClient,
	}
	s3Handler := s3Service.NewS3Handler(presigner)

	emailService := appEmail.NewEmailService()

	userStore := users.NewStore(db, emailService)
	userHandler := users.NewUserHandler(userStore)

	reviewStore := review.NewStore(db)
	reviewHandler := review.NewProjectHandler(reviewStore, userStore)

	c := cache.New(3*time.Minute, 5*time.Minute)
	projectStore := projects.NewStore(db, c, serverS3Service)
	projectHandler := projects.NewProjectHandler(projectStore, userStore, serverS3Service)

	captchaStore := captcha.NewStore(c)
	captchaHandler := captcha.NewCaptchaHandler(captchaStore)

	addressStore := address.NewStore(db, c)
	addressHandler := address.NewAddressHandler(addressStore)

	mux.Route("/api/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("API landing page"))
		})

		r.Get("/review/criteria/{criteriaVersion}", mw.IsLoggedIn(projectHandler.GetProjectCriteria))
		r.Get("/applicant/criteria/{applicantCriteriaVersion}", mw.IsApplicant(projectHandler.GetApplicantCriteria))
		r.Get("/applicant/project/details/{projectCode}", mw.IsLoggedIn(projectHandler.GetApplicantProjectDetails))

		r.Post("/project/reviewer", mw.IsReviewer(projectHandler.GetReviewerDashboard))
		r.Get("/project/review-period", mw.IsLoggedIn(projectHandler.GetReviewPeriod))
		r.Post("/project/review/{projectCode}", mw.IsLoggedIn(projectHandler.GetReviewerProjectDetails))
		r.Post("/project", mw.IsApplicant(projectHandler.AddProject))
		r.Post("/project/addition-files", mw.IsLoggedIn(projectHandler.AddProjectAdditionFiles))
		r.Get("/project/applicant/dashboard", mw.IsApplicant(projectHandler.GetAllProjectDashboardByApplicantId))

		r.Post("/admin/project/{projectCode}", mw.IsAdmin(projectHandler.AdminUpdateProject))
		r.Post("/admin/dashboard/summary", mw.IsAdmin(projectHandler.GetAdminSummary))
		r.Post("/admin/dashboard/request", mw.IsAdmin(projectHandler.GetAdminRequestDashboard))
		r.Post("/admin/dashboard/started", mw.IsAdmin(projectHandler.GetAdminStartedDashboard))
		r.Post("/admin/report", mw.IsAdmin(projectHandler.GenerateAdminReport))

		r.Post("/project/review", mw.IsReviewer(reviewHandler.AddReview))

		r.Post("/user/activate-email", userHandler.ActivateUser)
		r.Post("/user/password/forgot", mw.ValidateCaptcha(userHandler.ForgotPassword, captchaStore))
		r.Post("/user/password/reset", userHandler.ResetPassword)

		r.Get("/auth/current", mw.IsLoggedIn(userHandler.GetCurrentUser))
		r.Get("/user/full-name/{userId}", mw.IsLoggedIn(userHandler.GetUserFullNameById))
		r.Post("/auth/register", userHandler.SignUp)
		r.Post("/auth/login", userHandler.SignIn)
		r.Post("/auth/logout", userHandler.SignOut)
		r.Post("/auth/refresh-token", userHandler.RefreshAccessToken)

		r.Post("/captcha/generate", captchaHandler.GenerateCaptcha)

		r.Get("/address/provinces", mw.IsLoggedIn(addressHandler.GetProvinces))
		r.Get("/address/districts/{provinceId}", mw.IsLoggedIn(addressHandler.GetDistrictsByProvince))
		r.Get("/address/subdistricts/{districtId}", mw.IsLoggedIn(addressHandler.GetSubdistrictsByProvince))
		r.Get("/address/postcodes/{subdistrictId}", mw.IsLoggedIn(addressHandler.GetPostcodeBySubdistrict))

		r.Post("/s3/presigned", mw.IsLoggedIn(s3Handler.GeneratePresignedUrl))
		r.Post("/s3/objects", mw.IsLoggedIn(projectHandler.ListApplicantFiles))
	})

	return mux
}
