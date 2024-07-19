package bucket

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"io"
)

func MustNewGoogleStorageClient(ctx context.Context, bucketName, credentialsFile string) *GoogleStorageClient {
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		panic(err.Error())
	}
	return &GoogleStorageClient{
		Bucket:        bucketName,
		StorageClient: storageClient,
	}
}

type GoogleStorageClient struct {
	Bucket        string
	StorageClient *storage.Client
}

func (c *GoogleStorageClient) SaveImage(ctx context.Context, fileName string, fileReader io.Reader) (string, error) {
	bucketName := c.Bucket
	writer := c.StorageClient.Bucket(bucketName).Object(fileName).NewWriter(ctx)
	if _, err := io.Copy(writer, fileReader); err != nil {
		return "", fmt.Errorf("io.Copy: %w", err)
	}
	if err := writer.Close(); err != nil {
		//Error say:  Error 403: The billing account for the owning project is disabled in state closed, accountDisabled
		return "", fmt.Errorf("Writer.Close: %w", err)
	}
	fmt.Fprintf(writer, "Blob %v uploaded.\n", fileName)
	return "https://storage.googleapis.com/" + "nddbao_bucket_test/" + fileName, nil
}
