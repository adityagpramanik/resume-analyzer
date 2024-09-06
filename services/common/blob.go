package commonservices

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go"
)

var (
	once     sync.Once
	minioClient *minio.Client
	err      error
)

func init() {
	godotenv.Load()
	HOST := os.Getenv("MINIO_HOST")
	ACCESS_KEY := os.Getenv("MINIO_ACCESS_KEY");
	SECRET_KEY := os.Getenv("MINIO_SECRET_KEY");

	once.Do(func() {
		minioClient, err = minio.New(HOST, ACCESS_KEY, SECRET_KEY, false);
		if err != nil {
			log.Fatalf("Failed to create MinIO client: %v", err)
		}
	})
}

func GetMinioClient() *minio.Client {
	return minioClient
}

func bucketExists(bucketName string) error {
	exists, err := minioClient.BucketExists(bucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %v", err)
	}

	if !exists {
		err = minioClient.MakeBucket(bucketName, "")
	}

	if err != nil {
		return fmt.Errorf("failed to create bucket: %v", err)
	}
	return nil
}

func UploadFile(bucketName, objectName, filePath string) error {
	bucketExists(bucketName);
	_, err = minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}
	return nil
}

func UploadFileBuffer(bucketName string, objectName string, reader io.Reader, size int64) error {
	bucketExists(bucketName);
	_, err := minioClient.PutObject(bucketName, objectName, reader, size, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}
	return nil
}

func GetFileUrl(bucketName, objectName string) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := minioClient.PresignedGetObject(bucketName, objectName, time.Second*24*60*60, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate public url for the file: %v", err)
	}
	return presignedURL.String(), nil;
}