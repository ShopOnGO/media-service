package media

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type Storage interface {
  Save(file *multipart.FileHeader) (string, error)
  GenerateURL(key string) string
}

type LocalStorage struct {
  BasePath string // пример: "./uploads"
  BaseURL  string // пример: "http://localhost:8084/media"
}

func NewLocalStorage(basePath, baseURL string) *LocalStorage {
  return &LocalStorage{BasePath: basePath, BaseURL: baseURL}
}

func (s *LocalStorage) Save(file *multipart.FileHeader) (string, error) {
  key := file.Filename
  dest := filepath.Join(s.BasePath, key)
  if err := os.MkdirAll(s.BasePath, 0755); err != nil {
      return "", err
  }
  if err := saveUploadedFile(file, dest); err != nil {
      return "", err
  }
  return key, nil
}

func (s *LocalStorage) GenerateURL(key string) string {
  return fmt.Sprintf("%s%s", s.BaseURL, key)
}

func saveUploadedFile(file *multipart.FileHeader, dest string) error {
  src, err := file.Open()
  if err != nil {
      return err
  }
  defer src.Close()

  out, err := os.Create(dest)
  if err != nil {
      return err
  }
  defer out.Close()

  _, err = io.Copy(out, src)
  return err
}