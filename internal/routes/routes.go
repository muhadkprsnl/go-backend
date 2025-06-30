package routes

import (
	// "local/qa-report/internal/controllers"
	// "local/qa-report/internal/repositories"
	// "local/qa-report/pkg/middleware"
	"net/http"

	"github.com/muhadkprsnl/go-backend/interna/pkg/middleware"
	"github.com/muhadkprsnl/go-backend/internal/controllers"
	"github.com/muhadkprsnl/go-backend/internal/repositories"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func SetupRouter(mongoClient *mongo.Client, logger *zap.Logger) *mux.Router {
	router := mux.NewRouter()

	// Initialize database and repositories
	db := mongoClient.Database("QA")
	reportRepo := repositories.NewReportRepository(db, logger)
	chartRepo := repositories.NewChartRepository(db, logger)

	// Initialize controllers
	reportController := controllers.NewReportController(reportRepo, logger)
	chartController := controllers.NewChartController(chartRepo, logger)
	sprintController := controllers.NewSprintController(db.Collection("Report_2")) // ✅ ADD THIS
	summaryController := controllers.NewSummaryController(reportRepo)
	releaseController := controllers.NewReleaseController(reportRepo)
	authController := controllers.NewAuthController(mongoClient)

	// Middleware stack
	router.Use(middleware.CORS)                 // CORS for all routes
	router.Use(middleware.Logging(logger))      // Request logging
	router.Use(middleware.RecoverPanic(logger)) // Panic recovery

	// API routes with versioning
	api := router.PathPrefix("/api").Subrouter()
	v1 := api.PathPrefix("/v1").Subrouter()

	// Report endpoints (versioned)
	reportRouter := v1.PathPrefix("/reports").Subrouter()
	reportRouter.HandleFunc("", reportController.CreateReport).Methods("POST")
	reportRouter.HandleFunc("", reportController.GetAllReports).Methods("GET")
	reportRouter.HandleFunc("/{id}", reportController.UpdateReport).Methods("PUT", "PATCH", "OPTIONS")
	reportRouter.HandleFunc("/{id}", reportController.DeleteReport).Methods("DELETE", "OPTIONS")

	// Analytics endpoints (versioned)
	v1.HandleFunc("/analytics/donut", chartController.GetDonutChart).Methods("GET", "OPTIONS")

	// ✅ Sprint endpoint
	router.HandleFunc("/api/v1/sprints", sprintController.GetSprints).Methods("GET", "OPTIONS")

	// RoutesofSummary
	v1.HandleFunc("/summary", summaryController.GetSummary).Methods("GET", "OPTIONS")

	//SprintError
	router.HandleFunc("/api/v1/sprint-error-comparison", summaryController.GetSprintErrorComparison).Methods("GET")

	//REleasecontroller
	router.HandleFunc("/api/v1/releases", releaseController.GetReleases).Methods("GET", "OPTIONS")

	// Legacy endpoints (unversioned) — for dev/prod form submission
	api.HandleFunc("/devform", reportController.CreateReport).Methods("POST", "OPTIONS")
	api.HandleFunc("/prodform", reportController.CreateReport).Methods("POST", "OPTIONS")

	// System health endpoints
	router.HandleFunc("/health", healthCheck).Methods("GET")
	router.HandleFunc("/ready", readinessCheck).Methods("GET")

	// Auth routes
	// router.HandleFunc("/api/auth/register", authController.Register).Methods("POST", "OPTIONS")
	// router.HandleFunc("/api/auth/login", authController.Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/login", authController.Login).Methods("POST", "OPTIONS")

	return router
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

func readinessCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ready"}`))
}
