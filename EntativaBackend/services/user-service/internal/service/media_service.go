package service

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type MediaService struct {
	s3Client   *s3.S3
	bucketName string
	cdnURL     string
}

type S3Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	Region          string
	CDN             URL string
	UsePathStyle    bool // true for MinIO
}

// NewMediaService creates a new media upload service (S3/MinIO)
func NewMediaService(cfg *S3Config) (*MediaService, error) {
	if cfg == nil || cfg.AccessKeyID == "" {
		return nil, nil // Media service not configured
	}

	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Endpoint:         aws.String(cfg.Endpoint),
		Region:           aws.String(cfg.Region),
		S3ForcePathStyle: aws.Bool(cfg.UsePathStyle),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create S3 session: %w", err)
	}

	return &MediaService{
		s3Client:   s3.New(sess),
		bucketName: cfg.BucketName,
		cdnURL:     cfg.CDNURL,
	}, nil
}

// UploadProfilePicture uploads a profile picture and returns the URL
func (s *MediaService) UploadProfilePicture(userID uuid.UUID, file *multipart.FileHeader) (string, error) {
	if s == nil {
		return "", fmt.Errorf("media service not configured")
	}

	// Validate file type
	if !isValidImageType(file.Filename) {
		return "", fmt.Errorf("invalid file type. Only jpg, jpeg, png, gif allowed")
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return "", fmt.Errorf("file too large. Maximum size is 5MB")
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Read file content
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("profile-pictures/%s/%s%s", userID.String(), uuid.New().String(), ext)

	// Upload to S3/MinIO
	_, err = s.s3Client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(filename),
		Body:          bytes.NewReader(fileBytes),
		ContentType:   aws.String(file.Header.Get("Content-Type")),
		ContentLength: aws.Int64(file.Size),
		ACL:           aws.String("public-read"),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Generate URL
	if s.cdnURL != "" {
		return fmt.Sprintf("%s/%s", s.cdnURL, filename), nil
	}

	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucketName, filename), nil
}

// UploadCoverPhoto uploads a cover photo and returns the URL
func (s *MediaService) UploadCoverPhoto(userID uuid.UUID, file *multipart.FileHeader) (string, error) {
	if s == nil {
		return "", fmt.Errorf("media service not configured")
	}

	// Validate file type
	if !isValidImageType(file.Filename) {
		return "", fmt.Errorf("invalid file type. Only jpg, jpeg, png, gif allowed")
	}

	// Validate file size (max 10MB for cover photos)
	if file.Size > 10*1024*1024 {
		return "", fmt.Errorf("file too large. Maximum size is 10MB")
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Read file content
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("cover-photos/%s/%s%s", userID.String(), uuid.New().String(), ext)

	// Upload to S3/MinIO
	_, err = s.s3Client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(filename),
		Body:          bytes.NewReader(fileBytes),
		ContentType:   aws.String(file.Header.Get("Content-Type")),
		ContentLength: aws.Int64(file.Size),
		ACL:           aws.String("public-read"),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Generate URL
	if s.cdnURL != "" {
		return fmt.Sprintf("%s/%s", s.cdnURL, filename), nil
	}

	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucketName, filename), nil
}

// DeleteFile deletes a file from S3/MinIO
func (s *MediaService) DeleteFile(fileURL string) error {
	if s == nil {
		return nil
	}

	// Extract key from URL
	key := extractKeyFromURL(fileURL, s.bucketName, s.cdnURL)
	if key == "" {
		return fmt.Errorf("invalid file URL")
	}

	_, err := s.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})

	return err
}

// isValidImageType checks if the file has a valid image extension
func isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

// extractKeyFromURL extracts the S3 key from a URL
func extractKeyFromURL(url, bucket, cdnURL string) string {
	// Handle CDN URL
	if cdnURL != "" && strings.HasPrefix(url, cdnURL) {
		return strings.TrimPrefix(url, cdnURL+"/")
	}

	// Handle S3 URL
	s3Prefix := fmt.Sprintf("https://%s.s3.amazonaws.com/", bucket)
	if strings.HasPrefix(url, s3Prefix) {
		return strings.TrimPrefix(url, s3Prefix)
	}

	return ""
}
