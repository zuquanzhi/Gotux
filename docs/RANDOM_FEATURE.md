# 随机图片 API 功能说明

## 新增功能概述

新增了**随机图片 API**,允许通过简单的 HTTP 请求获取图库中的随机图片,无需增加图库(Gallery)层级,利用现有的标签(tags)和用户(user_id)系统即可实现灵活的分类和筛选。

## 为什么不需要图库层级?

### 当前方案的优势:

1. **简单高效**: 利用现有的 `tags` 字段实现分类
2. **足够灵活**: 通过 `user_id` 和 `tags` 组合可以实现复杂筛选
3. **即刻可用**: 无需额外的数据库迁移和前端开发
4. **向后兼容**: 不影响现有功能

### 如果需要图库,可以这样扩展:

将来如果确实需要图库概念,可以通过以下方式实现:

1. **使用特殊标签**: 例如 `gallery:风景` 作为图库标识
2. **新增 Gallery 表**: 添加图库管理功能
3. **多对多关系**: 一张图片可以属于多个图库

## 新增 API 端点

### 1. `/api/random` - 获取随机图片信息

**用途**: 返回 JSON 格式的图片详细信息

**支持参数**:
- `user_id` - 筛选指定用户的图片
- `tags` - 筛选带特定标签的图片

**示例**:
```bash
GET /api/random
GET /api/random?user_id=1
GET /api/random?tags=风景
GET /api/random?user_id=1&tags=风景,自然
```

### 2. `/api/random/image` - 直接返回图片文件

**用途**: 直接返回图片二进制数据,适合在 HTML 中使用

**应用场景**:
- 网站随机背景
- 占位图服务
- 图片轮播

**示例**:
```html
<img src="http://yourdomain.com/api/random/image?tags=风景">
<div style="background-image: url('http://yourdomain.com/api/random/image')"></div>
```

### 3. `/api/random/redirect` - 重定向到图片

**用途**: 返回 302 重定向到图片的永久链接

**应用场景**:
- SEO 优化
- 社交媒体分享
- 需要永久链接的场合

## 技术实现

### Controller 层

在 `backend/controllers/image.go` 中新增了三个函数:

```go
func GetRandomImage(c *gin.Context)      // 返回 JSON
func ServeRandomImage(c *gin.Context)    // 返回图片文件
func RedirectRandomImage(c *gin.Context) // 重定向到图片
```

### 核心逻辑

```go
// 使用 SQLite 的 RANDOM() 函数
query := models.DB.Where("is_public = ?", true)

// 支持按用户筛选
if userID := c.Query("user_id"); userID != "" {
    query = query.Where("user_id = ?", userID)
}

// 支持按标签筛选
if tags := c.Query("tags"); tags != "" {
    query = query.Where("tags LIKE ?", "%"+tags+"%")
}

// 随机获取
var image models.Image
query.Order("RANDOM()").First(&image)
```

### 路由配置

在 `backend/routes/routes.go` 中添加:

```go
// 随机图片 API
api.GET("/random", controllers.GetRandomImage)
api.GET("/random/image", controllers.ServeRandomImage)
api.GET("/random/redirect", controllers.RedirectRandomImage)
```

## 使用示例

### 1. 网站随机背景

```html
<body style="background-image: url('http://localhost:8080/api/random/image?tags=壁纸')">
```

### 2. 动态图片展示

```javascript
// 每 10 秒刷新一次
setInterval(() => {
    document.getElementById('random-img').src = 
        `/api/random/image?t=${Date.now()}`;
}, 10000);
```

### 3. API 集成

```javascript
// 获取随机图片信息
fetch('/api/random?tags=风景')
    .then(res => res.json())
    .then(data => {
        console.log('随机图片:', data);
        const imageUrl = `/i/${data.uuid}`;
    });
```

### 4. Markdown 中使用

```markdown
![随机图片](http://localhost:8080/api/random/image)
```

## 性能考虑

### SQLite RANDOM() 性能

- 适合中小型数据集(< 10 万张图片)
- 如果数据量很大,可以考虑:
  - 使用缓存策略
  - 预生成随机序列
  - 使用 ID 范围随机

### 缓存策略

响应头设置了缓存:
```go
c.Header("Cache-Control", "public, max-age=3600")
```

## 安全性

### 隐私保护

- 只返回 `is_public = true` 的公开图片
- 私有图片不会出现在随机结果中
- 支持用户级别的权限控制

### 访问统计

- `ServeRandomImage` 会增加访问计数
- `RedirectRandomImage` 会增加访问计数
- `GetRandomImage` (JSON) 不增加访问计数

## 演示页面

创建了一个完整的演示页面: `docs/random-demo.html`

**功能**:
- 可视化筛选条件设置
- 实时加载随机图片
- 展示图片详细信息
- 包含完整的使用示例代码

**使用方法**:
1. 启动后端服务
2. 在浏览器中打开 `docs/random-demo.html`
3. 设置筛选条件,点击"加载图片"

## 文档更新

### 新增文档:
- `docs/RANDOM_API.md` - 详细的随机 API 使用指南

### 更新文档:
- `docs/API.md` - 添加随机 API 端点说明
- `docs/README.md` - 添加随机 API 文档链接
- `README.md` - 在功能列表中添加随机 API

## 未来扩展建议

如果后续有需求,可以添加:

### 1. 图库(Gallery)功能

```go
type Gallery struct {
    ID          uint
    Name        string
    Description string
    CoverImage  uint  // 封面图片 ID
    UserID      uint
    IsPublic    bool
}

// 图片与图库的多对多关系
type ImageGallery struct {
    ImageID   uint
    GalleryID uint
}
```

**新增端点**:
```
GET /api/gallery/:name/random
GET /api/gallery/:id/images
POST /api/gallery
```

### 2. 更多筛选选项

```
GET /api/random?min_width=1920&min_height=1080
GET /api/random?format=jpg
GET /api/random?sort=popular  // 按热度随机
```

### 3. 批量随机

```
GET /api/random/batch?count=10&tags=风景
```

返回多张随机图片的数组,减少 API 调用次数。

### 4. 加权随机

```go
// 根据浏览量、点赞数等加权
query.Order("RANDOM() * (1 + view_count * 0.001)")
```

热门图片出现概率更高。

## 总结

✅ **已实现**:
- 随机图片 API (三种返回方式)
- 按用户和标签筛选
- 访问统计
- 完整文档和演示

❌ **暂不需要**:
- 图库(Gallery)层级
- 复杂的分类系统

💡 **建议**:
- 先使用现有的标签系统
- 根据实际使用情况决定是否增加图库功能
- 当前方案已经足够灵活和强大

## 测试清单

启动后端后,可以测试以下端点:

```bash
# 1. 获取随机图片信息
curl http://localhost:8080/api/random

# 2. 下载随机图片
curl http://localhost:8080/api/random/image -o random.jpg

# 3. 按标签筛选
curl http://localhost:8080/api/random?tags=风景

# 4. 按用户筛选
curl http://localhost:8080/api/random?user_id=1

# 5. 组合筛选
curl http://localhost:8080/api/random/image?user_id=1&tags=风景
```

打开 `docs/random-demo.html` 查看可视化演示。
