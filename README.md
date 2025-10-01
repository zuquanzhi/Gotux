# Gotux

A modern self-hosted image hosting solution built with Go and Vue.js.

## Features

- Drag-and-drop image upload with batch support
- **Secure UUID-based image links** to prevent enumeration attacks
- Multiple export formats: URL, Markdown, HTML, BBCode
- Custom domain support for image links
- Image compression with adjustable quality
- Watermark support with custom text and position
- User authentication and role-based access control
- Storage quota management
- View count tracking
- Admin dashboard for user and image management
- MD5-based deduplication
- Responsive UI design

## Quick Start

### Requirements

- Go 1.21 or higher
- Node.js 18 or higher

### Installation

Clone the repository:

```bash
git clone https://github.com/zuquanzhi/Gotux.git
cd Gotux
```

Start with Docker Compose (recommended):

```bash
docker-compose up -d
```

Or start manually:

```bash
# Backend
cd backend
go mod download
go run main.go

# Frontend (in a new terminal)
cd frontend
npm install
npm run dev
```

Access the application:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080

### Default Admin Account

```
Username: admin
Password: admin123
```

**Important:** Change the default password after first login.

## Configuration

### Backend Configuration

Create a `.env` file in the `backend` directory:

```env
SERVER_PORT=8080
SERVER_MODE=release
JWT_SECRET=your-secret-key-change-in-production
```

### Frontend Configuration

API endpoint is configured in `frontend/vite.config.js`. Update the proxy settings if needed.

## User Settings

Users can customize their experience through the profile settings:

- **Custom Domain**: Use your own domain for image links
- **Link Format**: Default format for copied links (URL/Markdown/HTML/BBCode)
- **Image Compression**: Automatic compression with quality control (1-100)
- **Watermark**: Add text watermark with customizable position
- **Upload Limits**: File size and format restrictions
- **Storage Quota**: Monitor storage usage

## Tech Stack

**Backend**
- Go 1.21
- Gin Web Framework
- GORM (SQLite)
- JWT Authentication

**Frontend**
- Vue 3
- Vite
- Element Plus
- Pinia
- Axios

## Project Structure

```
.
├── backend/
│   ├── config/         # Configuration
│   ├── controllers/    # Request handlers
│   ├── middleware/     # Middleware (auth, etc.)
│   ├── models/         # Data models
│   ├── routes/         # API routes
│   └── main.go         # Entry point
├── frontend/
│   └── src/
│       ├── api/        # API clients
│       ├── components/ # Vue components
│       ├── router/     # Vue Router
│       ├── stores/     # Pinia stores
│       ├── views/      # Page components
│       └── main.js     # Entry point
└── docker-compose.yml
```

## API Reference

### Authentication

```
POST /api/register      # Register new user
POST /api/login         # User login
```

### Images

```
GET    /api/images           # List images (paginated)
POST   /api/images/upload    # Upload images
GET    /api/images/:id       # Get image details
DELETE /api/images/:id       # Delete image
GET    /api/images/:id/links # Get image links in various formats
```

### User

```
GET  /api/user/profile     # Get user profile
PUT  /api/user/profile     # Update profile
PUT  /api/user/password    # Change password
GET  /api/user/stats       # Get user statistics
GET  /api/user/settings    # Get user settings
PUT  /api/user/settings    # Update user settings
```

### Admin (requires admin role)

```
GET    /api/admin/users        # List all users
PUT    /api/admin/users/:id    # Update user
DELETE /api/admin/users/:id    # Delete user
GET    /api/admin/images       # List all images
DELETE /api/admin/images/:id   # Delete any image
GET    /api/admin/stats        # System statistics
```

## Documentation

- [Deployment Guide](./docs/DEPLOYMENT.md) - Production deployment with Docker, Nginx, SSL
- [API Reference](./docs/API.md) - Complete API documentation
- [Changelog](./docs/CHANGELOG.md) - Version history

## Database Migration

If you're upgrading from a previous version, run the UUID migration:

```bash
cd backend/cmd/migrate_uuid
go run main.go
```

This adds UUID support to all existing images for secure access.

## Deployment

See [docs/DEPLOYMENT.md](./docs/DEPLOYMENT.md) for production deployment instructions including:

- Docker deployment
- Nginx/Caddy reverse proxy configuration
- SSL certificate setup with Let's Encrypt
- Custom domain configuration
- CDN integration

## Development

### Backend

```bash
cd backend
go run main.go
```

### Frontend

```bash
cd frontend
npm run dev
```

### Database Migration

```bash
cd backend/cmd/migrate
go run main.go
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see the LICENSE file for details

## Acknowledgments

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [GORM](https://gorm.io/) - ORM library
- [Vue.js](https://vuejs.org/) - Progressive JavaScript framework
- [Element Plus](https://element-plus.org/) - Vue 3 UI library
- [Vite](https://vitejs.dev/) - Build tool

## � 文档

- [功能更新日志](./CHANGELOG.md) - 版本更新历史
- [个人设置指南](./SETTINGS_GUIDE.md) - 详细的设置功能说明
- [域名绑定配置指南](./DOMAIN_BINDING_GUIDE.md) - 如何绑定自己的域名
- [设置快速参考](./SETTINGS_QUICKREF.md) - 快速查询设置选项
- [UI 极简化说明](./UI_MINIMALISM.md) - 界面设计理念
- [部署说明](./DEPLOY.md) - 生产环境部署指南

## �🛠️ 开发指南

### 后端开发

```bash
cd backend
go run main.go
```

### 前端开发

```bash
cd frontend
npm run dev
```

### 代码规范

- 后端遵循 Go 标准代码规范
- 前端使用 ESLint 进行代码检查
- 提交前请确保代码通过测试

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📝 许可证

MIT License

## 🙏 致谢

- [Gin](https://github.com/gin-gonic/gin) - Go Web 框架
- [GORM](https://gorm.io/) - Go ORM 库
- [Vue 3](https://vuejs.org/) - 渐进式 JavaScript 框架
- [Element Plus](https://element-plus.org/) - Vue 3 组件库
- [Vite](https://vitejs.dev/) - 下一代前端构建工具

## 📧 联系方式

如有问题或建议，请提交 Issue。

---

**⭐ 如果这个项目对你有帮助，请给个 Star！**
