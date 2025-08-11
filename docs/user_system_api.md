# 用户系统API文档

## 概述

用户系统基于微信小程序云开发环境，通过拦截器自动处理用户认证和用户信息管理。

## 请求头说明

每个请求都需要包含以下微信小程序云开发环境的请求头：

- `X-WX-OPENID`: 小程序用户 openid，资源复用情况下不存在，参考 X-WX-FROM-OPENID
- `X-WX-APPID`: 小程序 AppID
- `X-WX-UNIONID`: 小程序用户 unionid，并且满足 unionid 获取条件时有
- `X-WX-FROM-OPENID`: 资源复用情况下，小程序用户 openid
- `X-WX-FROM-APPID`: 资源复用情况下，使用方小程序 AppID
- `X-WX-FROM-UNIONID`: 资源复用情况下，小程序用户 unionid，并且满足 unionid 获取条件时有
- `X-WX-ENV`: 所在云环境 ID
- `X-WX-SOURCE`: 调用来源（本次运行是被什么触发）
- `X-Original-Forwarded-For`: 客户端 IPv4 或IPv6 地址

## 用户认证逻辑

1. 拦截器会检查请求头中的 `X-WX-OPENID` 或 `X-WX-FROM-OPENID`
2. 如果用户表中存在该 openId，则查询用户信息并传递给后端接口
3. 如果用户表中不存在该 openId，则返回 401 错误，提示用户需要先注册
4. 用户信息会被存储在请求上下文中，供后续处理器使用

## 用户注册流程

1. 前端首先调用 `/api/auth/check` 检查用户是否存在
2. 如果用户不存在，调用 `/api/auth/register` 进行用户注册
3. 注册成功后，用户可以使用所有需要认证的API接口

## API接口

### 1. 检查用户是否存在

**请求**
```
GET /api/auth/check
```

**响应**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "exists": true,
    "user": {
      "id": 1,
      "username": "user123456",
      "nickname": "用户user123456",
      "avatar": "",
      "bio": "",
      "level": 1,
      "isVerified": false,
      "openid": "wx_openid_123",
      "unionid": "wx_unionid_123",
      "appid": "wx_appid_123",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  }
}
```

### 2. 用户注册

**请求**
```
POST /api/auth/register
Content-Type: application/json

{
  "nickname": "用户昵称",
  "avatar": "头像URL",
  "bio": "个人简介"
}
```

**响应**
```json
{
  "code": 0,
  "msg": "User registered successfully",
  "data": {
    "id": 1,
    "username": "user123456",
    "nickname": "用户昵称",
    "avatar": "头像URL",
    "bio": "个人简介",
    "level": 1,
    "isVerified": false,
    "openid": "wx_openid_123",
    "unionid": "wx_unionid_123",
    "appid": "wx_appid_123",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

### 3. 用户登录

**请求**
```
POST /api/auth/login
```

**响应**
```json
{
  "code": 0,
  "msg": "Login successful",
  "data": {
    "id": 1,
    "username": "user123456",
    "nickname": "用户昵称",
    "avatar": "头像URL",
    "bio": "个人简介",
    "level": 1,
    "isVerified": false,
    "openid": "wx_openid_123",
    "unionid": "wx_unionid_123",
    "appid": "wx_appid_123",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

### 4. 获取当前用户信息

**请求**
```
GET /api/user/profile
```

**响应**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "username": "user123456",
    "nickname": "用户user123456",
    "avatar": "",
    "bio": "",
    "level": 1,
    "isVerified": false,
    "openid": "wx_openid_123",
    "unionid": "wx_unionid_123",
    "appid": "wx_appid_123",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

### 5. 更新用户资料

**请求**
```
PUT /api/user/profile
Content-Type: application/json

{
  "nickname": "新昵称",
  "avatar": "头像URL",
  "bio": "个人简介"
}
```

**响应**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "username": "user123456",
    "nickname": "新昵称",
    "avatar": "头像URL",
    "bio": "个人简介",
    "level": 1,
    "isVerified": false,
    "openid": "wx_openid_123",
    "unionid": "wx_unionid_123",
    "appid": "wx_appid_123",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

### 6. 获取用户列表（管理员功能）

**请求**
```
GET /api/user/list
```

**响应**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "users": [
      {
        "id": 1,
        "username": "user123456",
        "nickname": "用户user123456",
        "avatar": "",
        "bio": "",
        "level": 1,
        "isVerified": false,
        "openid": "wx_openid_123",
        "unionid": "wx_unionid_123",
        "appid": "wx_appid_123",
        "createdAt": "2024-01-01T00:00:00Z",
        "updatedAt": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 20
  }
}
```

## 错误码说明

- `401 Unauthorized`: 缺少必要的微信小程序请求头或用户认证失败
- `400 Bad Request`: 请求参数错误
- `500 Internal Server Error`: 服务器内部错误

## 使用示例

### 前端调用示例

```javascript
// 1. 检查用户是否存在
wx.cloud.callContainer({
  config: {
    env: 'your-env-id'
  },
  path: '/api/auth/check',
  method: 'GET',
  header: {
    'X-WX-OPENID': 'user_openid',
    'X-WX-APPID': 'your_appid',
    'X-WX-ENV': 'your-env-id'
  },
  success: function(res) {
    if (!res.data.data.exists) {
      // 用户不存在，跳转到注册页面
      console.log('用户不存在，需要注册');
    } else {
      // 用户存在，可以直接使用
      console.log('用户已存在:', res.data.data.user);
    }
  }
});

// 2. 用户注册
wx.cloud.callContainer({
  config: {
    env: 'your-env-id'
  },
  path: '/api/auth/register',
  method: 'POST',
  header: {
    'X-WX-OPENID': 'user_openid',
    'X-WX-APPID': 'your_appid',
    'X-WX-ENV': 'your-env-id',
    'Content-Type': 'application/json'
  },
  data: {
    nickname: '用户昵称',
    avatar: '头像URL',
    bio: '个人简介'
  },
  success: function(res) {
    console.log('注册成功:', res.data);
  }
});

// 3. 用户登录
wx.cloud.callContainer({
  config: {
    env: 'your-env-id'
  },
  path: '/api/auth/login',
  method: 'POST',
  header: {
    'X-WX-OPENID': 'user_openid',
    'X-WX-APPID': 'your_appid',
    'X-WX-ENV': 'your-env-id'
  },
  success: function(res) {
    if (res.data.code === 0) {
      console.log('登录成功:', res.data.data);
    } else {
      console.log('用户不存在，需要注册');
    }
  }
});

// 4. 获取用户信息（需要认证）
wx.cloud.callContainer({
  config: {
    env: 'your-env-id'
  },
  path: '/api/user/profile',
  method: 'GET',
  header: {
    'X-WX-OPENID': 'user_openid',
    'X-WX-APPID': 'your_appid',
    'X-WX-ENV': 'your-env-id'
  },
  success: function(res) {
    console.log('用户信息:', res.data);
  }
});

// 5. 更新用户资料
wx.cloud.callContainer({
  config: {
    env: 'your-env-id'
  },
  path: '/api/user/profile',
  method: 'PUT',
  header: {
    'X-WX-OPENID': 'user_openid',
    'X-WX-APPID': 'your_appid',
    'X-WX-ENV': 'your-env-id',
    'Content-Type': 'application/json'
  },
  data: {
    nickname: '新昵称',
    avatar: '头像URL',
    bio: '个人简介'
  },
  success: function(res) {
    console.log('更新成功:', res.data);
  }
});
```

## 注意事项

1. 认证相关接口（`/api/auth/*`）不需要通过用户拦截器
2. 其他API接口都需要通过用户拦截器进行认证
3. 用户拦截器只负责认证，不负责创建用户
4. 新用户注册时会自动生成随机用户名（格式：user + 6位随机数字）
5. 用户信息会自动传递给所有后续的处理器
6. 可以通过 `GetUserFromContext(r)` 函数获取当前用户信息
7. 避免了并发创建用户的问题，确保用户数据的唯一性
