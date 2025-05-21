package main

import (
	"context"
	"github.com/blacktalenthubs/go-service-api/database"
	"github.com/blacktalenthubs/go-service-api/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

// Middleware for logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the request
		log.Printf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Database configuration
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "consultancy"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Initialize database
	db, err := database.New(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize handlers
	consultantHandler := handlers.NewConsultantHandler(db)
	skillHandler := handlers.NewSkillHandler(db)

	// Initialize router
	r := mux.NewRouter()

	// Apply middleware
	r.Use(loggingMiddleware)

	// API routes
	apiRouter := r.PathPrefix("/api").Subrouter()

	// Consultant routes
	apiRouter.HandleFunc("/consultants", consultantHandler.GetAll).Methods("GET")
	apiRouter.HandleFunc("/consultants/{id:[0-9]+}", consultantHandler.Get).Methods("GET")
	apiRouter.HandleFunc("/consultants", consultantHandler.Create).Methods("POST")
	apiRouter.HandleFunc("/consultants/{id:[0-9]+}", consultantHandler.Update).Methods("PUT")
	apiRouter.HandleFunc("/consultants/{id:[0-9]+}", consultantHandler.Delete).Methods("DELETE")
	apiRouter.HandleFunc("/consultants/skills/{skill_id:[0-9]+}", consultantHandler.GetBySkill).Methods("GET")

	// Skill routes
	apiRouter.HandleFunc("/skills", skillHandler.GetAll).Methods("GET")
	apiRouter.HandleFunc("/skills/{id:[0-9]+}", skillHandler.Get).Methods("GET")
	apiRouter.HandleFunc("/skills", skillHandler.Create).Methods("POST")
	apiRouter.HandleFunc("/skills/{id:[0-9]+}", skillHandler.Update).Methods("PUT")
	apiRouter.HandleFunc("/skills/{id:[0-9]+}", skillHandler.Delete).Methods("DELETE")

	// Start server with graceful shutdown
	startServerWithGracefulShutdown(r)
}

// Helper function to get environment variable with default
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get environment variable as int with default
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func startServerWithGracefulShutdown(r *mux.Router) {
	// Define server
	srv := &http.Server{
		Addr:         ":" + getEnv("PORT", "8080"),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// Channel for server errors
	serverErrors := make(chan error, 1)

	// Start server
	go func() {
		log.Printf("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// Channel for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Wait for interrupt signal or server error
	select {
	case err := <-serverErrors:
		log.Fatalf("Server error: %v", err)
	case <-stop:
		log.Println("Shutting down server...")

		// Create a deadline for server shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}

		log.Println("Server gracefully stopped")
	}
}
