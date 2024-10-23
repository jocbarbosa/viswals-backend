package reader

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	adapters "github.com/jocbarbosa/viswals-backend/internals/adapter/logger"
	rabbitmq "github.com/jocbarbosa/viswals-backend/internals/adapter/messaging/rabbitmq"
	"github.com/jocbarbosa/viswals-backend/internals/application/controllers"
)

func StartReader() {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	filepath := os.Getenv("USERS_FILE_PATH")
	if filepath == "" {
		log.Fatal("filepath is not set")
	}

	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Error initializing zap logger")
	}
	defer zapLogger.Sync()

	logger := adapters.NewZapAdapter(zapLogger)

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("RabbitMQ URL is not set")
	}

	queueName := os.Getenv("RABBITMQ_QUEUE")
	if queueName == "" {
		log.Fatal("RabbitMQ queue is not set")
	}

	messaging, err := rabbitmq.NewRabbitMQAdapter(rabbitURL, queueName)
	if err != nil {
		log.Fatal("Error initializing RabbitMQ adapter:", err)
	}
	defer messaging.Close()

	r := controllers.NewFileReader(logger, messaging, filepath)

	err = r.ReadFile(ctx)
	if err != nil {
		logger.Error("Error reading file", err)
	}
}
