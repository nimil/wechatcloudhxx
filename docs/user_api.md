# 用户API文档

## 接口地址
`POST /api/test/user` - 创建用户
`GET /api/test/user` - 分页查询用户列表

## 功能描述
- POST: 创建新用户，传入用户名和密码
- GET: 分页查询所有用户列表

## 请求参数

### POST 创建用户
```json
{
  "username": "testuser",
  "password": "123456"
}
```

### GET 分页查询用户
查询参数：
- `page`: 页码（可选，默认1）
- `pageSize`: 每页大小（可选，默认10，最大100）

示例：`GET /api/test/user?page=1&pageSize=10`

### 微信请求头（POST创建用户时自动获取）
系统会自动从请求头中获取以下微信用户信息：
- `X-WX-OPENID`: 小程序用户 openid
- `X-WX-UNIONID`: 小程序用户 unionid
- `X-WX-APPID`: 小程序 AppID
- `X-WX-FROM-OPENID`: 资源复用情况下的小程序用户 openid
- `X-WX-FROM-UNIONID`: 资源复用情况下的小程序用户 unionid
- `X-WX-FROM-APPID`: 资源复用情况下的小程序 AppID

## 响应格式

### POST 成功响应
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "username": "testuser",
    "openid": "wx_openid_123",
    "unionid": "wx_unionid_456",
    "appid": "wx_appid_789"
  }
}
```

### GET 成功响应
```json
{
  "code": 0,
  "data": {
    "users": [
      {
        "id": 1,
        "username": "user1",
        "openid": "wx_openid_123",
        "unionid": "wx_unionid_456",
        "appid": "wx_appid_789"
      },
      {
        "id": 2,
        "username": "user2",
        "openid": "wx_openid_456",
        "unionid": "wx_unionid_789",
        "appid": "wx_appid_123"
      }
    ],
    "total": 25,
    "page": 1,
    "pageSize": 10,
    "pages": 3
  }
}
```

### 失败响应
```json
{
  "code": -1,
  "errorMsg": "用户名已存在"
}
```

## 错误码说明
- `code: 0` - 成功
- `code: -1` - 失败，具体错误信息在 `errorMsg` 字段

## 常见错误
- 用户名不能为空
- 密码不能为空
- 用户名已存在
- 该微信用户已存在
- 解析请求参数失败
- 创建用户失败
- 页码参数无效
- 每页大小参数无效，范围1-100
- 查询用户列表失败

## 使用示例

### POST 创建用户

#### curl 命令
```bash
curl -X POST http://localhost:80/api/test/user \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456"
  }'
```

#### JavaScript 示例
```javascript
fetch('/api/test/user', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    username: 'testuser',
    password: '123456'
  })
})
.then(response => response.json())
.then(data => console.log(data));
```

### GET 分页查询用户

#### curl 命令
```bash
# 查询第1页，每页10条
curl "http://localhost:80/api/test/user?page=1&pageSize=10"

# 查询第2页，每页20条
curl "http://localhost:80/api/test/user?page=2&pageSize=20"
```

#### JavaScript 示例
```javascript
// 查询第1页，每页10条
fetch('/api/test/user?page=1&pageSize=10')
.then(response => response.json())
.then(data => console.log(data));

// 查询第2页，每页20条
fetch('/api/test/user?page=2&pageSize=20')
.then(response => response.json())
.then(data => console.log(data));
```

## 注意事项
1. 用户名具有唯一性，不能重复
2. 微信OpenID和UnionID具有唯一性，不能重复
3. 密码目前以明文存储，生产环境建议加密
4. POST 接口用于创建用户，GET 接口用于分页查询用户
5. 返回的用户信息不包含密码字段
6. 分页查询默认每页10条，最大100条
7. 分页查询按用户ID倒序排列（最新创建的用户在前）
8. 系统会自动从请求头中获取微信用户信息（X-WX-OPENID、X-WX-UNIONID、X-WX-APPID）
9. 支持资源复用场景，会自动获取X-WX-FROM-*字段 