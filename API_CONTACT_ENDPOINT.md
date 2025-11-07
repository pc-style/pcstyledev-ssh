# API Contact Endpoint Documentation

**Endpoint:** `/api/contact`  
**Base URL:** `https://pcstyle.dev/api/contact`  
**Version:** 1.0

---

## Overview

The `/api/contact` endpoint handles contact form submissions from both the web interface and SSH terminal interface. It validates input, applies rate limiting, and forwards messages to Discord via webhook.

**Features:**
- Input validation using Zod schema
- Rate limiting (5 requests per minute per IP)
- Discord webhook integration with rich embeds
- Source tracking (web vs SSH)
- Comprehensive error handling
- Health check endpoint

---

## Endpoints

### POST `/api/contact`

Submit a contact form message.

#### Request

**Headers:**
```
Content-Type: application/json
```

**Body Schema:**
```typescript
{
  message: string;        // Required, 1-2000 characters
  name?: string;          // Optional, max 100 characters
  email?: string;         // Optional, valid email format, max 100 characters
  discord?: string;       // Optional, max 100 characters
  phone?: string;         // Optional, max 50 characters
  facebook?: string;      // Optional, valid URL format, max 200 characters
  source?: 'web' | 'ssh'; // Optional, defaults to 'web'
}
```

**Example Request:**
```json
{
  "message": "Hello! I'd like to discuss a project.",
  "name": "John Doe",
  "email": "john@example.com",
  "discord": "@johndoe",
  "source": "web"
}
```

#### Response

**Success (200 OK):**
```json
{
  "success": true,
  "message": "Message sent successfully! Thanks for reaching out."
}
```

**Validation Error (400 Bad Request):**
```json
{
  "error": "Invalid input",
  "details": [
    {
      "code": "too_small",
      "minimum": 1,
      "type": "string",
      "inclusive": true,
      "exact": false,
      "message": "Message is required",
      "path": ["message"]
    }
  ]
}
```

**Rate Limit Exceeded (429 Too Many Requests):**
```json
{
  "error": "Too many requests. Please try again later."
}
```

**Server Error (500 Internal Server Error):**
```json
{
  "error": "Failed to send message. Please try again."
}
```

or

```json
{
  "error": "An unexpected error occurred. Please try again."
}
```

---

### GET `/api/contact`

Health check endpoint to verify the API is operational.

#### Response

**Success (200 OK):**
```json
{
  "status": "ok",
  "endpoint": "/api/contact",
  "methods": ["POST"]
}
```

---

## Rate Limiting

The endpoint implements in-memory rate limiting:
- **Window:** 60 seconds (1 minute)
- **Limit:** 5 requests per IP address per window
- **Identification:** Uses `x-forwarded-for` or `x-real-ip` headers (Vercel proxy headers), falls back to `'unknown'` if unavailable

When rate limit is exceeded, the endpoint returns `429 Too Many Requests` with a clear error message.

**Note:** Rate limiting uses in-memory storage. For high-traffic scenarios, consider migrating to Redis or another persistent store.

---

## Validation Rules

### Required Fields
- `message`: Must be a non-empty string between 1 and 2000 characters

### Optional Fields
- `name`: String, max 100 characters
- `email`: Valid email format, max 100 characters (can be empty string)
- `discord`: String, max 100 characters
- `phone`: String, max 50 characters
- `facebook`: Valid URL format, max 200 characters (can be empty string)
- `source`: Enum `'web'` or `'ssh'`, defaults to `'web'`

### Validation Errors

The endpoint uses Zod for validation. When validation fails, the response includes detailed error information in the `details` array, with each error containing:
- `code`: Error code (e.g., `too_small`, `invalid_type`, `invalid_string`)
- `message`: Human-readable error message
- `path`: Array indicating the field path where the error occurred

---

## Discord Integration

Messages are forwarded to Discord via webhook using rich embeds.

### Embed Format

**Title:** `New Contact Form Submission`

**Description:** The message content

**Color:**
- `0x00e5ff` (cyan) for `source: 'ssh'`
- `0xe6007e` (magenta) for `source: 'web'`

**Fields:**
- `Name`: If provided
- `Source`: Always included (WEB or SSH)
- `Contact Methods`: If any contact methods provided (Email, Discord, Phone, Facebook)

**Footer:** `pcstyle.dev contact form`

**Timestamp:** ISO 8601 timestamp of submission

### Environment Variable

The endpoint requires `DISCORD_WEBHOOK_URL` environment variable to be set. If not configured, the endpoint will return a 500 error.

**Setup:**
1. Create a Discord webhook in your server settings
2. Copy the webhook URL
3. Add to Vercel environment variables: `DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/...`

---

## Error Handling

The endpoint handles several error scenarios:

1. **Rate Limit Exceeded (429)**
   - Triggered when IP exceeds 5 requests per minute
   - Returns clear error message

2. **Validation Errors (400)**
   - Triggered by Zod schema validation failures
   - Returns detailed error information

3. **Discord Webhook Failure (500)**
   - Triggered when webhook URL is missing or request fails
   - Returns generic error message (doesn't expose internal details)

4. **Unexpected Errors (500)**
   - Catches any other errors
   - Logs to console (Vercel logs)
   - Returns generic error message

---

## Usage Examples

### cURL

**Basic submission:**
```bash
curl -X POST https://pcstyle.dev/api/contact \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Hello from cURL!",
    "name": "Test User",
    "email": "test@example.com",
    "source": "web"
  }'
```

**SSH source submission:**
```bash
curl -X POST https://pcstyle.dev/api/contact \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Hello from SSH terminal!",
    "name": "SSH User",
    "discord": "@sshuser",
    "source": "ssh"
  }'
```

**Health check:**
```bash
curl https://pcstyle.dev/api/contact
```

### JavaScript/TypeScript

```typescript
async function submitContact(data: {
  message: string;
  name?: string;
  email?: string;
  discord?: string;
  phone?: string;
  facebook?: string;
  source?: 'web' | 'ssh';
}) {
  const response = await fetch('https://pcstyle.dev/api/contact', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to submit contact form');
  }

  return await response.json();
}

// Usage
try {
  const result = await submitContact({
    message: 'Hello!',
    name: 'John Doe',
    email: 'john@example.com',
    source: 'web',
  });
  console.log(result.message); // "Message sent successfully! Thanks for reaching out."
} catch (error) {
  console.error('Error:', error.message);
}
```

### Python

```python
import requests
import json

def submit_contact(message, name=None, email=None, discord=None, phone=None, facebook=None, source='web'):
    url = 'https://pcstyle.dev/api/contact'
    data = {
        'message': message,
        'source': source
    }
    
    if name:
        data['name'] = name
    if email:
        data['email'] = email
    if discord:
        data['discord'] = discord
    if phone:
        data['phone'] = phone
    if facebook:
        data['facebook'] = facebook
    
    response = requests.post(url, json=data)
    
    if response.status_code == 200:
        return response.json()
    else:
        raise Exception(f"Error {response.status_code}: {response.json().get('error', 'Unknown error')}")

# Usage
try:
    result = submit_contact(
        message='Hello from Python!',
        name='Python User',
        email='python@example.com',
        source='web'
    )
    print(result['message'])
except Exception as e:
    print(f'Error: {e}')
```

---

## Implementation Details

### Rate Limiting Algorithm

The rate limiting uses a simple in-memory Map:
- Key: IP address (or 'unknown')
- Value: `{ count: number, resetTime: number }`

On each request:
1. Check if record exists and if current time > resetTime
2. If expired or missing, create new record with count=1
3. If exists and not expired, increment count
4. If count >= MAX_REQUESTS, return false (rate limited)

### IP Detection

The endpoint attempts to detect client IP in this order:
1. `x-forwarded-for` header (Vercel proxy)
2. `x-real-ip` header (Vercel proxy)
3. `'unknown'` fallback

**Note:** In production on Vercel, the proxy headers are automatically set. For local development, you may need to set these manually or the rate limiting will use `'unknown'` for all requests.

---

## Testing

### Local Development

1. Set up environment variable:
```bash
# .env.local
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/...
```

2. Start dev server:
```bash
npm run dev
```

3. Test endpoint:
```bash
curl -X POST http://localhost:3000/api/contact \
  -H "Content-Type: application/json" \
  -d '{"message":"Test message","source":"web"}'
```

### Production Testing

Use the health check endpoint first:
```bash
curl https://pcstyle.dev/api/contact
```

Then test a submission:
```bash
curl -X POST https://pcstyle.dev/api/contact \
  -H "Content-Type: application/json" \
  -d '{"message":"Production test","source":"web"}'
```

---

## Security Considerations

1. **Rate Limiting:** Prevents abuse and spam
2. **Input Validation:** Prevents injection attacks and malformed data
3. **Error Messages:** Generic error messages don't expose internal details
4. **Environment Variables:** Discord webhook URL stored securely in Vercel
5. **IP Detection:** Uses proxy headers for accurate IP detection

---

## Limitations

1. **In-Memory Rate Limiting:** Rate limit state is lost on server restart. For production at scale, consider Redis.
2. **Single Webhook:** All submissions go to one Discord channel. For multiple channels, modify the webhook logic.
3. **No Email Fallback:** If Discord webhook fails, there's no alternative delivery method.

---

## Changelog

### Version 1.0
- Initial implementation
- POST endpoint for contact submissions
- GET endpoint for health checks
- Rate limiting (5 req/min per IP)
- Discord webhook integration
- Source tracking (web/ssh)
- Comprehensive validation

---

## Support

For issues or questions:
- Check Vercel logs for server-side errors
- Verify `DISCORD_WEBHOOK_URL` is set correctly
- Test with health check endpoint first
- Review validation error details for input issues

