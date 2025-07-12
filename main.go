package main

import (
	"context"
	"fmt"

	// "local/qa-report/internal/config"
	// "local/qa-report/internal/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/muhadkprsnl/go-backend/internal/config"
	"github.com/muhadkprsnl/go-backend/internal/routes"
	"go.uber.org/zap"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found (that's OK in production)")
	}

	fmt.Println("Loaded URI:", os.Getenv("MONGODB_URI"))

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("Failed to sync logger: %v", err)
		}
	}()

	// Initialize MongoDB connection
	mongoClient, err := config.ConnectMongoDB()
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(ctx); err != nil {
			logger.Error("Failed to disconnect MongoDB", zap.Error(err))
		}
	}()

	// Initialize router
	router := routes.SetupRouter(mongoClient, logger)

	// Test route for confirmation
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"pong"}`))
	}).Methods("GET")

	// Get dynamic port (Render will provide this via PORT env var)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001" // Default for local development
	}

	// Configure HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for OS signals for graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		logger.Info("🚀 Starting server", zap.String("address", server.Addr))
		fmt.Printf("Server listening at http://localhost:%s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("❌ Server failed", zap.Error(err))
		}
	}()

	// Wait for termination signal
	<-done
	logger.Info("🛑 Shutting down server...")

	// Shutdown with context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("❌ Server shutdown failed", zap.Error(err))
	} else {
		logger.Info("✅ Server stopped gracefully")
	}
}

// package main

// import (
// 	"context"
// 	"fmt"
// 	"local/qa-report/internal/config"
// 	"local/qa-report/internal/routes"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"go.uber.org/zap"
// )

// func main() {
// 	// Initialize logger
// 	logger, err := zap.NewProduction()
// 	if err != nil {
// 		log.Fatalf("Failed to initialize logger: %v", err)
// 	}
// 	defer func() {
// 		if err := logger.Sync(); err != nil {
// 			log.Printf("Failed to sync logger: %v", err)
// 		}
// 	}()

// 	// Initialize MongoDB connection
// 	mongoClient, err := config.ConnectMongoDB()
// 	if err != nil {
// 		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
// 	}
// 	defer func() {
// 		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 		defer cancel()
// 		if err := mongoClient.Disconnect(ctx); err != nil {
// 			logger.Error("Failed to disconnect MongoDB", zap.Error(err))
// 		}
// 	}()

// 	// Initialize router from routes package
// 	router := routes.SetupRouter(mongoClient, logger)

// 	// 🔍 Add test route for confirmation
// 	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte(`{"message":"pong"}`))
// 	}).Methods("GET")

// 	// Configure HTTP server
// 	server := &http.Server{
// 		Addr:         ":3001",
// 		Handler:      router,
// 		ReadTimeout:  15 * time.Second,
// 		WriteTimeout: 15 * time.Second,
// 		IdleTimeout:  60 * time.Second,
// 	}

// 	// Graceful shutdown channel
// 	done := make(chan os.Signal, 1)
// 	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

// 	// Start server in goroutine
// 	go func() {
// 		logger.Info("🚀 Starting server", zap.String("address", server.Addr))
// 		fmt.Println("Server listening at http://localhost:3001")
// 		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			logger.Fatal("❌ Server failed", zap.Error(err))
// 		}
// 	}()

// 	// Wait for shutdown signal
// 	<-done
// 	logger.Info("🛑 Shutting down server...")

// 	// Create shutdown context with timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// Attempt graceful shutdown
// 	if err := server.Shutdown(ctx); err != nil {
// 		logger.Error("❌ Server shutdown failed", zap.Error(err))
// 	} else {
// 		logger.Info("✅ Server stopped gracefully")
// 	}
// }
