# 随机图片 API - 完整实现总结

## 实现日期: 2025-10-02

## 需求分析

用户希望增加一个 API 调用功能,可以使用域名和请求来访问图库中的随机图片。

### 关键决策: 不增加图库(Gallery)层级

**理由**:
1. 当前系统已有 `tags` 字段,可以实现分类
2. 通过 `user_id` 可以区分不同用户的图片
3. 两者结合已经足够灵活,无需额外复杂度
4. 可以快速实现,无需数据库迁移

**如果未来需要图库**:
- 可以使用特殊标签格式(如 `gallery:风景`)
- 或新增 Gallery 表和多对多关系

## 功能实现

### 1. 后端 API 端点

#### 新增 Controller 函数 (`backend/controllers/image.go`)

```go
// GetRandomImage - 返回随机图片信息(JSON)
func GetRandomImage(c *gin.Context) {
    query := models.DB.Where("is_public = ?", true)
    
    // 支持按 user_id 筛选
    if userID := c.Query("user_id"); userID != "" {
        query = query.Where("user_id = ?", userID)
    }
    
    // 支持按 tags 筛选
    if tags := c.Query("tags"); tags != "" {
        query = query.Where("tags LIKE ?", "%"+tags+"%")
    }
    
    var image models.Image
    query.Order("RANDOM()").First(&image)
    c.JSON(http.StatusOK, image)
}

// ServeRandomImage - 直接返回随机图片文件
func ServeRandomImage(c *gin.Context) {
    // 同样的筛选逻辑
    // ...
    models.IncrementViewCount(image.ID)
    c.File(fullPath)
}

// RedirectRandomImage - 重定向到随机图片
func RedirectRandomImage(c *gin.Context) {
    // 同样的筛选逻辑
    // ...
    imageURL := fmt.Sprintf("%s/i/%s", baseURL, image.UUID)
    c.Redirect(http.StatusFound, imageURL)
}
```

#### 路由配置 (`backend/routes/routes.go`)

```go
// 公开路由,无需认证
api.GET("/random", controllers.GetRandomImage)
api.GET("/random/image", controllers.ServeRandomImage)
api.GET("/random/redirect", controllers.RedirectRandomImage)
```

### 2. API 端点说明

| 端点 | 方法 | 说明 | 返回 |
|------|------|------|------|
| `/api/random` | GET | 获取随机图片信息 | JSON |
| `/api/random/image` | GET | 直接返回图片文件 | Binary |
| `/api/random/redirect` | GET | 重定向到图片 | 302 Redirect |

### 3. 查询参数

所有端点都支持:
- `user_id` - 筛选指定用户的图片
- `tags` - 筛选带特定标签的图片(部分匹配)

### 4. 使用示例

#### HTML 中使用

```html
<!-- 直接显示随机图片 -->
<img src="http://localhost:8080/api/random/image" alt="随机图片">

<!-- 带标签筛选 -->
<img src="http://localhost:8080/api/random/image?tags=风景" alt="风景图">

<!-- 指定用户 -->
<img src="http://localhost:8080/api/random/image?user_id=1" alt="用户1的图片">

<!-- 作为背景 -->
<div style="background-image: url('http://localhost:8080/api/random/image?tags=壁纸')"></div>
```

#### JavaScript 调用

```javascript
// 获取随机图片信息
fetch('/api/random?tags=风景')
  .then(res => res.json())
  .then(data => {
    console.log('图片UUID:', data.uuid);
    console.log('尺寸:', data.width, 'x', data.height);
    console.log('浏览量:', data.stats?.view_count);
  });

// 动态更换图片
function refreshRandomImage() {
  const img = document.getElementById('random-img');
  img.src = `/api/random/image?t=${Date.now()}`;
}

// 每10秒刷新
setInterval(refreshRandomImage, 10000);
```

#### Markdown 中使用

```markdown
![随机图片](http://localhost:8080/api/random/image)
![风景图](http://localhost:8080/api/random/image?tags=风景)
```

### 5. 前端 API 封装 (`frontend/src/api/image.js`)

```javascript
// 获取随机图片信息(使用原生 fetch 避免添加 token)
export function getRandomImage(params) {
  const queryString = new URLSearchParams(params).toString()
  const url = `/api/random${queryString ? '?' + queryString : ''}`
  return fetch(url).then(res => res.json())
}

// 获取随机图片 URL
export function getRandomImageUrl(params) {
  const queryString = new URLSearchParams(params).toString()
  return `/api/random/image${queryString ? '?' + queryString : ''}`
}
```

## 文档更新

### 新增文档

1. **`docs/RANDOM_API.md`** (350+ 行)
   - 详细的 API 使用指南
   - 所有端点说明和示例
   - 各种使用场景
   - HTML/JS/Python 完整示例
   - 高级用法和最佳实践

2. **`docs/RANDOM_FEATURE.md`** (200+ 行)
   - 功能实现总结
   - 技术决策说明
   - 未来扩展建议
   - 测试清单

3. **`docs/random-demo.html`** (400+ 行)
   - 交互式演示页面
   - 可视化筛选条件
   - 实时加载随机图片
   - 完整代码示例

### 更新文档

1. **`docs/API.md`**
   - 添加随机 API 端点说明
   - 请求/响应示例

2. **`docs/README.md`**
   - 添加随机 API 文档链接

3. **`README.md`**
   - 功能列表添加"随机图片 API"
   - 文档链接更新

4. **`docs/CHANGELOG.md`**
   - 新增 v2.2.0 版本记录

## Bug 修复

### Bug #4: 管理员页面图片预览遮挡

**问题**: 在以下页面查看大图时,UI 元素会遮挡预览窗口:
- 管理员图片管理页面
- 上传成功页面
- 我的图片页面

**修复**: 为所有 `el-image` 组件添加:
```vue
<el-image
  :z-index="9999"
  :preview-teleported="true"
  ...
/>
```

**影响文件**:
- `frontend/src/views/admin/Images.vue`
- `frontend/src/views/Upload.vue`
- `frontend/src/views/Images.vue`
- `frontend/src/views/Dashboard.vue` (已在之前修复)

### Bug #5: 随机 API 认证问题

**问题**: 前端使用 `request.js` 调用随机 API 时,会自动添加 Authorization 头,导致可能出现 403 错误。

**原因**: `request.js` 的拦截器会为所有请求添加 token(如果存在)。

**解决方案**: 
1. 创建独立的随机图片 API 函数
2. 使用原生 `fetch` 而不是 axios
3. 避免触发请求拦截器

**代码**:
```javascript
// 不经过 request.js,直接使用 fetch
export function getRandomImage(params) {
  const queryString = new URLSearchParams(params).toString()
  const url = `/api/random${queryString ? '?' + queryString : ''}`
  return fetch(url).then(res => res.json())
}
```

## 技术实现细节

### 随机算法

使用 SQLite 的 `RANDOM()` 函数:
```go
query.Order("RANDOM()").First(&image)
```

**性能考虑**:
- 适合中小型数据集(< 10 万张)
- 大数据量时可以考虑:
  - 缓存策略
  - 预生成随机序列
  - ID 范围随机

### 隐私保护

```go
query := models.DB.Where("is_public = ?", true)
```
只返回公开图片,保护用户隐私。

### 访问统计

- `ServeRandomImage` - 增加访问计数 ✅
- `RedirectRandomImage` - 增加访问计数 ✅
- `GetRandomImage` - 不增加(仅查询信息)

### 缓存策略

```go
c.Header("Cache-Control", "public, max-age=3600")
c.Header("X-Image-UUID", image.UUID)
c.Header("X-Image-ID", strconv.Itoa(int(image.ID)))
```

- 缓存 1 小时
- 通过响应头提供 UUID 和 ID 信息

## 测试方法

### 命令行测试

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
curl "http://localhost:8080/api/random/image?user_id=1&tags=风景" -o landscape.jpg

# 6. 查看响应头
curl -I http://localhost:8080/api/random/image
```

### 浏览器测试

1. 打开 `docs/random-demo.html`
2. 设置筛选条件
3. 点击"加载图片"
4. 查看随机图片展示
5. 检查控制台输出

### 前端集成测试

```javascript
import { getRandomImage, getRandomImageUrl } from '@/api/image'

// 测试 1: 获取随机图片信息
const image = await getRandomImage({ tags: '风景' })
console.log('随机图片:', image)

// 测试 2: 在页面中显示
const imageUrl = getRandomImageUrl({ user_id: 1 })
document.getElementById('img').src = imageUrl
```

## 应用场景

### 1. 网站随机背景

```html
<body style="background-image: url('/api/random/image?tags=壁纸')">
```

### 2. 占位图服务

```html
<img src="/api/random/image" width="800" height="600">
```

### 3. 图片轮播

```javascript
const slides = [
  '/api/random/image?tags=风景',
  '/api/random/image?tags=动物',
  '/api/random/image?tags=建筑'
];
```

### 4. API 对外服务

```
https://img.yourdomain.com/api/random/image?tags=科技
```

### 5. 社交媒体卡片

```html
<meta property="og:image" content="https://img.yourdomain.com/api/random/redirect">
```

## 性能优化建议

### 当前方案
- SQLite RANDOM() 函数
- 适合 < 10 万张图片
- 简单直接,无需额外存储

### 优化方案(大数据量)

#### 1. 缓存随机列表
```go
var randomCache []uint
func getRandomImageID() uint {
    if len(randomCache) == 0 {
        // 重新生成缓存
        DB.Model(&Image{}).Where("is_public = ?", true).Pluck("id", &randomCache)
        // 打乱顺序
        rand.Shuffle(len(randomCache), func(i, j int) {
            randomCache[i], randomCache[j] = randomCache[j], randomCache[i]
        })
    }
    id := randomCache[0]
    randomCache = randomCache[1:]
    return id
}
```

#### 2. ID 范围随机
```go
var minID, maxID uint
DB.Model(&Image{}).Select("MIN(id), MAX(id)").Scan(&minID, &maxID)
randomID := rand.Intn(int(maxID-minID)) + int(minID)
DB.Where("id >= ? AND is_public = ?", randomID, true).First(&image)
```

#### 3. Redis 缓存
```go
// 缓存热门随机图片
redisClient.SAdd("random:images", imageID)
randomID := redisClient.SRandMember("random:images")
```

## 未来扩展建议

### 1. 图库(Gallery)功能

如果用户反馈需要更复杂的分类:

```go
type Gallery struct {
    ID          uint
    Name        string
    Description string
    CoverImageID uint
    UserID      uint
    IsPublic    bool
}

// 多对多关系
type ImageGallery struct {
    ImageID   uint
    GalleryID uint
}

// 新端点
GET /api/gallery/:name/random
GET /api/gallery/:id/images
```

### 2. 更多筛选选项

```
GET /api/random?min_width=1920&min_height=1080
GET /api/random?format=jpg
GET /api/random?order=popular  // 按热度加权随机
```

### 3. 批量随机

```
GET /api/random/batch?count=10&tags=风景
```

返回数组,减少 API 调用。

### 4. 加权随机

```go
// 热门图片出现概率更高
query.Order("RANDOM() * (1 + view_count * 0.001)")
```

### 5. 时间筛选

```
GET /api/random?date_from=2025-01-01&date_to=2025-12-31
```

### 6. 排除已显示

```
GET /api/random?exclude=1,2,3,4,5
```

避免短时间内重复。

## 总结

### ✅ 已完成

- [x] 随机图片 API 三种返回方式
- [x] 用户和标签筛选
- [x] 访问统计
- [x] 自定义域名支持
- [x] 完整文档(400+ 行)
- [x] 交互式演示页面
- [x] 前端 API 封装
- [x] Bug 修复(图片预览遮挡)
- [x] 测试验证

### ⭐ 核心优势

1. **无需图库层级** - 利用现有 tags 系统
2. **开箱即用** - 无需额外配置
3. **灵活筛选** - user_id + tags 组合
4. **多种返回** - JSON/Binary/Redirect
5. **完整文档** - 详细示例和最佳实践

### 🚀 快速开始

```bash
# 1. 启动后端
cd backend
go run main.go

# 2. 测试 API
curl http://localhost:8080/api/random

# 3. 浏览器查看演示
open docs/random-demo.html
```

### 📚 文档位置

- 使用指南: `docs/RANDOM_API.md`
- 功能说明: `docs/RANDOM_FEATURE.md`
- 演示页面: `docs/random-demo.html`
- API 文档: `docs/API.md`
- 更新日志: `docs/CHANGELOG.md`

---

**实现完成时间**: 2025-10-02  
**版本**: v2.2.0  
**作者**: GitHub Copilot
