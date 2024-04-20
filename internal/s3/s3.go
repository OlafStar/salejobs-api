package s3

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/olafstar/salejobs-api/internal/env"
)

type S3Client struct {
	s3 *s3.S3
}

func InitS3Client() *S3Client {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("auto"),
		Endpoint: aws.String(env.GoEnv("S3_CLIENT_URL")),
		Credentials: credentials.NewStaticCredentials(env.GoEnv("R2_ACCESS_KEY_ID"), env.GoEnv("R2_SECRET_ACCESS_KEY"), ""),
	}))
	return &S3Client{
		s3: s3.New(sess),
	}
}

func (s3client *S3Client) GetPresignedURL(bucket, key string) (string, error) {
	req, _ := s3client.s3.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return req.Presign(15 * time.Minute)
}

func (s3client *S3Client) UploadData(data []byte, bucketName, key string) error {
	url, err := s3client.GetPresignedURL(bucketName, key)
	if err != nil {
		return fmt.Errorf("failed to get presigned URL: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload data, status: %d", resp.StatusCode)
	}

	fmt.Println("Data uploaded successfully")
	return nil
}

