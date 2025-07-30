package services

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"cloud.google.com/go/storage"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"net/url"
)

// StorageService defines the interface for file storage operations
type StorageService interface {
	// UploadFile uploads a file to the storage backend
	UploadFile(ctx context.Context, key string, reader io.Reader, contentType string) (string, error)
	
	// DeleteFile deletes a file from the storage backend
	DeleteFile(ctx context.Context, key string) error
	
	// GetFileURL returns a public URL for accessing the file
	GetFileURL(ctx context.Context, key string) (string, error)
	
	// GetSignedURL returns a signed URL for temporary access
	GetSignedURL(ctx context.Context, key string, expiration int64) (string, error)
}

// LocalStorageService implements StorageService for local file system
type LocalStorageService struct {
	baseURL string
}

// NewLocalStorageService creates a new local storage service
func NewLocalStorageService(baseURL string) *LocalStorageService {
	return &LocalStorageService{
		baseURL: baseURL,
	}
}

// UploadFile uploads a file to local storage
func (s *LocalStorageService) UploadFile(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
	// TODO: Implement local file upload
	// For now, return a placeholder URL
	return s.baseURL + "/files/" + key, nil
}

// DeleteFile deletes a file from local storage
func (s *LocalStorageService) DeleteFile(ctx context.Context, key string) error {
	// TODO: Implement local file deletion
	return nil
}

// GetFileURL returns the public URL for a local file
func (s *LocalStorageService) GetFileURL(ctx context.Context, key string) (string, error) {
	return s.baseURL + "/files/" + key, nil
}

// GetSignedURL returns a signed URL (same as public URL for local storage)
func (s *LocalStorageService) GetSignedURL(ctx context.Context, key string, expiration int64) (string, error) {
	return s.GetFileURL(ctx, key)
}

// S3StorageService implements StorageService for AWS S3
type S3StorageService struct {
	session    *session.Session
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	s3Client   *s3.S3
	bucket     string
	region     string
}

// NewS3StorageService creates a new S3 storage service
func NewS3StorageService(bucket, region string) (*S3StorageService, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return &S3StorageService{
		session:    sess,
		uploader:   s3manager.NewUploader(sess),
		downloader: s3manager.NewDownloader(sess),
		s3Client:   s3.New(sess),
		bucket:     bucket,
		region:     region,
	}, nil
}

// UploadFile uploads a file to S3
func (s *S3StorageService) UploadFile(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
	uploadInput := &s3manager.UploadInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
		ACL:         aws.String("public-read"), // Make files publicly accessible
	}

	result, err := s.uploader.UploadWithContext(ctx, uploadInput)
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	return result.Location, nil
}

// DeleteFile deletes a file from S3
func (s *S3StorageService) DeleteFile(ctx context.Context, key string) error {
	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	_, err := s.s3Client.DeleteObjectWithContext(ctx, deleteInput)
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}

	return nil
}

// GetFileURL returns the public URL for an S3 file
func (s *S3StorageService) GetFileURL(ctx context.Context, key string) (string, error) {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, key), nil
}

// GetSignedURL returns a signed URL for temporary access to S3 file
func (s *S3StorageService) GetSignedURL(ctx context.Context, key string, expiration int64) (string, error) {
	req, _ := s.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	urlStr, err := req.Presign(time.Duration(expiration) * time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return urlStr, nil
}

// GCSStorageService implements StorageService for Google Cloud Storage
type GCSStorageService struct {
	client     *storage.Client
	bucket     string
	projectID  string
}

// NewGCSStorageService creates a new GCS storage service
func NewGCSStorageService(bucket, projectID string) (*GCSStorageService, error) {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}

	return &GCSStorageService{
		client:    client,
		bucket:    bucket,
		projectID: projectID,
	}, nil
}

// UploadFile uploads a file to GCS
func (s *GCSStorageService) UploadFile(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
	bucket := s.client.Bucket(s.bucket)
	obj := bucket.Object(key)

	w := obj.NewWriter(ctx)
	w.ContentType = contentType
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}} // Make publicly readable

	if _, err := io.Copy(w, reader); err != nil {
		w.Close()
		return "", fmt.Errorf("failed to copy data to GCS: %w", err)
	}

	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close GCS writer: %w", err)
	}

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.bucket, key), nil
}

// DeleteFile deletes a file from GCS
func (s *GCSStorageService) DeleteFile(ctx context.Context, key string) error {
	bucket := s.client.Bucket(s.bucket)
	obj := bucket.Object(key)

	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete from GCS: %w", err)
	}

	return nil
}

// GetFileURL returns the public URL for a GCS file
func (s *GCSStorageService) GetFileURL(ctx context.Context, key string) (string, error) {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.bucket, key), nil
}

// GetSignedURL returns a signed URL for temporary access to GCS file
func (s *GCSStorageService) GetSignedURL(ctx context.Context, key string, expiration int64) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(time.Duration(expiration) * time.Second),
	}

	urlStr, err := storage.SignedURL(s.bucket, key, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return urlStr, nil
}

// AzureBlobStorageService implements StorageService for Azure Blob Storage
type AzureBlobStorageService struct {
	serviceURL   azblob.ServiceURL
	containerURL azblob.ContainerURL
	accountName  string
	accountKey   string
	container    string
}

// NewAzureBlobStorageService creates a new Azure Blob storage service
func NewAzureBlobStorageService(accountName, accountKey, container string) (*AzureBlobStorageService, error) {
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure credentials: %w", err)
	}

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	serviceURL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))
	service := azblob.NewServiceURL(*serviceURL, p)
	containerURL := service.NewContainerURL(container)

	return &AzureBlobStorageService{
		serviceURL:   service,
		containerURL: containerURL,
		accountName:  accountName,
		accountKey:   accountKey,
		container:    container,
	}, nil
}

// UploadFile uploads a file to Azure Blob Storage
func (s *AzureBlobStorageService) UploadFile(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
	blobURL := s.containerURL.NewBlobURL(key)
	blockBlobURL := blobURL.ToBlockBlobURL()

	_, err := azblob.UploadStreamToBlockBlob(ctx, reader, blockBlobURL, azblob.UploadStreamToBlockBlobOptions{
		BufferSize: 4 * 1024 * 1024, // 4MB buffer
		MaxBuffers: 16,
		BlobHTTPHeaders: azblob.BlobHTTPHeaders{
			ContentType: contentType,
		},
		AccessConditions: azblob.BlobAccessConditions{},
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to Azure Blob: %w", err)
	}

	return fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", s.accountName, s.container, key), nil
}

// DeleteFile deletes a file from Azure Blob Storage
func (s *AzureBlobStorageService) DeleteFile(ctx context.Context, key string) error {
	blobURL := s.containerURL.NewBlobURL(key)

	_, err := blobURL.Delete(ctx, azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{})
	if err != nil {
		return fmt.Errorf("failed to delete from Azure Blob: %w", err)
	}

	return nil
}

// GetFileURL returns the public URL for an Azure Blob file
func (s *AzureBlobStorageService) GetFileURL(ctx context.Context, key string) (string, error) {
	return fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", s.accountName, s.container, key), nil
}

// GetSignedURL returns a signed URL for temporary access to Azure Blob file
func (s *AzureBlobStorageService) GetSignedURL(ctx context.Context, key string, expiration int64) (string, error) {
	blobURL := s.containerURL.NewBlobURL(key)

	// Create credential from account key
	credential, err := azblob.NewSharedKeyCredential(s.accountName, s.accountKey)
	if err != nil {
		return "", fmt.Errorf("failed to create credential: %w", err)
	}

	// Create a SAS token
	sasQueryParams, err := azblob.BlobSASSignatureValues{
		Protocol:      azblob.SASProtocolHTTPS,
		ExpiryTime:    time.Now().Add(time.Duration(expiration) * time.Second),
		Permissions:   azblob.BlobSASPermissions{Read: true}.String(),
		ContainerName: s.container,
		BlobName:      key,
	}.NewSASQueryParameters(credential)
	if err != nil {
		return "", fmt.Errorf("failed to generate SAS token: %w", err)
	}

	return fmt.Sprintf("%s?%s", blobURL.String(), sasQueryParams.Encode()), nil
}