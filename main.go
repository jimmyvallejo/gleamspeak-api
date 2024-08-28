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

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jimmyvallejo/gleamspeak-api/internal/api/middleware"
	"github.com/jimmyvallejo/gleamspeak-api/internal/api/routes"
	"github.com/jimmyvallejo/gleamspeak-api/internal/api/v1/handlers"
	"github.com/jimmyvallejo/gleamspeak-api/internal/database"
	"github.com/jimmyvallejo/gleamspeak-api/internal/redis"
	"github.com/jimmyvallejo/gleamspeak-api/internal/websocket"
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
	awsAccess := os.Getenv("AWS_ACCESS_KEY")
	awsSecret := os.Getenv("AWS_SECRET_KEY")

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

	credProvider := credentials.NewStaticCredentialsProvider(awsAccess, awsSecret, "")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credProvider),
		config.WithRegion("us-east-1"),
	)

	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	})

	apiCfg := APIConfig{
		Port:      port,
		DB:        dbQueries,
		RDB:       rdb,
		JwtSecret: jwtSecret,
		S3:        s3Client,
	}

	h := handlers.NewHandlers(apiCfg.DB, apiCfg.JwtSecret, apiCfg.S3)
	m := middleware.NewMiddleware(apiCfg.DB, apiCfg.RDB, apiCfg.JwtSecret)
	w := websocket.NewManager(apiCfg.DB, apiCfg.RDB, h)

	apiCfg.Handlers = h

	router := routes.NewRouter(h, m, w)
	router.SetupV1Routes()

	handler := c.Handler(router.GetHandler())

	srv := &http.Server{
		Addr:    ":" + apiCfg.Port,
		Handler: handler,
	}

	go func() {
		log.Printf("Serving on port: %s\n", apiCfg.Port)
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
