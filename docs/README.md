# Gotux Documentation

## Getting Started

- [Main README](../README.md) - Overview and quick start guide

## Deployment

- [Deployment Guide](./DEPLOYMENT.md) - Production deployment with Docker, Nginx, SSL

## API Reference

- [API Documentation](./API.md) - Complete API reference with examples
- [Random Image API](./RANDOM_API.md) - Random image endpoints usage guide

## Changelog

- [Changelog](./CHANGELOG.md) - Version history and release notes

## Database

To add UUID support to existing images, run the migration:

```bash
cd backend/cmd/migrate_uuid
go run main.go
```

## Security Features

### UUID-based Image Links

All images are assigned a unique UUID for secure access:

```
https://your-domain.com/i/550e8400-e29b-41d4-a716-446655440000
```

Benefits:
- Prevents directory enumeration
- Hides actual file paths
- Enables access control per image
- Automatic view tracking

### Legacy Path Support

Direct path access is still supported for backwards compatibility:

```
https://your-domain.com/uploads/2025/01/02/image.jpg
```

However, UUID-based links are recommended for better security.

## Configuration

### User Settings

Users can configure:
- Custom domains
- Default link formats
- Image compression (1-100 quality)
- Watermarks with positioning
- Upload limits and allowed formats
- Storage quotas

### Admin Settings

Admins have access to:
- User management (activate/deactivate)
- System-wide image management
- Storage and usage statistics

## Support

For issues and questions, please open an issue on GitHub.
