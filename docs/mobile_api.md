# Zero Mobile API Documentation

This document outlines the REST API endpoints available for mobile applications integrating with the Zero note-taking platform.

## Implementation Status

🟢 **Implemented**: Core authentication, note management, and social features  
🟡 **Partial**: Some endpoints may have limited functionality  
🔴 **Planned**: Advanced features not yet implemented  

## Base URL

```
https://your-domain.com/api/v1
```

## API Structure

The API is organized into the following main sections:
- **Health Check**: System status endpoint
- **Authentication**: User registration, login, and password management
- **Profile**: User profile management
- **Contact**: Contact form submission
- **Files**: File upload and management
- **Tasks**: Background task creation
- **Search**: Content search functionality
- **Admin**: Administrative operations
- **WhatsApp Webhooks**: WhatsApp integration endpoints

## Authentication

Most endpoints require authentication using JWT tokens. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

Tokens are valid for 24 hours and are returned upon successful login or registration.

## Response Format

All responses follow this general format:

**Success Response:**
```json
{
  "message": "Success message",
  "data": { ... }
}
```

**Error Response:**
```json
{
  "error": "Error message",
  "fields": { ... } // Only for validation errors
}
```

## Health Check Endpoint

### Health Check

**GET** `/health`

Check API availability and system status.

**Response (200 OK):**
```json
{
  "status": "healthy",
  "service": "external-api",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## WhatsApp Webhook Endpoints

### Verify Webhook

**GET** `/webhook`

Webhook verification endpoint for WhatsApp integration (360dialog).

**Query Parameters:**
- `hub.mode`: Should be "subscribe"
- `hub.verify_token`: Verification token
- `hub.challenge`: Challenge string to return

**Response (200 OK):**
Returns the challenge string if verification is successful.

### Handle Webhook

**POST** `/webhook`

Receive and process WhatsApp messages and status updates.

**Request Body:**
```json
{
  "messages": [
    {
      "from": "+1234567890",
      "id": "message_id",
      "timestamp": "1234567890",
      "type": "text",
      "text": {
        "body": "Hello"
      }
    }
  ],
  "statuses": [
    {
      "id": "message_id",
      "status": "delivered",
      "timestamp": "1234567890",
      "recipient_id": "+1234567890"
    }
  ]
}
```

**Response (200 OK):**
```
OK
```

## Authentication Endpoints

### Register 🟢

**POST** `/register`

Create a new user account.

**Request Body:**
```json
{
  "name": "John Doe",
  "phone_number": "+1234567890",
  "password": "securepassword123"
}
```

**Response (201 Created):**
```json
{
  "message": "Account created successfully",
  "user_id": 1,
  "verification_code": "42",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "verified": false
}
```

### Login 🟢

**POST** `/login`

Authenticate a user and receive a JWT token.

**Request Body:**
```json
{
  "phone_number": "+1234567890",
  "password": "securepassword123"
}
```

**Response (200 OK):**
```json
{
  "message": "Login successful",
  "user_id": 1,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "verified": true,
  "admin": false,
  "name": "John Doe"
}
```

### Logout 🟢

**POST** `/logout`

*Requires Authentication*

Logout the current user (client-side token invalidation).

**Response (200 OK):**
```json
{
  "message": "Logout successful"
}
```

### Forgot Password 🟢

**POST** `/forgot-password`

Request a password reset via WhatsApp.

**Request Body:**
```json
{
  "phone_number": "+1234567890"
}
```

**Response (200 OK):**
```json
{
  "message": "If your phone number is registered, you will receive a password reset message on WhatsApp"
}
```

### Reset Password 🟢

**POST** `/reset-password`

*Requires Authentication*

Reset user password (used with reset token from WhatsApp).

**Request Body:**
```json
{
  "password": "newpassword123"
}
```

**Response (200 OK):**
```json
{
  "message": "Password updated successfully"
}
```

### Resend Verification

**POST** `/resend-verification`

*Requires Authentication*

Resend phone verification message via WhatsApp.

**Response (200 OK):**
```json
{
  "message": "Verification message has been resent to your WhatsApp"
}
```

## Profile Endpoints

### Get Profile 🟢

**GET** `/profile`

*Requires Authentication*

Retrieve the current user's profile information.

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "John Doe",
  "phone_number": "+1234567890",
  "verified": true,
  "admin": false,
  "profile_picture": "https://example.com/profile.jpg",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Update Profile 🟢

**PUT** `/profile`

*Requires Authentication*

Update the current user's profile information.

**Request Body:**
```json
{
  "name": "Jane Doe"
}
```

**Response (200 OK):**
```json
{
  "message": "Profile updated successfully",
  "user": {
    "id": 1,
    "name": "Jane Doe",
    "phone_number": "+1234567890",
    "verified": true,
    "admin": false,
    "profile_picture": "https://example.com/profile.jpg",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

### Update Profile Picture

**POST** `/picture`

*Requires Authentication*

Upload a new profile picture.

**Request:** Multipart form data with `picture` field containing image file.

**Response (200 OK):**
```json
{
  "message": "Profile picture updated successfully",
  "profile_picture": "profile_1_1234567890.jpg"
}
```

### Change Password

**POST** `/change-password`

*Requires Authentication*

Change user password.

**Request Body:**
```json
{
  "current_password": "oldpassword123",
  "new_password": "newpassword123"
}
```

**Response (200 OK):**
```json
{
  "message": "Password updated successfully"
}
```

### Deactivate Account

**POST** `/deactivate`

*Requires Authentication*

Deactivate user account (permanent deletion).

**Request Body:**
```json
{
  "password": "userpassword123",
  "confirm": true
}
```

**Response (200 OK):**
```json
{
  "message": "Account deactivated successfully"
}
```

## Note Management Endpoints

### Create Note 🟢

**POST** `/notes`

*Requires Authentication*

Create a new note with title, description, content, and resources.

**Request Body:**
```json
{
  "title": "My Study Notes",
  "description": "Notes about advanced mathematics",
  "content": "Detailed content of the note...",
  "visibility": "private",
  "permission_level": "read_only",
  "resources": [
    {
      "type": "url",
      "url": "https://example.com/resource",
      "title": "External Resource"
    },
    {
      "type": "file",
      "filename": "document.pdf",
      "size": 1024000,
      "url": "https://storage.example.com/files/document.pdf"
    }
  ]
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "note": {
      "id": 456,
      "title": "My Study Notes",
      "description": "Notes about advanced mathematics",
      "content": "Detailed content of the note...",
      "visibility": "private",
      "permission_level": "read_only",
      "share_token": "abc123def456",
      "ai_processing": false,
      "resources": [
        {
          "type": "url",
          "url": "https://example.com/resource",
          "title": "External Resource"
        }
      ],
      "author": {
        "id": 123,
        "name": "John Doe"
      },
      "likes_count": 0,
      "reposts_count": 0,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  },
  "message": "Note created successfully",
  "error": null
}
```

## Contact Endpoints

### Submit Contact Form

**POST** `/contact`

*Requires Authentication*

Submit a contact form message.

**Request Body:**
```json
{
  "subject": "Bug Report",
  "message": "I found a bug in the application..."
}
```

**Response (200 OK):**
```json
{
  "message": "Contact form submitted successfully"
}
```

## File Endpoints

### List Files

**GET** `/files`

*Requires Authentication*

Retrieve a list of user's uploaded files.

**Response (200 OK):**
```json
{
  "files": [
    {
      "id": 1,
      "filename": "document.pdf",
      "original_name": "my-document.pdf",
      "size": 1024000,
      "content_type": "application/pdf",
      "uploaded_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### Upload File

**POST** `/upload`

*Requires Authentication*

Upload a file to the server.

**Request Body:** `multipart/form-data`
- `file`: The file to upload

**Response (200 OK):**
```json
{
  "message": "File uploaded successfully",
  "file": {
    "id": 1,
    "filename": "document_1_1234567890.pdf",
    "original_name": "my-document.pdf",
    "size": 1024000,
    "content_type": "application/pdf",
    "uploaded_at": "2024-01-01T00:00:00Z"
  }
}
```

## Task Endpoints

### Create Task

**POST** `/tasks`

*Requires Authentication*

Create a background task.

**Request Body:**
```json
{
  "type": "data_processing",
  "parameters": {
    "input_file": "document.pdf",
    "options": ["extract_text", "analyze"]
  }
}
```

**Response (200 OK):**
```json
{
  "message": "Task created successfully",
  "task_id": "task_123456"
}
```

## Search Endpoints

### Search

**GET** `/search`

*Requires Authentication*

Search for content across the platform.

**Query Parameters:**
- `q`: Search query (required)
- `type`: Content type filter (optional)
- `limit`: Number of results (optional, default: 10)

**Response (200 OK):**
```json
{
  "results": [
    {
      "id": 1,
      "type": "file",
      "title": "Document Title",
      "snippet": "...relevant content snippet...",
      "score": 0.95
    }
  ],
  "total": 1,
  "query": "search term"
}
```

## Admin Endpoints

*All admin endpoints require authentication and admin privileges.*

### Admin Overview

**GET** `/overview`

Get admin dashboard statistics.

**Response (200 OK):**
```json
{
  "stats": {
    "total_users": 150,
    "verified_users": 120,
    "admin_users": 3
  }
}
```

### List Users

**GET** `/users?page=1&limit=25`

Get paginated list of users.

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 25, max: 100)

**Response (200 OK):**
```json
{
  "users": [
    {
      "id": 1,
      "name": "John Doe",
      "phone_number": "+1234567890",
      "verified": true,
      "admin": false,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 25,
    "total": 150,
    "pages": 6
  }
}
```

### Get User

**GET** `/users/{id}`

Get specific user details.

**Response (200 OK):**
```json
{
  "user": {
    "id": 1,
    "name": "John Doe",
    "phone_number": "+1234567890",
    "email": "john@example.com",
    "verified": true,
    "admin": false,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### Verify User

**POST** `/users/{id}/verify`

Manually verify a user account.

**Response (200 OK):**
```json
{
  "message": "User verified successfully"
}
```

## Error Codes

- **400 Bad Request**: Invalid request format or validation errors
- **401 Unauthorized**: Authentication required or invalid token
- **403 Forbidden**: Insufficient privileges (admin required)
- **404 Not Found**: Resource not found
- **409 Conflict**: Resource already exists (e.g., duplicate phone number)
- **500 Internal Server Error**: Server error

## Rate Limiting

API endpoints may be rate limited. Check response headers for rate limit information:

- `X-RateLimit-Limit`: Maximum requests per time window
- `X-RateLimit-Remaining`: Remaining requests in current window
- `X-RateLimit-Reset`: Time when rate limit resets

## File Upload Limits

- Profile pictures: Maximum 5MB, image files only
- General file uploads: Check with server configuration

## Security Notes

1. Always use HTTPS in production
2. Store JWT tokens securely on the client side
3. Implement proper token refresh mechanisms
4. Validate all user inputs on the client side before sending
5. Handle sensitive operations (password changes, account deletion) with extra confirmation

## Example Usage

### JavaScript/React Native Example

```javascript
// Register a new user
const registerResponse = await fetch('http://localhost:8000/api/v1/register', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    name: 'John Doe',
    phone_number: '+1234567890',
    password: 'securepassword123'
  })
});

// Login
const loginResponse = await fetch('http://localhost:8000/api/v1/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    phone_number: '+1234567890',
    password: 'securepassword123'
  })
});

const { token } = await loginResponse.json();

// Get profile (authenticated request)
const profileResponse = await fetch('http://localhost:8000/api/v1/profile', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});
```

### Swift/iOS Example

```swift
// Login
func login(phoneNumber: String, password: String) async throws -> LoginResponse {
    let url = URL(string: "/api/v1/login")!
    var request = URLRequest(url: url)
    request.httpMethod = "POST"
    request.setValue("application/json", forHTTPHeaderField: "Content-Type")
    
    let body = [
        "phone_number": phoneNumber,
        "password": password
    ]
    request.httpBody = try JSONSerialization.data(withJSONObject: body)
    
    let (data, response) = try await URLSession.shared.data(for: request)
    
    if let httpResponse = response as? HTTPURLResponse, httpResponse.statusCode == 200 {
        return try JSONDecoder().decode(LoginResponse.self, from: data)
    } else {
        throw APIError.loginFailed
    }
}

// Authenticated request
func getProfile() async throws -> UserProfile {
    guard let token = KeychainHelper.getToken() else {
        throw APIError.noToken
    }
    
    let url = URL(string: "/api/v1/profile")!
    var request = URLRequest(url: url)
    request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
    
    let (data, _) = try await URLSession.shared.data(for: request)
    return try JSONDecoder().decode(UserProfile.self, from: data)
}
```

This API provides complete functionality for web and mobile applications with a unified endpoint structure.