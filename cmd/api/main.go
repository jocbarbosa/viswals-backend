package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jocbarbosa/viswals-backend/internals/adapter/cache"
	"github.com/jocbarbosa/viswals-backend/internals/adapter/logger"
	"github.com/jocbarbosa/viswals-backend/internals/adapter/messaging/rabbitmq"
	"github.com/jocbarbosa/viswals-backend/internals/adapter/repository"
	"github.com/jocbarbosa/viswals-backend/internals/application/api"
	"github.com/jocbarbosa/viswals-backend/internals/application/controllers"
	"github.com/jocbarbosa/viswals-backend/internals/core/model"
	"github.com/jocbarbosa/viswals-backend/internals/core/services"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartAPIServer() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default values")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Error initializing zap logger")
	}
	defer zapLogger.Sync()
	loggerAdapter := logger.NewZapAdapter(zapLogger)

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("Failed to migrate the database: %v", err)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR is not set")
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0
	redisAdapter := cache.NewRedisClient(redisAddr, redisPassword, redisDB)

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	rabbitMQQueue := os.Getenv("RABBITMQ_QUEUE")
	if rabbitMQURL == "" || rabbitMQQueue == "" {
		log.Fatal("RABBITMQ_URL or RABBITMQ_QUEUE is not set")
	}
	rabbitMqAdapter, err := rabbitmq.NewRabbitMQAdapter(rabbitMQURL, rabbitMQQueue)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rabbitMqAdapter.Close()

	userRepo := repository.NewUserRepository(db)
	userController := controllers.NewUserController(userRepo, loggerAdapter, redisAdapter)

	userSvc := services.NewUserService(loggerAdapter, userRepo, redisAdapter, rabbitMqAdapter)
	go userSvc.StartConsuming(ctx)

	router := api.NewRouter(userController)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Starting API server on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :%s: %v", port, err)
		}
	}()

	<-ch
	log.Println("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
