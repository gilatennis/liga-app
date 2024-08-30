package minio

import (
	"context"
	"fmt"
	"io"
	"liga-app/config"
	"mime/multipart"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	*minio.Client
}

func NewClient(c *config.Config) (*Client, error) {
	minioClient, err := minio.New(c.Minio.URL, &minio.Options{
		Creds:  credentials.NewStaticV4(c.Minio.AccessKey, c.Minio.SecretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, fmt.Errorf("minio.New(): failed to create new minio client: %w", err)
	}

	return &Client{minioClient}, nil
}

func (c *Client) SaveFile(ctx context.Context, bucketName string, fileName string, file multipart.File) error {
	var size int64 = -1
	_, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("infra.minio.client.SaveFile(): %w", err)
	}
	size, err = file.Seek(0, io.SeekCurrent)
	if err != nil {
		return fmt.Errorf("infra.minio.client.SaveFile(): %w", err)
	}

	file.Seek(0, io.SeekStart)

	_, err = c.PutObject(ctx, bucketName, fileName, file, size, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return fmt.Errorf("infra.minio.client.SaveFile(): %w", err)
	}

	return nil
}

func (c *Client) GetFileURL(ctx context.Context, bucketName string, fileName string) (string, error) {
	presignedURL, err := c.PresignedGetObject(ctx, bucketName, fileName, time.Second*24*60*60, nil)
	if err != nil {
		return "", fmt.Errorf("infra.minio.client.GetFileURL(): %w", err)
	}

	return presignedURL.String(), nil
}

func (c *Client) FileExists(ctx context.Context, bucketName string, fileName string) (bool, error) {
	_, err := c.StatObject(ctx, bucketName, fileName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("infra.minio.client.FileExists(): %w", err)
	}

	return true, nil
}
