package media

import (
    "mime/multipart"
)

type MediaService struct {
    storage Storage
}

func NewMediaService(s Storage) *MediaService {
    return &MediaService{storage: s}
}

func (s *MediaService) UploadFile(file *multipart.FileHeader) (string, error) {
    f, err := file.Open()
    if err != nil {
        return "", err
    }
    defer f.Close()

    key, err := s.storage.Save(file)
    if err != nil {
        return "", err
    }
    return s.storage.GenerateURL(key), nil
}

func (s *MediaService) GenerateURL(key string) string {
	return s.storage.GenerateURL(key)
}
