package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jimmyvallejo/gleamspeak-api/internal/api/v1/handlers"
	"github.com/jimmyvallejo/gleamspeak-api/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading it: %v", err)
	}

	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()

	APICfg := APIConfig{
		Port: port,
		DB:   dbQueries,
	}

	h := handlers.NewHandlers(APICfg.DB)
	// m := middleware.NewMiddleware(APICfg.DB)

	// Test readiness

	mux.HandleFunc("GET /v1/healthz", handlers.HandlerReadiness)
	mux.HandleFunc("GET /v1/err", handlers.HandlerError)

	// User Routes

	mux.HandleFunc("POST /v1/users", h.CreateUser)

	srv := &http.Server{
		Addr:    ":" + APICfg.Port,
		Handler: mux,
	}

	go func() {
		log.Printf("Serving on port: %s\n", APICfg.Port)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown signal received, exiting...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
}
