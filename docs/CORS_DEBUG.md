# CORS 403 问题调试与解决

## 问题描述

在访问 `/api/random` 端点时,持续返回 403 Forbidden 错误:

```
[GIN] 2025/10/02 - 02:48:07 | 403 | 0s | ::1 | GET "/api/random"
```

关键信息:
- 请求来自 `::1` (IPv6 localhost)
- 从 `127.0.0.1` (IPv4) 访问正常返回 200
- 只有随机 API 端点受影响

## 根本原因

### CORS Origin 不匹配

浏览器在使用 IPv6 时,发送的 Origin 头是:
```
Origin: http://[::1]:5173
```

但原始 CORS 配置只允许:
```go
AllowOrigins: []string{"http://localhost:5173", "http://localhost:3000"}
```

`http://[::1]:5173` ≠ `http://localhost:5173`

因此 CORS 中间件拒绝了请求,返回 403。

## 尝试的解决方案

### 方案 1: 添加 IPv6 地址 (❌ 未完全解决)

```go
AllowOrigins: []string{
    "http://localhost:5173", 
    "http://127.0.0.1:5173", 
    "http://[::1]:5173"
}
```

**问题**: 需要列举所有可能的端口和协议组合。

### 方案 2: 使用 AllowOriginFunc (❌ 复杂且可能有坑)

```go
AllowOriginFunc: func(origin string) bool {
    if origin == "" {
        return true  // file:// 协议
    }
    return strings.HasPrefix(origin, "http://localhost:") || 
           strings.HasPrefix(origin, "http://127.0.0.1:") ||
           strings.HasPrefix(origin, "http://[::1]:")
}
```

**问题**: 
- 需要导入 `strings` 包
- 与 `AllowOrigins` 同时存在可能冲突
- 逻辑复杂,容易出错

### 方案 3: AllowAllOrigins (✅ 开发环境最佳)

```go
r.Use(cors.New(cors.Config{
    AllowAllOrigins:  true,
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length", "X-Image-UUID", "X-Image-ID"},
    AllowCredentials: false, // 必须设为 false
}))
```

**优点**:
- 简单直接
- 适合开发环境
- 支持所有 localhost 变体(IPv4/IPv6)
- 支持 file:// 协议(本地 HTML 文件)

## 最终解决方案

### 开发环境配置

```go
// backend/main.go
r.Use(cors.New(cors.Config{
    AllowAllOrigins:  true,
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length", "X-Image-UUID", "X-Image-ID"},
    AllowCredentials: false,
}))
```

### 生产环境配置

```go
// backend/main.go
r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{
        "https://yourdomain.com",
        "https://www.yourdomain.com",
    },
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length", "X-Image-UUID", "X-Image-ID"},
    AllowCredentials: true,
}))
```

## 重要注意事项

### AllowAllOrigins vs AllowCredentials

不能同时设置:
```go
AllowAllOrigins: true   // ❌
AllowCredentials: true  // ❌
```

原因: 这是 CORS 规范的安全限制。

正确用法:
- 开发环境: `AllowAllOrigins: true`, `AllowCredentials: false`
- 生产环境: 明确列出域名,`AllowCredentials: true`

### localhost 的多种形式

浏览器可能使用:
- `localhost` (DNS 名称)
- `127.0.0.1` (IPv4)
- `::1` (IPv6)
- `[::1]` (IPv6 在 URL 中的表示)

每种都被视为不同的 Origin!

## 测试验证

### 命令行测试

```bash
# IPv4
curl http://127.0.0.1:8080/api/random

# IPv6
curl http://[::1]:8080/api/random

# localhost (可能解析为 IPv4 或 IPv6)
curl http://localhost:8080/api/random
```

### 浏览器测试

1. 打开 `docs/random-demo.html`
2. 检查浏览器控制台网络选项卡
3. 查看请求的 Origin 头
4. 确认响应状态码为 200

### 检查 CORS 头

```bash
curl -H "Origin: http://localhost:5173" \
     -H "Access-Control-Request-Method: GET" \
     -H "Access-Control-Request-Headers: Content-Type" \
     -X OPTIONS \
     http://localhost:8080/api/random -v
```

应该看到:
```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
```

## 常见错误

### 错误 1: 403 但没有 CORS 错误

**症状**: 服务器返回 403,但浏览器控制台没有 CORS 错误提示。

**原因**: 403 是服务器返回的,不是 CORS 预检失败。

**检查**: 
- 查看服务器日志
- 确认路由配置正确
- 检查是否有其他中间件拦截

### 错误 2: OPTIONS 请求失败

**症状**: 浏览器发送 OPTIONS 预检请求失败。

**原因**: CORS 配置未包含 OPTIONS 方法。

**解决**: 
```go
AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
```

### 错误 3: 自定义头无法访问

**症状**: 响应头存在,但 JavaScript 无法读取。

**原因**: 未在 `ExposeHeaders` 中声明。

**解决**:
```go
ExposeHeaders: []string{"X-Image-UUID", "X-Image-ID"}
```

## 相关配置文件

- `backend/main.go` - CORS 配置位置
- `docs/DEPLOYMENT.md` - 生产环境 CORS 配置指南
- `docs/FIXES.md` - Bug 修复记录

## 总结

### 问题本质
CORS Origin 不匹配,浏览器使用 IPv6 但配置只允许特定域名。

### 解决方案
开发环境使用 `AllowAllOrigins: true`,生产环境明确列出允许的域名。

### 经验教训
1. localhost 不是唯一的,有多种表示形式
2. CORS 配置需要考虑 IPv4/IPv6
3. 开发和生产环境应使用不同配置
4. `AllowAllOrigins` 和 `AllowCredentials` 不能同时为 true

### 参考资料
- [MDN CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
- [gin-contrib/cors](https://github.com/gin-contrib/cors)
- [CORS Spec](https://fetch.spec.whatwg.org/#http-cors-protocol)
