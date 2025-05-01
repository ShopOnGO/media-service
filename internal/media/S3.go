package media

import (
    "context"
    "fmt"
    "mime/multipart"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Storage struct {
    Bucket string
    Client *s3.Client
}

func NewS3Storage(bucket string, region string) (*S3Storage, error) {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
    if err != nil {
        return nil, err
    }
    client := s3.NewFromConfig(cfg)
    return &S3Storage{
        Bucket: bucket,
        Client: client,
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
        Bucket: aws.String(s.Bucket),
        Key:    aws.String(key),
        Body:   src,
        ContentType: aws.String(file.Header.Get("Content-Type")),
        ACL:    types.ObjectCannedACLPublicRead,
    })
    if err != nil {
        return "", err
    }

    return key, nil
}

func (s *S3Storage) GenerateURL(key string) string {
    return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.Bucket, key)
}
