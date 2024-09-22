package interfaces

import (
	"context"
	"mime/multipart"
)

type IFileStorage interface {
	GetFileURL(ctx context.Context, bucketName string, fileName string) (string, error)
	SaveFile(ctx context.Context, bucketName string, fileName string, file multipart.File) error
	FileExists(ctx context.Context, bucketName string, fileName string) (bool, error)
	CreateFileName(ctx context.Context, file multipart.File) (string, error)
}
