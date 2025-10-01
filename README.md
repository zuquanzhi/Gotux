# Gotux - 图床管理系统

一个基于 Golang + Vue 的现代化图床管理系统，支持图片上传、管理、链接导出等功能。

## ✨ 功能特性

### 核心功能
- 🖼️ **图片上传**：支持拖拽上传、批量上传、多种图片格式
- 📁 **图片管理**：查看、编辑、删除、搜索图片
- 🔗 **链接导出**：支持 URL、Markdown、HTML、BBCode 等多种格式
- 📊 **数据统计**：图片数量、存储空间、访问量统计
- 🔍 **图片搜索**：按文件名、描述、标签搜索

### 用户系统
- 👤 **用户注册/登录**：支持邮箱注册
- 🔐 **JWT 认证**：安全的身份验证机制
- 👥 **用户管理**：修改个人信息、修改密码
- 🎯 **权限控制**：普通用户和管理员角色

### 管理功能
- 🛠️ **用户管理**：管理员可以查看所有用户、禁用/激活用户
- 📷 **图片管理**：管理员可以查看和管理所有图片
- 📈 **系统统计**：用户数、图片数、存储使用情况

### 技术特性
- ⚡ **高性能**：Golang 后端，Vue 3 前端
- 🗄️ **数据持久化**：SQLite 数据库
- 🔄 **图片去重**：基于 MD5 哈希的自动去重
- 📦 **文件组织**：按日期自动组织存储
- 🎨 **现代化 UI**：Element Plus 组件库
- 📱 **响应式设计**：支持各种屏幕尺寸

## 🚀 快速开始

### 系统要求

- Go 1.21+
- Node.js 18+
- npm 或 yarn

### 后端部署

```bash
cd backend

# 安装依赖
go mod download

# 复制配置文件
cp .env.example .env

# 编辑 .env 文件，修改配置（可选）

# 运行服务
go run main.go
```

后端服务将在 `http://localhost:8080` 启动

### 前端部署

```bash
cd frontend

# 安装依赖
npm install

# 开发模式
npm run dev

# 生产构建
npm run build
```

前端开发服务器将在 `http://localhost:5173` 启动

### 默认管理员账户

- 用户名: `admin`
- 密码: `admin123`

**⚠️ 请在首次登录后立即修改默认密码！**

## 📁 项目结构

```
Gotux/
├── backend/                 # 后端代码
│   ├── config/             # 配置
│   ├── controllers/        # 控制器
│   ├── middleware/         # 中间件
│   ├── models/             # 数据模型
│   ├── routes/             # 路由
│   ├── main.go            # 入口文件
│   └── go.mod             # Go 依赖
├── frontend/               # 前端代码
│   ├── src/
│   │   ├── api/           # API 接口
│   │   ├── components/    # 组件
│   │   ├── layout/        # 布局
│   │   ├── router/        # 路由
│   │   ├── stores/        # 状态管理
│   │   ├── utils/         # 工具函数
│   │   └── views/         # 页面
│   ├── package.json       # Node 依赖
│   └── vite.config.js     # Vite 配置
└── README.md              # 项目说明
```

## 🔧 配置说明

### 后端配置 (.env)

```env
SERVER_PORT=8080                                    # 服务端口
SERVER_MODE=release                                 # 运行模式
JWT_SECRET=your-secret-key-change-in-production    # JWT 密钥
```

### 前端配置 (vite.config.js)

代理配置已设置好，开发模式下自动代理到后端服务。

## 📸 功能截图

### 用户界面
- 登录/注册页面
- 仪表盘（数据统计）
- 图片上传页面
- 图片管理页面
- 个人中心

### 管理员界面
- 用户管理
- 系统图片管理
- 系统统计

## 🌐 API 文档

详细的 API 文档请查看 [backend/README.md](backend/README.md)

### 主要 API 端点

- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `POST /api/images/upload` - 上传图片
- `GET /api/images` - 获取图片列表
- `GET /api/images/:id/links` - 获取图片链接
- `DELETE /api/images/:id` - 删除图片
- `GET /api/user/profile` - 获取个人信息
- `GET /api/admin/stats` - 系统统计（管理员）

## 🚀 生产部署

### 使用 Docker

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### Nginx 配置

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    # 前端
    root /var/www/gotux/frontend/dist;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # 后端 API
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
    
    # 图片文件
    location /uploads {
        proxy_pass http://localhost:8080;
    }
}
```

## 🛠️ 开发指南

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
