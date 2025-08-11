# 帖子API文档

## 概述

帖子API提供了帖子的创建、查询、详情查看和删除功能。

## API接口

### 1. 获取帖子列表

**请求**
```
GET /api/posts?page=1&pageSize=10&category=all&sort=latest
```

**参数说明**
- `page`: 页码，默认为1
- `pageSize`: 每页数量，默认为10，最大50
- `category`: 分类代码，默认为"all"表示所有分类
- `sort`: 排序方式，可选值：latest（最新）、hot（热门）、recommend（推荐）

**响应**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "帖子标题",
        "excerpt": "帖子摘要",
        "content": "帖子内容",
        "author": {
          "id": 1,
          "nickname": "用户昵称",
          "avatar": "头像URL",
          "bio": "个人简介",
          "level": 1,
          "isVerified": false
        },
        "category": "tech",
        "categoryName": "技术",
        "tags": ["标签1", "标签2"],
        "images": ["图片1", "图片2"],
        "stats": {
          "likes": 10,
          "comments": 5,
          "views": 100,
          "shares": 2
        },
        "isLiked": false,
        "isCollected": false,
        "createdAt": "2024-01-01T00:00:00Z",
        "updatedAt": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "current": 1,
      "pageSize": 10,
      "total": 100,
      "hasMore": true
    }
  }
}
```

### 2. 创建帖子

**请求**
```
POST /api/posts
Content-Type: application/json

{
  "title": "帖子标题",
  "content": "帖子内容",
  "category": "tech",
  "tags": ["标签1", "标签2"],
  "images": ["图片1", "图片2"],
  "isPublic": true
}
```

**响应**
```json
{
  "code": 200,
  "message": "发布成功",
  "data": {
    "postId": 1,
    "createdAt": "2024-01-01T00:00:00Z",
    "url": "https://api.example.com/posts/1"
  }
}
```

### 3. 获取帖子详情

**请求**
```
GET /api/posts/{postId}
```

**响应**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "title": "帖子标题",
    "excerpt": "帖子摘要",
    "content": "帖子完整内容",
    "author": {
      "id": 1,
      "nickname": "用户昵称",
      "avatar": "头像URL",
      "bio": "个人简介",
      "level": 1,
      "isVerified": false
    },
    "category": "tech",
    "categoryName": "技术",
    "tags": ["标签1", "标签2"],
    "images": ["图片1", "图片2"],
    "stats": {
      "likes": 10,
      "comments": 5,
      "views": 100,
      "shares": 2
    },
    "isLiked": false,
    "isCollected": false,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

**特性**
- 自动增加浏览量
- 返回完整的帖子内容
- 包含用户点赞状态

### 4. 删除帖子（逻辑删除）

**请求**
```
DELETE /api/posts/{postId}
```

**响应**
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

**权限要求**
- 只有帖子作者可以删除自己的帖子
- 需要用户认证

### 5. 帖子点赞

**请求**
```
POST /api/posts/{postId}/like
Content-Type: application/json

{
  "action": "like" // 或 "unlike"
}
```

**响应**
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "isLiked": true,
    "likesCount": 11
  }
}
```

### 6. 获取帖子评论

**请求**
```
GET /api/posts/{postId}/comments?page=1&pageSize=10
```

**响应**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "content": "评论内容",
        "author": {
          "id": 1,
          "nickname": "用户昵称",
          "avatar": "头像URL"
        },
        "likes": 5,
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "current": 1,
      "pageSize": 10,
      "total": 50,
      "hasMore": true
    }
  }
}
```

### 7. 发表评论

**请求**
```
POST /api/posts/{postId}/comments
Content-Type: application/json

{
  "content": "评论内容",
  "parentId": null // 父评论ID，回复评论时使用
}
```

**响应**
```json
{
  "code": 200,
  "message": "评论成功",
  "data": {
    "id": 1,
    "content": "评论内容",
    "author": {
      "id": 1,
      "nickname": "用户昵称",
      "avatar": "头像URL"
    },
    "createdAt": "2024-01-01T00:00:00Z"
  }
}
```

## 错误码说明

- `200`: 成功
- `400`: 请求参数错误
- `401`: 未授权（需要登录）
- `403`: 无权限（如删除他人帖子）
- `404`: 资源不存在
- `500`: 服务器内部错误

## 逻辑删除说明

1. **删除方式**: 使用逻辑删除，不会物理删除数据
2. **删除字段**: `is_deleted` 字段标记为 `true`
3. **查询过滤**: 所有查询都会自动过滤已删除的帖子
4. **权限控制**: 只有帖子作者可以删除自己的帖子
5. **数据恢复**: 支持通过 `Restore` 方法恢复已删除的帖子

## 前端调用示例

```javascript
// 获取帖子详情
wx.cloud.callContainer({
  config: {
    env: 'your-env-id'
  },
  path: '/api/posts/123',
  method: 'GET',
  header: {
    'X-WX-OPENID': 'user_openid',
    'X-WX-APPID': 'your_appid',
    'X-WX-ENV': 'your-env-id'
  },
  success: function(res) {
    console.log('帖子详情:', res.data);
  }
});

// 删除帖子
wx.cloud.callContainer({
  config: {
    env: 'your-env-id'
  },
  path: '/api/posts/123',
  method: 'DELETE',
  header: {
    'X-WX-OPENID': 'user_openid',
    'X-WX-APPID': 'your_appid',
    'X-WX-ENV': 'your-env-id'
  },
  success: function(res) {
    console.log('删除成功:', res.data);
  }
});
```

## 注意事项

1. 所有接口都需要通过用户拦截器进行认证
2. 帖子详情接口会自动增加浏览量
3. 删除帖子采用逻辑删除，数据不会丢失
4. 只有帖子作者可以删除自己的帖子
5. 已删除的帖子不会出现在列表中
