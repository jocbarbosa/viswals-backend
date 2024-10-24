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

	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Server panicked: %v", r)
		}
	}()

	err := godotenv.Load()
	if err != nil {
		log.Println("error loading .env file, using default values")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("error initializing zap logger")
	}
	defer zapLogger.Sync()
	loggerAdapter := logger.NewZapAdapter(zapLogger)

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("failed to migrate the database: %v", err)
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
		log.Fatalf("failed to initialize RabbitMQ: %v", err)
	}
	defer rabbitMqAdapter.Close()

	userRepo := repository.NewUserRepository(db)

	userSvc := services.NewUserService(loggerAdapter, userRepo, redisAdapter, rabbitMqAdapter)
	go userSvc.StartConsuming(ctx)

	userController := controllers.NewUserController(userSvc, loggerAdapter, redisAdapter)

	router := api.NewRouter(userController)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("starting API server on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not listen on :%s: %v", port, err)
		}
	}()

	<-ch
	log.Println("shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server exiting")
}
