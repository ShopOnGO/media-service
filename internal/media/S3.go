package media

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	Bucket   string
	Endpoint string
	Client   *s3.Client
}

func NewS3Storage(bucket, region, endpoint, accessKey, secretKey string) (*S3Storage, error) {
	creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(creds),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	return &S3Storage{
		Bucket:   bucket,
		Endpoint: endpoint,
		Client:   client,
	}, nil
}

func (s *S3Storage) Save(file *multipart.FileHeader) (string, error) {
	key := file.Filename
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	_, err = s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(key),
		Body:        src,
		ContentType: aws.String(file.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", err
	}

	return key, nil
}

func (s *S3Storage) GenerateURL(key string) string {
	
	publicEndpoint := strings.Replace(s.Endpoint, "/s3", "/object/public", 1)

	return fmt.Sprintf("%s/%s/%s", publicEndpoint, s.Bucket, key)
}
