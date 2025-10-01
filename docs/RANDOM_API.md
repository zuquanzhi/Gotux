# 随机图片 API 文档

## 概述

随机图片 API 允许您通过简单的 HTTP 请求获取图库中的随机图片,非常适合用作壁纸、头像、占位图等场景。

## 功能特点

- ✅ 随机返回公开图片
- ✅ 支持按用户筛选
- ✅ 支持按标签筛选
- ✅ 三种返回方式(JSON/直接图片/重定向)
- ✅ 自动统计访问次数
- ✅ 支持自定义域名

## API 端点

### 1. 获取随机图片信息 (JSON)

```http
GET /api/random
```

**用途**: 获取随机图片的完整信息(包括 UUID、文件名、尺寸等)

**响应示例**:
```json
{
  "id": 42,
  "uuid": "550e8400-e29b-41d4-a716-446655440000",
  "file_name": "sunset.jpg",
  "original_name": "beautiful_sunset.jpg",
  "width": 1920,
  "height": 1080,
  "file_size": 524288,
  "mime_type": "image/jpeg",
  "tags": "风景,日落,自然",
  "created_at": "2025-01-15T10:30:00Z",
  "stats": {
    "view_count": 128
  }
}
```

### 2. 直接获取随机图片文件

```http
GET /api/random/image
```

**用途**: 直接返回图片二进制数据,适合在 HTML 中使用

**使用示例**:
```html
<!-- 作为背景图 -->
<div style="background-image: url('https://yourdomain.com/api/random/image')">
</div>

<!-- 作为图片标签 -->
<img src="https://yourdomain.com/api/random/image" alt="随机图片">
```

**响应头**:
- `Content-Type`: 图片的 MIME 类型
- `Cache-Control`: `public, max-age=3600`
- `X-Image-UUID`: 图片的 UUID
- `X-Image-ID`: 图片的数据库 ID

### 3. 重定向到随机图片

```http
GET /api/random/redirect
```

**用途**: 返回 302 重定向到图片的永久链接 (`/i/:uuid`)

**使用场景**:
- 需要永久链接的场合
- SEO 优化
- 社交媒体分享

## 查询参数

所有端点都支持以下查询参数:

### user_id

按用户 ID 筛选图片

```http
GET /api/random/image?user_id=1
```

### tags

按标签筛选图片(部分匹配)

```http
GET /api/random/image?tags=风景
GET /api/random/image?tags=自然,风景
```

### 组合使用

```http
GET /api/random/image?user_id=1&tags=风景
```

## 使用场景

### 1. 网站随机背景

```html
<body style="background-image: url('https://img.example.com/api/random/image?tags=壁纸')">
```

### 2. 占位图服务

```html
<img src="https://img.example.com/api/random/image" width="800" height="600">
```

### 3. API 调用

```javascript
// 获取随机图片信息
fetch('https://img.example.com/api/random?tags=风景')
  .then(res => res.json())
  .then(data => {
    console.log('随机图片:', data);
    // 使用 data.uuid 访问图片
    const imageUrl = `https://img.example.com/i/${data.uuid}`;
  });
```

### 4. Markdown 随机图片

```markdown
![随机图片](https://img.example.com/api/random/image)
```

### 5. 社交媒体卡片

```javascript
// 每次分享都是不同的图片
const shareImage = 'https://img.example.com/api/random/redirect?user_id=1';
```

## 注意事项

### 隐私设置

- 只返回 `is_public = true` 的公开图片
- 私有图片不会出现在随机结果中

### 性能考虑

- 使用 SQLite 的 `RANDOM()` 函数
- 建议为 `is_public` 字段添加索引
- 缓存时间设置为 1 小时

### 访问统计

- `/api/random/image` 会增加访问计数
- `/api/random/redirect` 会增加访问计数
- `/api/random` (JSON) 不增加访问计数

## 高级用法

### 1. 构建自己的随机图床

```javascript
// 前端随机图片组件
function RandomImage({ tags, userId }) {
  const src = `/api/random/image?tags=${tags}&user_id=${userId}`;
  return <img src={src} key={Date.now()} alt="随机图片" />;
}
```

### 2. 定时刷新随机图

```javascript
// 每10秒更换一次背景
setInterval(() => {
  document.body.style.backgroundImage = 
    `url('https://img.example.com/api/random/image?t=${Date.now()}')`;
}, 10000);
```

### 3. 结合标签系统

```javascript
// 根据不同标签获取随机图
const categories = ['风景', '动物', '建筑', '人物'];
const randomTag = categories[Math.floor(Math.random() * categories.length)];
const imageUrl = `/api/random/image?tags=${randomTag}`;
```

## 错误处理

### 404 - 没有找到图片

```json
{
  "error": "没有找到符合条件的图片"
}
```

**原因**:
- 没有公开图片
- 筛选条件过于严格(如指定的标签不存在)

**解决方案**:
- 确保有公开图片
- 放宽筛选条件
- 检查标签拼写

## 完整示例

### HTML 页面

```html
<!DOCTYPE html>
<html>
<head>
  <title>随机图片展示</title>
  <style>
    body {
      margin: 0;
      padding: 20px;
      font-family: Arial, sans-serif;
    }
    .gallery {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
      gap: 20px;
    }
    .card {
      border: 1px solid #ddd;
      border-radius: 8px;
      overflow: hidden;
    }
    .card img {
      width: 100%;
      height: 200px;
      object-fit: cover;
    }
    button {
      margin: 20px 0;
      padding: 10px 20px;
      font-size: 16px;
      cursor: pointer;
    }
  </style>
</head>
<body>
  <h1>随机图片画廊</h1>
  <button onclick="loadRandomImages()">刷新图片</button>
  
  <div class="gallery" id="gallery"></div>

  <script>
    function loadRandomImages() {
      const gallery = document.getElementById('gallery');
      gallery.innerHTML = '';
      
      // 加载6张随机图片
      for (let i = 0; i < 6; i++) {
        fetch('/api/random')
          .then(res => res.json())
          .then(data => {
            const card = document.createElement('div');
            card.className = 'card';
            card.innerHTML = `
              <img src="/i/${data.uuid}" alt="${data.original_name}">
              <div style="padding: 10px;">
                <h3>${data.original_name}</h3>
                <p>尺寸: ${data.width}x${data.height}</p>
                <p>浏览: ${data.stats?.view_count || 0} 次</p>
                ${data.tags ? `<p>标签: ${data.tags}</p>` : ''}
              </div>
            `;
            gallery.appendChild(card);
          });
      }
    }
    
    // 页面加载时自动加载
    loadRandomImages();
  </script>
</body>
</html>
```

### Python 脚本

```python
import requests

# 下载随机图片
def download_random_image(output_path, tags=None, user_id=None):
    url = 'https://img.example.com/api/random/image'
    params = {}
    if tags:
        params['tags'] = tags
    if user_id:
        params['user_id'] = user_id
    
    response = requests.get(url, params=params)
    if response.status_code == 200:
        with open(output_path, 'wb') as f:
            f.write(response.content)
        print(f'图片已保存到 {output_path}')
        print(f'UUID: {response.headers.get("X-Image-UUID")}')
    else:
        print('获取图片失败:', response.json())

# 使用示例
download_random_image('random.jpg', tags='风景')
```

## 未来扩展 (可选)

如果后续有需求,可以添加:

### 图库 (Gallery) 概念

```http
GET /api/gallery/{name}/random
```

**优势**:
- 更好的分类管理
- 可以设置图库封面
- 支持图库描述和权限
- 图库级别的统计

### 更多筛选选项

```http
GET /api/random?min_width=1920&min_height=1080&format=jpg
```

### 批量随机

```http
GET /api/random/batch?count=10&tags=风景
```

返回多张随机图片的数组

## 总结

随机图片 API 提供了灵活的方式来访问和展示您的图片库。通过简单的 URL 就能实现:

- ✅ 网站动态背景
- ✅ 随机占位图
- ✅ 图片轮播效果
- ✅ API 集成
- ✅ 社交媒体分享

无需复杂配置,开箱即用!
