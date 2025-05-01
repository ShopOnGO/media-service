package configs

import (
	"os"
	"strings"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	"github.com/joho/godotenv"
)

type Config struct {
	Media MediaConfig
	// Db DbConfig
	KafkaConsumer KafkaConsumerConfig
	KafkaProducer KafkaProducerConfig
}

// type DbConfig struct {
// 	Dsn string
// }

type MediaConfig struct {
	Port         string
	StorageType  string // "local" or "s3"
	LocalPath    string // path where to store locally
	BaseURL      string // public URL prefix for local media
	S3Bucket    string
	S3Region    string
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
}

type KafkaConsumerConfig struct {
	Brokers []string
	Topic   string
	GroupID string
	ClientID string
}

type KafkaProducerConfig struct {
	Brokers []string
	Topic   map[string]string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file, using default config", err.Error())
	}

	brokersRaw := os.Getenv("KAFKA_BROKERS")
	brokers := strings.Split(brokersRaw, ",")

	return &Config{
		Media: MediaConfig{
			Port:        getEnv("MEDIA_PORT", "8080"),
			StorageType: getEnv("MEDIA_STORAGE", "local"),
			LocalPath:   getEnv("MEDIA_LOCAL_PATH", "./uploads"),
			BaseURL:     getEnv("MEDIA_BASE_URL", "http://localhost:8080/media/"),
			S3Bucket:    os.Getenv("MEDIA_S3_BUCKET"),
			S3Region:    os.Getenv("MEDIA_S3_REGION"),
			S3Endpoint:  os.Getenv("MEDIA_S3_ENDPOINT"),
			S3AccessKey: os.Getenv("MEDIA_S3_ACCESS_KEY"),
			S3SecretKey: os.Getenv("MEDIA_S3_SECRET_KEY"),
		},
		// Db: DbConfig{
		// 	Dsn: os.Getenv("DSN"),
		// },
		KafkaConsumer: KafkaConsumerConfig{
			Brokers: brokers,
			Topic:   os.Getenv("KAFKA_CONSUMER_TOPIC"),
			GroupID: os.Getenv("KAFKA_CONSUMER_GROUP_ID"),
			ClientID: os.Getenv("KAFKA_CONSUMER_CLIENT_ID"),
		},
		KafkaProducer: KafkaProducerConfig{
			Brokers: brokers,
			Topic:   parseKafkaTopics(os.Getenv("KAFKA_PRODUCER_TOPIC")),
		},
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
	  return v
	}
	return def
}

func parseKafkaTopics(s string) map[string]string {
	topics := map[string]string{}
	pairs := strings.Split(s, ",")
	for _, p := range pairs {
		kv := strings.SplitN(p, ":", 2)
		if len(kv) == 2 {
			topics[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return topics
}