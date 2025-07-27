# WhatsApp Integration with 360dialog

This document describes the comprehensive WhatsApp integration using 360dialog for phone-based authentication, user verification, and password management.

## Overview

The system integrates 360dialog's WhatsApp Business API to provide:

- **Enhanced Phone Verification**: 2-digit code verification with interactive buttons for web registrations
- **Direct WhatsApp Registration**: Auto-verified accounts created from WhatsApp messages
- **Password Reset via WhatsApp**: Complete password reset flow through WhatsApp
- **Interactive Button System**: Rich button-based interactions for better UX
- **Webhook Processing**: Comprehensive handling of button responses and text messages

## Environment Configuration

Add these environment variables to your `.env` file:

```bash
# 360dialog API configuration
WHATSAPP_API_KEY=your_360dialog_api_key_here
WHATSAPP_VERIFY_TOKEN=your_webhook_verify_token_here
```

## Enhanced Security Features

### 2-Digit Code Verification System

When users register on the web:

1. A random 2-digit code (10-99) is generated and displayed on the web
2. WhatsApp verification message shows 3 numbered buttons with different 2-digit codes
3. User must select the button matching the code shown on their registration page
4. Only the correct code selection verifies the account

This prevents unauthorized verification attempts and ensures the user has access to both the web session and WhatsApp.

## API Integration

### Core Services (pkg/services/whatsapp.go)

The `APIService` provides comprehensive WhatsApp functionality:

- `SendWhatsAppMessage()` - Send simple text messages
- `SendVerificationMessage()` - Send enhanced verification with 3 numbered buttons
- `SendWelcomeMessage()` - Send welcome message with action buttons
- `SendPasswordResetMessage()` - Send password reset with interactive buttons
- `generateThreeButtonCodes()` - Generate 3 unique 2-digit codes for verification
- `generateTwoDigitCode()` - Create random 2-digit verification codes

### Webhook System (pkg/handlers/whatsapp_webhook.go)

Comprehensive webhook handling includes:

- **Button Response Processing**: Handle all interactive button types
- **Text Message Processing**: Process registration requests and password changes
- **Enhanced Verification**: Multi-code verification with security checks
- **Password Reset Handling**: Complete password reset flow via WhatsApp
- **Direct Password Setting**: Allow users to set passwords via WhatsApp text

### Task Queue Integration (pkg/tasks/)

Background task processing for:

- `PhoneVerificationTask` - Send verification messages with 2-digit codes
- `PasswordResetTask` - Send password reset messages with interactive buttons

## User Flows

### 1. Enhanced Web Registration Flow

1. User registers on website with phone number
2. System generates random 2-digit code (e.g., "87")
3. Web displays: "Your WhatsApp verification code is: 87"
4. WhatsApp message sent with 3 buttons: "🔢 87", "🔢 23", "🔢 45"
5. User must select the "🔢 87" button to verify
6. System validates selected code matches displayed code
7. Account verified and verification code cleared

### 2. Password Reset via WhatsApp Flow

1. User enters phone number on "Forgot Password" page
2. System generates JWT reset token and sends WhatsApp message
3. WhatsApp shows: "🔐 Reset Password" and "❓ Help" buttons
4. User clicks "🔐 Reset Password" button
5. System validates token and phone number
6. User can either:
   - Use web form to set new password
   - Reply with "NEW PASSWORD: newpassword123" in WhatsApp
7. Password updated and confirmation sent

### 3. WhatsApp-Only Registration Flow

1. User messages the WhatsApp bot with registration intent
2. System auto-creates verified account
3. System sends welcome message with login instructions
4. User can access web platform immediately

### 4. Direct Password Setting via WhatsApp

1. User replies with format: "NEW PASSWORD: mypassword123"
2. System validates password requirements (minimum 8 characters)
3. Password hashed and stored securely
4. Confirmation message sent

## Interactive Button System

### Button Types and Handlers

The system handles these interactive button types:

#### Verification Buttons

- `verify_{token}_{code}` - Enhanced verification with code validation
- `verify_{token}` - Legacy verification (backwards compatible)
- `help_verification` - Verification help message

#### Password Reset Buttons

- `reset_password_{token}` - Process password reset with token validation
- `help_reset` - Password reset help message

#### General Actions

- `login_web` - Web login instructions
- `get_started` - Getting started guide

## Message Templates

### Enhanced Verification Message

```json
{
  "to": "+1234567890",
  "type": "interactive",
  "interactive": {
    "type": "button",
    "header": {"type": "text", "text": "📱 Account Verification"},
    "body": {"text": "Hi John! Welcome to Zero.\n\nTo complete your registration, please verify your account by selecting the 2-digit code shown on your registration page.\n\n🔐 Choose the correct verification code:"},
    "footer": {"text": "This verification will expire in 15 minutes."},
    "action": {
      "buttons": [
        {"type": "reply", "reply": {"id": "verify_token123_87", "title": "🔢 87"}},
        {"type": "reply", "reply": {"id": "verify_token123_23", "title": "🔢 23"}},
        {"type": "reply", "reply": {"id": "verify_token123_45", "title": "🔢 45"}}
      ]
    }
  }
}
```

### Password Reset Message

```json
{
  "to": "+1234567890", 
  "type": "interactive",
  "interactive": {
    "type": "button",
    "header": {"type": "text", "text": "🔐 Password Reset"},
    "body": {"text": "Hi John! We received a request to reset your password.\n\nClick the button below to proceed with your password reset.\n\n⚠️ If you didn't request this, please ignore this message."},
    "footer": {"text": "This reset link expires in 1 hour for security."},
    "action": {
      "buttons": [
        {"type": "reply", "reply": {"id": "reset_password_token123", "title": "🔐 Reset Password"}},
        {"type": "reply", "reply": {"id": "help_reset", "title": "❓ Need Help?"}}
      ]
    }
  }
}
```

## Database Schema Updates

### User Entity Enhancements

```go
// Added to ent/schema/user.go
field.String("verification_code").
    Optional().
    Nillable().
    Comment("2-digit verification code shown on web for WhatsApp verification")
```

The verification code is:

- Generated during web registration
- Displayed to user on web interface  
- Validated against WhatsApp button selection
- Cleared after successful verification

## Setup Instructions

### 1. Get 360dialog API Key

- Sign up for 360dialog account
- Obtain API key from dashboard
- Configure your WhatsApp Business number

### 2. Configure Webhook

- Set webhook URL: `https://yourdomain.com/api/whatsapp/webhook`
- Set verify token (matches `WHATSAPP_VERIFY_TOKEN`)
- Enable webhook for messages and button interactions

### 3. Database Migration

```bash
# Generate new Ent code after schema changes
go run -mod=mod entgo.io/ent/cmd/ent generate ./ent/schema

# Run any pending migrations
go run ./cmd/migrate
```

### 4. Test Integration

- **Web Registration**: Register → receive 3-button verification → select correct code
- **Password Reset**: Request reset → receive WhatsApp buttons → reset via web or WhatsApp
- **WhatsApp Registration**: Message bot → auto-create account → receive welcome
- **Direct Password**: Reply "NEW PASSWORD: test123" → password updated

## Development & Testing

### Fallback Mode

If no API key is set, messages are logged to console:

```text
📱 WhatsApp Verification to +1234567890: Hi John! Your verification code is: 87
```

### Error Handling

- Invalid verification codes show helpful error messages
- Expired tokens redirect to registration
- Phone number mismatches are caught and reported
- Password validation ensures security requirements

### Security Features

- **JWT Tokens**: Secure verification and password reset flows
- **Code Validation**: 2-digit codes prevent unauthorized verification
- **Phone Verification**: Ensures phone number ownership
- **Token Expiry**: Time-limited security tokens
- **Password Hashing**: bcrypt encryption for passwords

## Production Considerations

### Security Best Practices

- Rate limiting on webhook endpoints
- Input validation on all user messages
- Secure token generation and validation
- Phone number format validation (E.164)
- HTTPS-only webhook URLs

### Monitoring & Metrics

Key metrics to track:

- Verification completion rates by method (correct vs incorrect code selection)
- Password reset completion rates
- WhatsApp message delivery success rates
- Button interaction rates
- Average verification time

### Performance Optimization

- Background task processing for message sending
- Database indexing on phone numbers and verification codes
- Webhook response time monitoring
- Error rate tracking and alerting

## API Reference

### Webhook Endpoints

- `GET /api/whatsapp/webhook` - Webhook verification for 360dialog
- `POST /api/whatsapp/webhook` - Handle incoming messages and button responses

### Response Formats

All webhook responses follow 360dialog specifications with proper error handling and logging for debugging.

## Troubleshooting

### Common Issues

1. **Verification buttons not working**: Check webhook URL and verify token
2. **Messages not sending**: Verify API key and 360dialog account status  
3. **Wrong code selected**: User will see error message, can request new verification
4. **Password reset not working**: Check JWT token generation and validation
5. **Database errors**: Ensure schema is up to date with verification_code field

### Debug Logging

Enable debug logging to see:

- Webhook payloads
- Button interaction details
- Token validation results
- Message sending attempts
