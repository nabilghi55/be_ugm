package storage

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"samplegadai-backend/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	Client     *minio.Client
	BucketName string
}

func NewMinioStorage(cfg *config.Config) *MinioStorage {
	minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		log.Printf("Warning: Failed to initialize MinIO client: %v", err)
		return &MinioStorage{Client: nil, BucketName: cfg.MinioBucketName}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	exists, err := minioClient.BucketExists(ctx, cfg.MinioBucketName)
	if err == nil && !exists {
		err = minioClient.MakeBucket(ctx, cfg.MinioBucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("Warning: Failed to create bucket %s: %v", cfg.MinioBucketName, err)
		} else {
			log.Printf("Created MinIO private bucket %s", cfg.MinioBucketName)
		}
	} else if err != nil {
		log.Printf("MinIO ping warning (%s): %v. Using fallback mode for local development.", cfg.MinioEndpoint, err)
	}

	return &MinioStorage{
		Client:     minioClient,
		BucketName: cfg.MinioBucketName,
	}
}

func (m *MinioStorage) GeneratePresignedUploadURL(ctx context.Context, objectKey string, expiresMinutes int) (string, error) {
	if m.Client == nil {
		return fmt.Sprintf("http://localhost:8080/api/v1/storage/mock-upload/%s", objectKey), nil
	}

	presignedURL, err := m.Client.PresignedPutObject(ctx, m.BucketName, objectKey, time.Duration(expiresMinutes)*time.Minute)
	if err != nil {
		log.Printf("MinIO PresignedPutObject warning: %v. Falling back to dev mock upload URL.", err)
		return fmt.Sprintf("http://localhost:8080/api/v1/storage/mock-upload/%s", objectKey), nil
	}
	return presignedURL.String(), nil
}

func (m *MinioStorage) GeneratePresignedViewURL(ctx context.Context, objectKey string, expiresMinutes int) (string, error) {
	if m.Client == nil || objectKey == "" {
		return "https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?auto=format&fit=crop&w=400&q=80", nil
	}

	reqParams := make(url.Values)
	presignedURL, err := m.Client.PresignedGetObject(ctx, m.BucketName, objectKey, time.Duration(expiresMinutes)*time.Minute, reqParams)
	if err != nil {
		log.Printf("MinIO PresignedGetObject warning: %v. Using fallback sample photo.", err)
		return "https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?auto=format&fit=crop&w=400&q=80", nil
	}
	return presignedURL.String(), nil
}
