# Gotux 部署指南

本文档提供了 Gotux 图床管理系统在不同环境下的部署方案。

## 目录

- [开发环境部署](#开发环境部署)
- [生产环境部署](#生产环境部署)
- [Docker 部署](#docker-部署)
- [云服务器部署](#云服务器部署)
- [安全配置](#安全配置)
- [性能优化](#性能优化)

## 开发环境部署

### Windows

```powershell
# 1. 克隆项目
git clone https://github.com/zuquanzhi/Gotux.git
cd Gotux

# 2. 使用启动脚本（推荐）
.\start.ps1

# 或手动启动

# 启动后端
cd backend
go mod download
go run main.go

# 新开一个终端，启动前端
cd frontend
npm install
npm run dev
```

### Linux/Mac

```bash
# 1. 克隆项目
git clone https://github.com/zuquanzhi/Gotux.git
cd Gotux

# 2. 使用启动脚本（推荐）
chmod +x start.sh
./start.sh

# 或手动启动

# 启动后端
cd backend
go mod download
go run main.go &

# 启动前端
cd frontend
npm install
npm run dev
```

## 生产环境部署

### 方案一：传统部署

#### 1. 编译后端

```bash
cd backend
go build -o gotux main.go
```

#### 2. 构建前端

```bash
cd frontend
npm install
npm run build
```

#### 3. 配置 Nginx

创建 `/etc/nginx/sites-available/gotux`:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    root /var/www/gotux/frontend/dist;
    index index.html;

    # 前端路由
    location / {
        try_files $uri $uri/ /index.html;
    }

    # 后端 API 代理
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket 支持（如需要）
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # 图片文件代理
    location /uploads {
        proxy_pass http://localhost:8080;
        
        # 缓存设置
        expires 7d;
        add_header Cache-Control "public, immutable";
    }

    # Gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript 
               application/x-javascript application/xml+rss 
               application/json application/javascript;

    # 安全头
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
}
```

启用配置：

```bash
ln -s /etc/nginx/sites-available/gotux /etc/nginx/sites-enabled/
nginx -t
systemctl restart nginx
```

#### 4. 配置 systemd 服务

创建 `/etc/systemd/system/gotux.service`:

```ini
[Unit]
Description=Gotux Image Hosting Service
After=network.target

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/var/www/gotux/backend
ExecStart=/var/www/gotux/backend/gotux
Restart=on-failure
RestartSec=5s

# 环境变量
Environment="SERVER_PORT=8080"
Environment="SERVER_MODE=release"
Environment="JWT_SECRET=your-production-secret-key-here"

# 日志
StandardOutput=append:/var/log/gotux/access.log
StandardError=append:/var/log/gotux/error.log

[Install]
WantedBy=multi-user.target
```

创建日志目录：

```bash
mkdir -p /var/log/gotux
chown www-data:www-data /var/log/gotux
```

启动服务：

```bash
systemctl enable gotux
systemctl start gotux
systemctl status gotux
```

### 方案二：使用 SSL (HTTPS)

#### 1. 安装 Certbot

```bash
# Ubuntu/Debian
apt install certbot python3-certbot-nginx

# CentOS
yum install certbot python3-certbot-nginx
```

#### 2. 获取 SSL 证书

```bash
certbot --nginx -d your-domain.com
```

#### 3. 更新 Nginx 配置

Certbot 会自动修改配置，或手动添加：

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;
    
    # SSL 配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # ... 其他配置同上
}

server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}
```

## Docker 部署

### 1. 使用 Docker Compose（推荐）

```bash
# 克隆项目
git clone https://github.com/zuquanzhi/Gotux.git
cd Gotux

# 编辑环境变量
nano .env

# 构建并启动
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 2. 手动构建

```bash
# 构建后端镜像
cd backend
docker build -t gotux-backend .

# 构建前端镜像
cd ../frontend
docker build -t gotux-frontend .

# 运行容器
docker run -d -p 8080:8080 -v $(pwd)/uploads:/app/uploads gotux-backend
docker run -d -p 80:80 gotux-frontend
```

## 云服务器部署

### 阿里云/腾讯云

#### 1. 购买服务器

- 配置建议：2核4G，带宽1-5M
- 系统：Ubuntu 20.04/22.04 或 CentOS 7/8

#### 2. 安全组配置

开放端口：
- 22 (SSH)
- 80 (HTTP)
- 443 (HTTPS)

#### 3. 安装依赖

```bash
# 更新系统
apt update && apt upgrade -y

# 安装 Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装 Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
apt install -y nodejs

# 安装 Nginx
apt install -y nginx
```

#### 4. 部署应用

按照[生产环境部署](#生产环境部署)步骤进行。

### AWS EC2

类似上述步骤，注意：
- 使用 Amazon Linux 2 或 Ubuntu
- 配置 IAM 角色和安全组
- 可选：使用 S3 存储图片

### 宝塔面板部署

1. 安装宝塔面板
2. 在面板中安装 Nginx、PM2
3. 创建站点，上传代码
4. 使用 PM2 管理后端进程
5. 配置 Nginx 反向代理

## 安全配置

### 1. 修改默认密码

首次部署后立即修改默认管理员密码。

### 2. 配置防火墙

```bash
# UFW (Ubuntu)
ufw allow 22
ufw allow 80
ufw allow 443
ufw enable

# firewalld (CentOS)
firewall-cmd --permanent --add-service=ssh
firewall-cmd --permanent --add-service=http
firewall-cmd --permanent --add-service=https
firewall-cmd --reload
```

### 3. 配置 JWT 密钥

编辑 `.env` 文件，设置强密码：

```env
JWT_SECRET=use-a-very-strong-random-string-here
```

生成随机密钥：

```bash
openssl rand -hex 32
```

### 4. 限制文件上传

在 `backend/config/config.go` 中调整：

```go
MaxSize: 10 * 1024 * 1024, // 10MB
AllowedTypes: []string{"image/jpeg", "image/png", "image/gif", "image/webp"},
```

### 5. 配置 CORS

根据实际域名修改 `backend/main.go` 中的 CORS 配置：

```go
AllowOrigins: []string{"https://your-domain.com"},
```

## 性能优化

### 1. 数据库优化

考虑使用 MySQL/PostgreSQL 替代 SQLite：

```go
// models/database.go
import "gorm.io/driver/mysql"

DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
```

### 2. 图片存储优化

使用 CDN 或对象存储：
- 阿里云 OSS
- 腾讯云 COS
- AWS S3
- 七牛云

### 3. 缓存配置

在 Nginx 中配置缓存：

```nginx
proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=image_cache:10m;

location /uploads {
    proxy_cache image_cache;
    proxy_cache_valid 200 7d;
    proxy_pass http://localhost:8080;
}
```

### 4. 启用 HTTP/2

```nginx
listen 443 ssl http2;
```

### 5. 图片压缩

考虑集成图片压缩服务，如 TinyPNG API。

## 监控和维护

### 1. 日志管理

```bash
# 查看应用日志
journalctl -u gotux -f

# 查看 Nginx 日志
tail -f /var/log/nginx/access.log
tail -f /var/log/nginx/error.log
```

### 2. 定期备份

备份脚本示例：

```bash
#!/bin/bash
DATE=$(date +%Y%m%d)
BACKUP_DIR="/backup/gotux"

# 备份数据库
cp /var/www/gotux/backend/gotux.db $BACKUP_DIR/gotux_$DATE.db

# 备份图片
tar -czf $BACKUP_DIR/uploads_$DATE.tar.gz /var/www/gotux/backend/uploads

# 删除 30 天前的备份
find $BACKUP_DIR -name "*.db" -mtime +30 -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +30 -delete
```

设置定时任务：

```bash
crontab -e
# 每天凌晨 2 点备份
0 2 * * * /path/to/backup.sh
```

### 3. 监控工具

推荐使用：
- Prometheus + Grafana
- 宝塔面板监控
- 云服务商自带监控

## 故障排查

### 后端无法启动

```bash
# 检查端口占用
lsof -i :8080

# 检查日志
journalctl -u gotux -n 50
```

### 前端无法访问

```bash
# 检查 Nginx 配置
nginx -t

# 重启 Nginx
systemctl restart nginx
```

### 图片上传失败

1. 检查上传目录权限
2. 检查磁盘空间
3. 查看应用日志

## 更新升级

```bash
# 备份数据
./backup.sh

# 拉取最新代码
git pull

# 后端
cd backend
go build -o gotux main.go
systemctl restart gotux

# 前端
cd frontend
npm install
npm run build

# 重启 Nginx
systemctl reload nginx
```

## 支持

如有问题，请查看：
- [GitHub Issues](https://github.com/zuquanzhi/Gotux/issues)
- [项目文档](https://github.com/zuquanzhi/Gotux)
