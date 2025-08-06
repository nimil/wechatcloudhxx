# 纯净社区小程序 API 接口文档

## 📋 文档说明

本文档描述了纯净社区小程序的核心API接口，包括首页帖子列表、发布帖子、用户交互等功能。

### 基础信息
- **接口域名**: https://api.example.com
- **请求格式**: application/json
- **响应格式**: application/json
- **认证方式**: Bearer Token（需要登录的接口）

## 🏠 首页帖子列表接口

### 1. 获取帖子列表

**接口地址**: `GET /api/posts`

**接口描述**: 获取首页帖子列表，支持分页、分类筛选和排序

**请求参数**:

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| page | number | 否 | 1 | 页码，从1开始 |
| pageSize | number | 否 | 10 | 每页数量，最大50 |
| category | string | 否 | "all" | 分类代码 |
| sort | string | 否 | "latest" | 排序方式 |

**分类代码说明**:
- `all` - 全部
- `tech` - 技术
- `life` - 生活
- `food` - 美食
- `travel` - 旅行
- `book` - 读书
- `sport` - 运动

**排序方式说明**:
- `latest` - 最新发布
- `hot` - 热门推荐
- `recommend` - 编辑推荐

**请求示例**:
```bash
# 获取最新技术帖子
GET /api/posts?page=1&pageSize=10&category=tech&sort=latest

# 获取热门推荐
GET /api/posts?page=1&pageSize=10&category=all&sort=hot
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": "post_001",
        "title": "分享一个实用的开发技巧",
        "excerpt": "今天在工作中发现了一个很实用的开发技巧，想和大家分享一下。这个技巧可以大大提高开发效率...",
        "content": "完整的帖子内容...",
        "author": {
          "id": "user_001",
          "nickname": "技术达人",
          "avatar": "https://example.com/avatar1.png",
          "bio": "热爱技术的开发者",
          "level": 5,
          "isVerified": true
        },
        "category": "tech",
        "categoryName": "技术",
        "tags": ["技术", "开发", "效率"],
        "images": [
          "https://example.com/post1.jpg",
          "https://example.com/post1-2.jpg"
        ],
        "stats": {
          "likes": 128,
          "comments": 32,
          "views": 1024,
          "shares": 15
        },
        "isLiked": false,
        "isCollected": false,
        "createdAt": "2024-01-15T10:30:00Z",
        "updatedAt": "2024-01-15T10:30:00Z"
      }
    ],
    "pagination": {
      "current": 1,
      "pageSize": 10,
      "total": 156,
      "hasMore": true
    }
  }
}
```

### 2. 获取热门话题

**接口地址**: `GET /api/topics/hot`

**接口描述**: 获取首页热门话题列表

**请求参数**: 无

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": "topic_001",
      "name": "技术",
      "icon": "💻",
      "code": "tech",
      "postCount": 1250,
      "followCount": 890,
      "isFollowed": false
    },
    {
      "id": "topic_002",
      "name": "生活",
      "icon": "🏠",
      "code": "life",
      "postCount": 890,
      "followCount": 567,
      "isFollowed": true
    },
    {
      "id": "topic_003",
      "name": "美食",
      "icon": "🍜",
      "code": "food",
      "postCount": 650,
      "followCount": 420,
      "isFollowed": false
    }
  ]
}
```

### 3. 获取分类列表

**接口地址**: `GET /api/categories`

**接口描述**: 获取所有分类列表

**请求参数**: 无

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": "cat_001",
      "name": "全部",
      "code": "all",
      "icon": "📋",
      "postCount": 2560
    },
    {
      "id": "cat_002",
      "name": "技术",
      "code": "tech",
      "icon": "💻",
      "postCount": 1250
    },
    {
      "id": "cat_003",
      "name": "生活",
      "code": "life",
      "icon": "🏠",
      "postCount": 890
    },
    {
      "id": "cat_004",
      "name": "美食",
      "code": "food",
      "icon": "🍜",
      "postCount": 650
    },
    {
      "id": "cat_005",
      "name": "旅行",
      "code": "travel",
      "icon": "✈️",
      "postCount": 450
    },
    {
      "id": "cat_006",
      "name": "读书",
      "code": "book",
      "icon": "📚",
      "postCount": 320
    },
    {
      "id": "cat_007",
      "name": "运动",
      "code": "sport",
      "icon": "🏃",
      "postCount": 280
    }
  ]
}
```

## ✍️ 发布帖子接口

### 1. 发布帖子

**接口地址**: `POST /api/posts`

**接口描述**: 发布新帖子

**请求头**:
```
Content-Type: application/json
Authorization: Bearer {token}
```

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| title | string | 是 | 帖子标题，最大50字符 |
| content | string | 是 | 帖子内容，最大1000字符 |
| category | string | 是 | 分类代码 |
| tags | array | 否 | 标签数组，最大5个 |
| images | array | 否 | 图片URL数组，最大9张 |
| isPublic | boolean | 否 | 是否公开，默认true |

**请求示例**:
```json
{
  "title": "分享一个实用的开发技巧",
  "content": "今天在工作中发现了一个很实用的开发技巧，想和大家分享一下。这个技巧可以大大提高开发效率，特别是在处理大量数据时。\n\n主要特点：\n1. 性能优化明显\n2. 代码更简洁\n3. 易于维护\n\n希望这个分享对大家有帮助！",
  "category": "tech",
  "tags": ["技术", "开发", "效率"],
  "images": [
    "https://example.com/uploaded_image1.jpg",
    "https://example.com/uploaded_image2.jpg"
  ],
  "isPublic": true
}
```

**响应数据**:
```json
{
  "code": 200,
  "message": "发布成功",
  "data": {
    "postId": "post_001",
    "createdAt": "2024-01-15T10:30:00Z",
    "url": "https://api.example.com/posts/post_001"
  }
}
```

### 2. 获取发布分类列表

**接口地址**: `GET /api/categories/publish`

**接口描述**: 获取可用于发布的分类列表

**请求参数**: 无

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": "cat_002",
      "name": "技术",
      "code": "tech",
      "icon": "💻",
      "description": "技术分享、开发经验、编程技巧"
    },
    {
      "id": "cat_003",
      "name": "生活",
      "code": "life",
      "icon": "🏠",
      "description": "日常生活、心情分享、生活感悟"
    },
    {
      "id": "cat_004",
      "name": "美食",
      "code": "food",
      "icon": "🍜",
      "description": "美食制作、餐厅推荐、食谱分享"
    }
  ]
}
```

## 💬 交互功能接口

### 1. 点赞/取消点赞

**接口地址**: `POST /api/posts/{postId}/like`

**接口描述**: 点赞或取消点赞帖子

**请求头**:
```
Content-Type: application/json
Authorization: Bearer {token}
```

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| action | string | 是 | like(点赞) 或 unlike(取消点赞) |

**请求示例**:
```json
{
  "action": "like"
}
```

**响应数据**:
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "isLiked": true,
    "likesCount": 129
  }
}
```

### 2. 发表评论

**接口地址**: `POST /api/posts/{postId}/comments`

**接口描述**: 对帖子发表评论

**请求头**:
```
Content-Type: application/json
Authorization: Bearer {token}
```

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| content | string | 是 | 评论内容，最大200字符 |
| parentId | string | 否 | 回复的评论ID |

**请求示例**:
```json
{
  "content": "这个技巧确实很实用，感谢分享！",
  "parentId": "comment_001"
}
```

**响应数据**:
```json
{
  "code": 200,
  "message": "评论成功",
  "data": {
    "commentId": "comment_002",
    "createdAt": "2024-01-15T11:00:00Z"
  }
}
```

### 3. 获取评论列表

**接口地址**: `GET /api/posts/{postId}/comments`

**接口描述**: 获取帖子的评论列表

**请求参数**:

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| page | number | 否 | 1 | 页码 |
| pageSize | number | 否 | 20 | 每页数量 |

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": "comment_001",
        "content": "这个技巧确实很实用，感谢分享！",
        "author": {
          "id": "user_002",
          "nickname": "开发者小王",
          "avatar": "https://example.com/avatar2.png"
        },
        "postId": "post_001",
        "parentId": null,
        "likes": 12,
        "isLiked": false,
        "createdAt": "2024-01-15T11:00:00Z",
        "replies": [
          {
            "id": "comment_002",
            "content": "我也觉得很有用！",
            "author": {
              "id": "user_003",
              "nickname": "前端工程师",
              "avatar": "https://example.com/avatar3.png"
            },
            "parentId": "comment_001",
            "likes": 5,
            "isLiked": false,
            "createdAt": "2024-01-15T11:30:00Z"
          }
        ]
      }
    ],
    "pagination": {
      "current": 1,
      "pageSize": 20,
      "total": 45,
      "hasMore": true
    }
  }
}
```

## 🔧 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 📝 注意事项

1. 所有需要认证的接口都需要在请求头中携带 `Authorization: Bearer {token}`
2. 图片上传需要先调用文件上传接口获取URL
3. 分页参数从1开始计数
4. 时间格式统一使用ISO 8601格式
5. 文件大小限制：图片最大5MB，支持jpg、png、gif格式 