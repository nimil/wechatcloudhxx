# 我的帖子API文档

## 接口概述

获取当前登录用户发布的未删除的所有帖子列表。

## 接口信息

- **接口地址**: `GET /api/posts/my`
- **请求方法**: GET
- **需要认证**: 是（需要用户登录）
- **内容类型**: application/json

## 请求参数

### 查询参数

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| page | int | 否 | 1 | 页码，从1开始 |
| pageSize | int | 否 | 10 | 每页数量，最大50 |

### 请求头

```
X-WX-OPENID: 用户openid
X-WX-APPID: 小程序AppID
X-WX-ENV: 云环境ID
```

## 响应格式

### 成功响应

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 123,
        "title": "帖子标题",
        "excerpt": "帖子摘要...",
        "content": "帖子完整内容",
        "author": {
          "id": 456,
          "nickname": "用户昵称",
          "avatar": "头像URL",
          "bio": "个人简介",
          "level": 1,
          "isVerified": false
        },
        "category": "general",
        "categoryName": "综合",
        "tags": ["标签1", "标签2"],
        "images": ["图片URL1", "图片URL2"],
        "stats": {
          "likes": 10,
          "comments": 5,
          "views": 100,
          "shares": 2
        },
        "isLiked": false,
        "isCollected": false,
        "createdAt": "2024-01-01T12:00:00Z",
        "updatedAt": "2024-01-01T12:00:00Z"
      }
    ],
    "pagination": {
      "current": 1,
      "pageSize": 10,
      "total": 25,
      "hasMore": true
    }
  }
}
```

### 错误响应

#### 未授权错误

```json
{
  "code": 401,
  "message": "Unauthorized"
}
```

#### 服务器错误

```json
{
  "code": 500,
  "message": "获取用户帖子列表失败: 数据库连接错误"
}
```

## 字段说明

### PostDetail 字段

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int64 | 帖子ID |
| title | string | 帖子标题 |
| excerpt | string | 帖子摘要（前200字符） |
| content | string | 帖子完整内容 |
| author | UserInfo | 作者信息 |
| category | string | 分类代码 |
| categoryName | string | 分类名称 |
| tags | []string | 标签列表 |
| images | []string | 图片URL列表 |
| stats | PostStats | 统计信息 |
| isLiked | bool | 是否已点赞（我的帖子默认为false） |
| isCollected | bool | 是否已收藏 |
| createdAt | string | 创建时间 |
| updatedAt | string | 更新时间 |

### UserInfo 字段

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int64 | 用户ID |
| nickname | string | 用户昵称 |
| avatar | string | 头像URL |
| bio | string | 个人简介 |
| level | int | 用户等级 |
| isVerified | bool | 是否认证用户 |

### PostStats 字段

| 字段名 | 类型 | 说明 |
|--------|------|------|
| likes | int | 点赞数 |
| comments | int | 评论数 |
| views | int | 浏览量 |
| shares | int | 分享数 |

### Pagination 字段

| 字段名 | 类型 | 说明 |
|--------|------|------|
| current | int | 当前页码 |
| pageSize | int | 每页数量 |
| total | int64 | 总记录数 |
| hasMore | bool | 是否有更多数据 |

## 使用示例

### 前端调用示例

```javascript
// 获取我的帖子列表
wx.cloud.callContainer({
  config: {
    env: 'your-env-id'
  },
  path: '/api/posts/my?page=1&pageSize=10',
  method: 'GET',
  header: {
    'X-WX-OPENID': 'user_openid',
    'X-WX-APPID': 'your_appid',
    'X-WX-ENV': 'your-env-id'
  },
  success: function(res) {
    if (res.data.code === 200) {
      console.log('我的帖子列表:', res.data.data);
      
      // 处理帖子列表
      const posts = res.data.data.list;
      const pagination = res.data.data.pagination;
      
      console.log('帖子数量:', posts.length);
      console.log('总帖子数:', pagination.total);
      console.log('是否有更多:', pagination.hasMore);
      
      // 遍历帖子
      posts.forEach(post => {
        console.log('帖子标题:', post.title);
        console.log('点赞数:', post.stats.likes);
        console.log('发布时间:', post.createdAt);
      });
    } else {
      console.log('获取失败:', res.data.message);
    }
  },
  fail: function(err) {
    console.log('请求失败:', err);
  }
});
```

### 分页加载示例

```javascript
let currentPage = 1;
let hasMore = true;

// 加载第一页
function loadMyPosts(page = 1) {
  if (!hasMore && page > 1) {
    return;
  }
  
  wx.cloud.callContainer({
    config: {
      env: 'your-env-id'
    },
    path: `/api/posts/my?page=${page}&pageSize=10`,
    method: 'GET',
    header: {
      'X-WX-OPENID': 'user_openid',
      'X-WX-APPID': 'your_appid',
      'X-WX-ENV': 'your-env-id'
    },
    success: function(res) {
      if (res.data.code === 200) {
        const data = res.data.data;
        
        if (page === 1) {
          // 第一页，替换数据
          this.setData({
            myPosts: data.list
          });
        } else {
          // 加载更多，追加数据
          this.setData({
            myPosts: [...this.data.myPosts, ...data.list]
          });
        }
        
        // 更新分页信息
        hasMore = data.pagination.hasMore;
        currentPage = data.pagination.current;
        
        console.log('当前页:', currentPage);
        console.log('是否有更多:', hasMore);
      }
    }
  });
}

// 加载更多
function loadMore() {
  if (hasMore) {
    loadMyPosts(currentPage + 1);
  }
}
```

## 注意事项

1. **认证要求**: 必须提供有效的用户认证信息
2. **分页限制**: 每页最大50条记录
3. **排序规则**: 按创建时间倒序排列（最新发布在前）
4. **数据范围**: 只返回未删除的帖子
5. **点赞状态**: 我的帖子列表中不显示点赞状态（默认为false）

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 401 | 未授权，需要用户登录 |
| 405 | 请求方法不允许 |
| 500 | 服务器内部错误 |

## 相关接口

- [获取帖子列表](/api/posts) - 获取所有公开帖子
- [获取帖子详情](/api/posts/{id}) - 获取单个帖子详情
- [删除帖子](/api/posts/{id}) - 删除指定帖子

