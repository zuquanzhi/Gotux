# Gotux Frontend

基于 Vue 3 + Vite + Element Plus 的图床前端应用

## 功能特性

- 🎨 现代化的 UI 设计
- 📱 响应式布局
- 🖼️ 图片拖拽上传
- 🔗 多种链接格式导出
- 🔍 图片搜索
- 📊 数据统计
- 👥 用户管理
- 🔐 JWT 认证

## 技术栈

- Vue 3
- Vite
- Vue Router
- Pinia
- Element Plus
- Axios

## 快速开始

### 安装依赖

```bash
npm install
```

### 开发模式

```bash
npm run dev
```

访问 http://localhost:5173

### 生产构建

```bash
npm run build
```

### 预览生产构建

```bash
npm run preview
```

## 项目结构

```
frontend/
├── src/
│   ├── api/              # API 接口
│   ├── assets/           # 静态资源
│   ├── components/       # 组件
│   ├── layout/           # 布局
│   ├── router/           # 路由
│   ├── stores/           # 状态管理
│   ├── utils/            # 工具函数
│   ├── views/            # 页面
│   ├── App.vue          # 根组件
│   └── main.js          # 入口文件
├── index.html
├── vite.config.js
└── package.json
```

## 默认账号

- 用户名: admin
- 密码: admin123

**请在生产环境中立即修改默认密码！**

## 部署

### Nginx 配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    root /var/www/gotux/frontend/dist;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
    
    location /uploads {
        proxy_pass http://localhost:8080;
    }
}
```

## 许可证

MIT
