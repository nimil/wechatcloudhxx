# 纯净社区小程序后端

[![Go Version](https://img.shields.io/badge/Go-1.16+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![WeChat Cloud](https://img.shields.io/badge/WeChat-Cloud%20Run-orange.svg)](https://developers.weixin.qq.com/miniprogram/dev/wxcloudrun/)

一个基于微信云托管的纯净社区小程序后端服务，提供完整的社区功能，包括用户系统、内容发布、互动交流等核心功能。

## 🚀 项目特性

- **微信生态集成** - 基于微信小程序和微信云托管
- **内容安全** - 集成微信内容安全检测，保障社区内容质量
- **用户系统** - 完整的用户注册、登录、认证体系
- **内容管理** - 支持帖子发布、分类管理、标签系统
- **互动功能** - 点赞、评论、收藏等社交功能
- **图片处理** - 支持图片上传和内容安全检测
- **RESTful API** - 标准的REST API设计，易于集成

## 📋 功能模块

### 用户系统
- 微信授权登录
- 用户信息管理
- 用户等级和认证体系
- 个人资料设置

### 内容系统
- 帖子发布和管理
- 分类和标签系统
- 内容审核和安全检测
- 图片上传和处理

### 互动功能
- 点赞和取消点赞
- 评论和回复系统
- 用户关注功能
- 内容分享

### 管理功能
- 内容审核
- 用户管理
- 数据统计
- 系统配置

## 🏗️ 技术架构

### 后端技术栈
- **语言**: Go 1.16+
- **框架**: 原生 HTTP 包
- **数据库**: MySQL
- **ORM**: GORM
- **部署**: 微信云托管
- **容器**: Docker

### 项目结构
```
wechatcloudhxx/
├── main.go                 # 应用入口
├── go.mod                  # Go模块文件
├── Dockerfile              # Docker构建文件
├── container.config.json   # 云托管配置
├── db/                     # 数据库层
│   ├── init.go            # 数据库初始化
│   ├── dao/               # 数据访问对象
│   └── model/             # 数据模型
├── service/               # 业务逻辑层
│   ├── auth_handler.go    # 认证处理器
│   ├── post_handler.go    # 帖子处理器
│   ├── user_handler.go    # 用户处理器
│   └── ...
├── docs/                  # 文档
└── sql/                   # 数据库脚本
```

## 🚀 快速开始

### 环境要求
- Go 1.16 或更高版本
- MySQL 5.7 或更高版本
- 微信开发者账号

### 本地开发

1. **克隆项目**
```bash
git clone <repository-url>
cd wechatcloudhxx
```

2. **安装依赖**
```bash
go mod download
```

3. **配置数据库**
```bash
# 创建数据库
mysql -u root -p
CREATE DATABASE wechat_community;

# 执行数据库脚本
mysql -u root -p wechat_community < sql/database_schema.sql
```

4. **配置环境变量**
```bash
# 创建 .env 文件
cp .env.example .env

# 编辑配置文件
vim .env
```

5. **运行项目**
```bash
go run main.go
```

### 微信云托管部署

1. **构建镜像**
```bash
docker build -t wechat-community .
```

2. **推送镜像**
```bash
# 登录微信云托管镜像仓库
docker login <registry-url>

# 推送镜像
docker push <registry-url>/wechat-community:latest
```

3. **配置云托管**
- 在微信开发者工具中创建云托管服务
- 配置容器规格和网络
- 设置环境变量和数据库连接

## 📚 API 文档

详细的API文档请参考：[API文档](docs/api_documentation.md)

### 主要接口

#### 用户认证
- `POST /api/auth/login` - 微信登录
- `GET /api/auth/userinfo` - 获取用户信息

#### 内容管理
- `GET /api/posts/` - 获取帖子列表
- `POST /api/posts/` - 发布帖子
- `GET /api/posts/{id}` - 获取帖子详情
- `DELETE /api/posts/{id}` - 删除帖子

#### 互动功能
- `POST /api/posts/{id}/like` - 点赞/取消点赞
- `POST /api/posts/{id}/comments` - 发表评论
- `GET /api/posts/{id}/comments` - 获取评论列表

#### 分类管理
- `GET /api/categories` - 获取分类列表
- `GET /api/topics/hot` - 获取热门话题

## 🔧 配置说明

### 数据库配置
```go
// db/init.go
type Config struct {
    Host     string
    Port     int
    User     string
    Password string
    Database string
}
```

### 微信配置
```go
// 微信小程序配置
AppID     string
AppSecret string
```

### 内容安全配置
```go
// 微信内容安全检测
ContentSecurityEnabled bool
ContentSecurityURL    string
```

## 🛡️ 安全特性

- **内容安全检测** - 集成微信内容安全API
- **用户认证** - 基于微信授权的用户认证
- **数据验证** - 完整的输入数据验证
- **SQL注入防护** - 使用GORM防止SQL注入
- **XSS防护** - 内容过滤和转义

## 📊 数据库设计

### 主要数据表
- `users` - 用户表
- `posts` - 帖子表
- `comments` - 评论表
- `user_likes` - 用户点赞表
- `categories` - 分类表
- `image_checks` - 图片检测记录表

详细的数据库设计请参考：[数据库设计](sql/database_schema.sql)

## 🚀 部署指南

### 微信云托管部署

1. **准备代码**
```bash
# 确保代码已提交到Git仓库
git add .
git commit -m "准备部署"
git push
```

2. **配置云托管**
- 在微信开发者工具中打开云托管控制台
- 创建新的服务实例
- 配置代码仓库和构建规则

3. **环境配置**
```json
{
  "env": {
    "DB_HOST": "your-db-host",
    "DB_PORT": "3306",
    "DB_USER": "your-db-user",
    "DB_PASSWORD": "your-db-password",
    "DB_DATABASE": "wechat_community",
    "WECHAT_APPID": "your-appid",
    "WECHAT_SECRET": "your-secret"
  }
}
```

4. **启动服务**
- 点击"部署"按钮
- 等待构建和部署完成
- 检查服务状态和日志

### 监控和日志

- **服务监控** - 通过微信云托管控制台监控服务状态
- **日志查看** - 在控制台查看实时日志
- **性能监控** - 监控CPU、内存、网络等指标

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 📞 联系方式

- 项目维护者：[amin]
- 邮箱：[nimilgg@qq.com]
- 微信：[aminzdx]

## 🙏 致谢

- [微信云托管](https://developers.weixin.qq.com/miniprogram/dev/wxcloudrun/) - 提供云托管服务
- [GORM](https://gorm.io/) - Go语言ORM框架
- [Go语言](https://golang.org/) - 编程语言

---

⭐ 如果这个项目对你有帮助，请给它一个星标！
