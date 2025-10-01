# Gotux Backend

基于 Golang + Gin + GORM 的图床后端服务

## 功能特性

- 用户注册、登录、JWT 认证
- 图片上传（支持批量上传）
- 图片管理（查看、编辑、删除、搜索）
- 链接导出（Markdown、HTML、BBCode 等格式）
- 图片去重（基于 MD5 哈希）
- 访问统计
- 管理员功能（用户管理、系统统计）
- 文件按日期组织存储

## 快速开始

### 安装依赖

```bash
go mod download
```

### 配置环境变量

复制 `.env.example` 到 `.env` 并修改配置：

```bash
cp .env.example .env
```

### 运行服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动

### 默认管理员账户

- 用户名: `admin`
- 密码: `admin123`

**请在生产环境中立即修改默认密码！**

## API 文档

### 认证相关

#### 注册
```
POST /api/auth/register
Content-Type: application/json

{
  "username": "user",
  "email": "user@example.com",
  "password": "password123"
}
```

#### 登录
```
POST /api/auth/login
Content-Type: application/json

{
  "username": "user",
  "password": "password123"
}
```

### 图片相关

#### 上传图片
```
POST /api/images/upload
Authorization: Bearer <token>
Content-Type: multipart/form-data

files: [图片文件]
```

#### 获取图片列表
```
GET /api/images?page=1&page_size=20&keyword=搜索关键词
Authorization: Bearer <token>
```

#### 获取图片详情
```
GET /api/images/:id
Authorization: Bearer <token>
```

#### 获取图片链接
```
GET /api/images/:id/links
Authorization: Bearer <token>
```

#### 更新图片信息
```
PUT /api/images/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "description": "图片描述",
  "tags": "标签1,标签2",
  "is_public": true
}
```

#### 删除图片
```
DELETE /api/images/:id
Authorization: Bearer <token>
```

#### 批量删除图片
```
POST /api/images/batch-delete
Authorization: Bearer <token>
Content-Type: application/json

{
  "image_ids": [1, 2, 3]
}
```

### 用户相关

#### 获取个人信息
```
GET /api/user/profile
Authorization: Bearer <token>
```

#### 更新个人信息
```
PUT /api/user/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "email": "newemail@example.com",
  "avatar": "avatar_url"
}
```

#### 修改密码
```
POST /api/user/change-password
Authorization: Bearer <token>
Content-Type: application/json

{
  "old_password": "oldpass",
  "new_password": "newpass"
}
```

#### 获取统计信息
```
GET /api/user/stats
Authorization: Bearer <token>
```

### 管理员相关

#### 获取所有用户
```
GET /api/admin/users?page=1&page_size=20
Authorization: Bearer <token>
```

#### 更新用户状态
```
PUT /api/admin/users/:id/status
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": "active" // 或 "disabled"
}
```

#### 获取所有图片
```
GET /api/admin/images?page=1&page_size=20
Authorization: Bearer <token>
```

#### 获取系统统计
```
GET /api/admin/stats
Authorization: Bearer <token>
```

## 项目结构

```
backend/
├── main.go              # 入口文件
├── config/              # 配置
│   └── config.go
├── models/              # 数据模型
│   ├── database.go
│   ├── user.go
│   └── image.go
├── controllers/         # 控制器
│   ├── auth.go
│   ├── image.go
│   └── admin.go
├── middleware/          # 中间件
│   └── auth.go
└── routes/              # 路由
    └── routes.go
```

## 部署

### 使用 Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env.example .env
RUN mkdir uploads
EXPOSE 8080
CMD ["./main"]
```

### 使用 systemd

创建 `/etc/systemd/system/gotux.service`:

```ini
[Unit]
Description=Gotux Image Hosting Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/gotux
ExecStart=/opt/gotux/main
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

启动服务:

```bash
sudo systemctl enable gotux
sudo systemctl start gotux
```

## 许可证

MIT
