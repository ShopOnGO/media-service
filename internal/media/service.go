package media

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	"github.com/ShopOnGO/media-service/internal/redisdb"
)

const mediaCacheTTL = 24 * time.Hour

type MediaService struct {
	storage Storage
	redis   *redisdb.RedisDB
}

func NewMediaService(s Storage, r *redisdb.RedisDB) *MediaService {
	return &MediaService{
		storage: s,
		redis:   r,
	}
}

func (s *MediaService) UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	key, err := s.storage.Save(file)
	if err != nil {
		return "", err
	}

	fullURL := s.storage.GenerateURL(key)

	err = s.redis.Client.Set(ctx, "media:"+key, fullURL, mediaCacheTTL).Err()
	if err != nil {
		logger.Error("Failed to cache in Redis", err)
	}

	return fullURL, nil
}

func (s *MediaService) GenerateURL(ctx context.Context, key string) string {
	val, err := s.redis.Client.Get(ctx, "media:" + key).Result()
	if err == nil {
		return val
	}

	url := s.storage.GenerateURL(key)

	_ = s.redis.Client.Set(ctx, "media:"+key, url, mediaCacheTTL).Err()

	return url
}
