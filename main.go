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

	"github.com/jimmyvallejo/gleamspeak-api/internal/api/middleware"
	"github.com/jimmyvallejo/gleamspeak-api/internal/api/v1/handlers"
	"github.com/jimmyvallejo/gleamspeak-api/internal/database"
	"github.com/jimmyvallejo/gleamspeak-api/internal/redis"
	"github.com/rs/cors"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading it: %v", err)
	}

	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB")
	jwtSecret := os.Getenv("JWT_SECRET")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	dbQueries := database.New(db)

	rdb, err := redis.NewClient()
	if err != nil {
		log.Print("Redis failed to initialize")
	}

	mux := http.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	APICfg := APIConfig{
		Port:      port,
		DB:        dbQueries,
		RDB:       rdb,
		JwtSecret: jwtSecret,
	}

	h := handlers.NewHandlers(APICfg.DB, APICfg.JwtSecret)
	m := middleware.NewMiddleware(APICfg.DB, APICfg.RDB, APICfg.JwtSecret)

	// Test readiness

	mux.HandleFunc("GET /v1/healthz", handlers.HandlerReadiness)
	mux.HandleFunc("GET /v1/err", handlers.HandlerError)

	// Auth Routes

	mux.HandleFunc("POST /v1/login", h.LoginUserStandard)
	mux.HandleFunc("POST /v1/logout", h.LogoutUserStandard)
	mux.HandleFunc("GET /v1/auth", m.IsAuthenticated(h.CheckAuthStatus))

	// User Routes

	mux.HandleFunc("POST /v1/users", h.CreateUserStandard)
	mux.HandleFunc("PUT /v1/users", m.IsAuthenticated(h.UpdateUser))

	// Server Routes

	mux.HandleFunc("POST /v1/servers", m.IsAuthenticated(h.CreateServer))
	mux.HandleFunc("POST /v1/servers/join", m.IsAuthenticated(h.JoinServer))
	mux.HandleFunc("DELETE /v1/servers/user", m.IsAuthenticated(h.LeaveServer))
	mux.HandleFunc("GET /v1/servers/user/many", m.IsAuthenticated(h.GetUserServers))

	// Text Channel Routes
	mux.HandleFunc("POST /v1/channels/text", m.IsAuthenticated(h.CreateTextChannel))
	mux.HandleFunc("GET /v1/channels/{serverID}", m.IsAuthenticated(h.GetServerTextChannels))

	//Token Routes
	mux.HandleFunc("POST /v1/refresh", h.RefreshToken)

	srv := &http.Server{
		Addr:    ":" + APICfg.Port,
		Handler: handler,
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
