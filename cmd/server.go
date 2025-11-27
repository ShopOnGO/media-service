package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ShopOnGO/ShopOnGO/pkg/kafkaService"
	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	"github.com/ShopOnGO/media-service/configs"
	"github.com/ShopOnGO/media-service/internal/media"
	"github.com/ShopOnGO/media-service/internal/redisdb"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

func main() {
	// migrations.CheckForMigrations()
	conf := configs.LoadConfig()

	kafkaProducers := kafkaService.InitKafkaProducers(
		conf.KafkaProducer.Brokers,
		conf.KafkaProducer.Topic,
	)

	redisClient := redisdb.NewRedisDB(conf)
	logger.Info("✅ Redis connected")

	// репозиторий для хранения (local/S3) через единый интерфейс
	var store media.Storage
	var err error

	switch conf.Media.StorageType {
	case "s3":
		store, err = media.NewS3Storage(
			conf.Media.S3Bucket,
			conf.Media.S3Region,
			conf.Media.S3Endpoint,
			conf.Media.S3AccessKey,
			conf.Media.S3SecretKey,
		)
		if err != nil {
			logger.Error("CRITICAL: Failed to initialize S3 storage", err.Error())
			os.Exit(1)
		}
		logger.Info("✅ S3 Storage initialized successfully")

	default:
		store = media.NewLocalStorage(conf.Media.LocalPath, conf.Media.BaseURL)
		logger.Info("✅ Local Storage initialized")
	}

	// service
	mediaSvc := media.NewMediaService(store, redisClient)

	// handler
	router := gin.Default()

	router.Static("/media", "./uploads")

	mediaHandler := media.NewMediaHandler(router, media.MediaHandlerDeps{
		Mediasvc: mediaSvc,
		Kafka:    kafkaProducers["media"],
	})

	go func() {
		if err := router.Run(":8084"); err != nil {
			fmt.Println("Ошибка при запуске HTTP-сервера:", err)
		}
	}()
	logger.Info("Media HTTP listening on 8084")

	kafkaConsumer := kafkaService.NewConsumer(
		conf.KafkaConsumer.Brokers,
		conf.KafkaConsumer.Topic,
		conf.KafkaConsumer.GroupID,
		conf.KafkaConsumer.ClientID,
	)
	defer kafkaConsumer.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Запускаем Kafka
	go func() {
		kafkaConsumer.Consume(ctx, func(msg kafka.Message) error {
			return mediaHandler.HandleMediaEvent(ctx, msg.Value)
		})
	}()

	<-ctx.Done()
	logger.Info("Shutting down media service...")
}
