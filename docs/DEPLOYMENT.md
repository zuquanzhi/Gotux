# Deployment Guide

This guide covers production deployment of Gotux.

## Docker Deployment

### Using Docker Compose (Recommended)

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Custom Docker Build

```bash
# Build backend
cd backend
docker build -t gotux-backend .

# Build frontend
cd ../frontend
docker build -t gotux-frontend .

# Run containers
docker run -d -p 8080:8080 -v $(pwd)/uploads:/app/uploads gotux-backend
docker run -d -p 80:80 gotux-frontend
```

## Reverse Proxy Configuration

### Nginx Configuration

Create `/etc/nginx/sites-available/gotux`:

```nginx
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;
    
    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;
    
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    
    client_max_body_size 50M;
    
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    location ~* \.(jpg|jpeg|png|gif|webp|bmp|svg)$ {
        proxy_pass http://127.0.0.1:8080;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
}
```

Enable the configuration:

```bash
sudo ln -s /etc/nginx/sites-available/gotux /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### Caddy Configuration (Simpler Alternative)

Create `/etc/caddy/Caddyfile`:

```caddy
your-domain.com {
    reverse_proxy localhost:8080
    
    @images {
        path *.jpg *.jpeg *.png *.gif *.webp *.bmp *.svg
    }
    header @images {
        Cache-Control "public, max-age=2592000"
    }
}
```

Start Caddy:

```bash
sudo systemctl enable caddy
sudo systemctl start caddy
```

## SSL Certificate Setup

### Using Let's Encrypt with Certbot

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx

# Obtain certificate
sudo certbot --nginx -d your-domain.com

# Test auto-renewal
sudo certbot renew --dry-run
```

### Using Caddy (Automatic)

Caddy automatically obtains and renews SSL certificates from Let's Encrypt.

## Custom Domain Configuration

### Step 1: DNS Configuration

Add an A record in your DNS provider:

```
Type: A
Name: @ (or subdomain like img)
Value: your-server-ip
TTL: 600
```

Verify DNS propagation:

```bash
nslookup your-domain.com
```

### Step 2: Configure in Gotux

1. Log in to Gotux
2. Go to Profile Settings
3. In "Custom Domain" field, enter: `https://your-domain.com`
4. Save settings

## Environment Variables

Production environment variables:

```env
SERVER_PORT=8080
SERVER_MODE=release
JWT_SECRET=change-this-to-a-secure-random-string
DB_PATH=./gotux.db
UPLOAD_PATH=./uploads
```

## Security Best Practices

1. **Change Default Credentials**: Update admin password immediately
2. **Use Strong JWT Secret**: Generate a random 32+ character secret
3. **Enable HTTPS**: Always use SSL/TLS in production
4. **Configure Firewall**: Only expose necessary ports (80, 443)
5. **Regular Backups**: Backup database and uploads directory
6. **Update Regularly**: Keep dependencies and system packages updated
7. **Monitor Logs**: Review access and error logs regularly

## CDN Integration

### Cloudflare Setup

1. Add your domain to Cloudflare
2. Update nameservers at your registrar
3. Enable "Proxied" for your domain
4. Configure caching rules for images

### Custom CDN

Configure your CDN to point to your origin server and update the custom domain in Gotux settings.

## Backup Strategy

### Database Backup

```bash
# Create backup
cp backend/gotux.db backend/gotux.db.backup

# Automated daily backup
echo "0 2 * * * cp /path/to/gotux.db /path/to/backups/gotux-$(date +\%Y\%m\%d).db" | crontab -
```

### Uploads Backup

```bash
# Sync to backup location
rsync -av backend/uploads/ /path/to/backup/uploads/

# Or use cloud storage
rclone sync backend/uploads/ remote:gotux-uploads/
```

## Monitoring

### Health Check Endpoint

```bash
curl http://localhost:8080/health
```

### Log Locations

- Nginx: `/var/log/nginx/`
- Application: Check stdout/stderr or configure file logging
- System: `/var/log/syslog`

## Troubleshooting

### Images Not Loading

1. Check upload directory permissions
2. Verify proxy configuration
3. Check firewall rules
4. Review error logs

### High Memory Usage

1. Check upload file size limits
2. Monitor concurrent uploads
3. Consider adding swap space
4. Use compression for large images

### Database Locked

1. Ensure only one backend instance is running
2. Check file permissions on database
3. Consider using PostgreSQL for high-traffic deployments

## Performance Optimization

1. Enable gzip compression in Nginx/Caddy
2. Configure browser caching for static assets
3. Use CDN for image delivery
4. Enable HTTP/2
5. Optimize image compression settings
