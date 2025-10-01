# API Documentation

Base URL: `http://localhost:8080/api`

## Authentication

All authenticated endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

## Endpoints

### Authentication

#### Register
```http
POST /api/register
Content-Type: application/json

{
  "username": "string",
  "email": "string",
  "password": "string"
}
```

Response:
```json
{
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "username": "string",
    "email": "string"
  }
}
```

#### Login
```http
POST /api/login
Content-Type: application/json

{
  "username": "string",
  "password": "string"
}
```

Response:
```json
{
  "token": "jwt_token",
  "user": {
    "id": 1,
    "username": "string",
    "email": "string",
    "is_admin": false
  }
}
```

### User

#### Get Profile
```http
GET /api/user/profile
Authorization: Bearer <token>
```

#### Update Profile
```http
PUT /api/user/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "email": "new@email.com"
}
```

#### Change Password
```http
POST /api/user/change-password
Authorization: Bearer <token>
Content-Type: application/json

{
  "old_password": "string",
  "new_password": "string"
}
```

#### Get User Statistics
```http
GET /api/user/stats
Authorization: Bearer <token>
```

Response:
```json
{
  "image_count": 100,
  "storage_used": 1048576,
  "total_views": 1000
}
```

#### Get User Settings
```http
GET /api/user/settings
Authorization: Bearer <token>
```

Response:
```json
{
  "custom_domain": "https://img.example.com",
  "default_link_format": "url",
  "enable_watermark": false,
  "watermark_text": "",
  "watermark_position": "bottom-right",
  "compress_image": true,
  "compress_quality": 80,
  "max_image_size": 10,
  "allowed_image_types": ["jpg", "png", "gif"],
  "enable_image_review": false,
  "storage_quota": 1073741824,
  "used_storage": 1048576
}
```

#### Update User Settings
```http
PUT /api/user/settings
Authorization: Bearer <token>
Content-Type: application/json

{
  "custom_domain": "https://img.example.com",
  "compress_quality": 85,
  ...
}
```

### Random Image API

#### Get Random Image Info (JSON)
```http
GET /api/random
GET /api/random?user_id=1
GET /api/random?tags=风景
GET /api/random?user_id=1&tags=风景
```

Returns random public image information in JSON format.

Response:
```json
{
  "id": 42,
  "uuid": "550e8400-e29b-41d4-a716-446655440000",
  "file_name": "sunset.jpg",
  "original_name": "beautiful_sunset.jpg",
  "width": 1920,
  "height": 1080,
  "file_size": 524288,
  "mime_type": "image/jpeg",
  "tags": "风景,日落",
  "created_at": "2025-01-15T10:30:00Z",
  "stats": {
    "view_count": 128
  }
}
```

#### Get Random Image File
```http
GET /api/random/image
GET /api/random/image?user_id=1
GET /api/random/image?tags=风景
```

Directly returns random image binary data. Suitable for use in HTML `<img>` tags or as background images.

Response Headers:
- `Content-Type`: Image MIME type
- `Cache-Control`: `public, max-age=3600`
- `X-Image-UUID`: Image UUID
- `X-Image-ID`: Image database ID

#### Redirect to Random Image
```http
GET /api/random/redirect
GET /api/random/redirect?user_id=1
GET /api/random/redirect?tags=风景
```

Returns 302 redirect to permanent image link (`/i/:uuid`).

For detailed usage examples, see [Random API Documentation](./RANDOM_API.md).

### Images

#### Upload Images
```http
POST /api/images/upload
Authorization: Bearer <token>
Content-Type: multipart/form-data

files: file1, file2, ...
```

Response:
```json
{
  "message": "Uploaded 2 images successfully",
  "images": [
    {
      "id": 1,
      "uuid": "550e8400-e29b-41d4-a716-446655440000",
      "file_name": "image.jpg",
      "original_name": "photo.jpg",
      "file_path": "2025/01/02/image.jpg",
      "file_size": 102400,
      "width": 1920,
      "height": 1080
    }
  ],
  "errors": []
}
```

#### List Images
```http
GET /api/images?page=1&page_size=20
Authorization: Bearer <token>
```

Response:
```json
{
  "images": [...],
  "total": 100,
  "page": 1,
  "page_size": 20
}
```

#### Get Image Details
```http
GET /api/images/:id
Authorization: Bearer <token>
```

#### Update Image
```http
PUT /api/images/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "description": "New description",
  "tags": "tag1,tag2",
  "is_public": true
}
```

#### Delete Image
```http
DELETE /api/images/:id
Authorization: Bearer <token>
```

#### Batch Delete Images
```http
POST /api/images/batch-delete
Authorization: Bearer <token>
Content-Type: application/json

{
  "ids": [1, 2, 3]
}
```

#### Get Image Links
```http
GET /api/images/:id/links
Authorization: Bearer <token>
```

Response:
```json
{
  "image": {...},
  "links": {
    "url": "https://img.example.com/i/550e8400-e29b-41d4-a716-446655440000",
    "direct_url": "https://img.example.com/uploads/2025/01/02/image.jpg",
    "html": "<img src=\"...\" alt=\"photo.jpg\" />",
    "markdown": "![photo.jpg](...)",
    "bbcode": "[img]...[/img]",
    "markdown_with_link": "[![photo.jpg](...)](...)"
  }
}
```

### Public Image Access

#### Get Image Info by UUID
```http
GET /api/i/:uuid
```

Response:
```json
{
  "image": {
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "file_name": "image.jpg",
    "width": 1920,
    "height": 1080,
    ...
  }
}
```

#### Serve Image File by UUID
```http
GET /i/:uuid
```

Returns the image file directly. View count is automatically incremented.

### Admin

All admin endpoints require admin role.

#### List All Users
```http
GET /api/admin/users?page=1&page_size=20
Authorization: Bearer <admin_token>
```

#### Update User Status
```http
PUT /api/admin/users/:id/status
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "is_active": false
}
```

#### List All Images
```http
GET /api/admin/images?page=1&page_size=20
Authorization: Bearer <admin_token>
```

#### Get System Statistics
```http
GET /api/admin/stats
Authorization: Bearer <admin_token>
```

Response:
```json
{
  "total_users": 100,
  "total_images": 10000,
  "total_storage": 10737418240,
  "total_views": 100000
}
```

## Error Responses

All endpoints may return error responses in the following format:

```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- `200 OK`: Success
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Rate Limiting

Currently not implemented. Consider adding rate limiting for production use.

## Image URL Formats

### UUID-based (Recommended - Secure)
```
GET /i/{uuid}
```
Example: `https://img.example.com/i/550e8400-e29b-41d4-a716-446655440000`

### Direct Path (Legacy - Less Secure)
```
GET /uploads/{year}/{month}/{day}/{filename}
```
Example: `https://img.example.com/uploads/2025/01/02/image.jpg`

The UUID-based format is recommended as it:
- Hides the actual file structure
- Prevents enumeration attacks
- Allows access control
- Automatically tracks view counts
