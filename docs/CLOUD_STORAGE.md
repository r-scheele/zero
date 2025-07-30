# Cloud Storage Configuration

This application supports multiple cloud storage providers for file uploads. By default, it uses local storage, but you can configure it to use AWS S3, Google Cloud Storage (GCS), or Azure Blob Storage.

## Configuration

Cloud storage is configured via environment variables. Set the appropriate variables for your chosen provider:

### AWS S3

```bash
export AWS_S3_BUCKET="your-bucket-name"
export AWS_REGION="us-east-1"  # Optional, defaults to us-east-1
export AWS_ACCESS_KEY_ID="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"
```

### Google Cloud Storage (GCS)

```bash
export GCS_BUCKET="your-bucket-name"
export GCP_PROJECT_ID="your-project-id"
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account.json"
```

### Azure Blob Storage

```bash
export AZURE_STORAGE_ACCOUNT="your-storage-account"
export AZURE_STORAGE_KEY="your-storage-key"
export AZURE_CONTAINER="files"  # Optional, defaults to 'files'
```

## Features

- **Automatic Provider Selection**: The application automatically detects which cloud provider to use based on environment variables
- **Fallback to Local Storage**: If no cloud storage is configured, files are stored locally
- **Asynchronous Upload**: Files are uploaded to cloud storage in the background using task queues
- **File Size Limits**: Individual files are limited to 40MB, with a total limit of 400MB per note
- **Multiple File Support**: Up to 20 files can be uploaded per note
- **Signed URLs**: Support for generating temporary signed URLs for secure file access

## File Upload Process

1. User uploads files through the web interface
2. Files are temporarily saved locally
3. A background task is queued for each file
4. The task uploads the file to the configured cloud storage
5. The file URL is saved to the database
6. Temporary local files are cleaned up

## Supported File Types

- Images: JPG, JPEG, PNG, GIF
- Documents: PDF, DOC, DOCX, TXT
- Videos: MP4, AVI, MOV
- Audio files

## Security

- All cloud storage credentials should be kept secure
- Use IAM roles and least-privilege access policies
- Consider using signed URLs for temporary file access
- Files are organized by note ID in cloud storage for better organization

## Monitoring

File upload tasks can be monitored through the admin panel's task queue interface, which shows:
- Task status (pending, processing, completed, failed)
- Upload progress
- Error messages for failed uploads
- Retry attempts