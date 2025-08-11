# 内容安全校验功能集成说明

## 概述

本项目已集成微信小程序内容安全检测功能，在用户发布帖子和评论时会自动进行内容安全校验，确保内容符合平台规范。

## 功能特性

- **自动检测**：发布帖子和评论时自动进行内容安全检测
- **场景适配**：根据内容类型使用不同的检测场景
- **错误提示**：检测到违规内容时提供友好的错误提示
- **无感知集成**：对现有业务逻辑无影响

## 场景值配置

根据微信小程序内容安全检测API，使用以下场景值：

| 场景值 | 场景名称 | 使用场景 |
|--------|----------|----------|
| 1 | 资料 | 用户资料、个人信息等 |
| 2 | 评论 | 评论、回复等 |
| 3 | 论坛 | 帖子、动态等 |
| 4 | 社交日志 | 社交分享、日志等 |

## 集成位置

### 1. 帖子发布

**文件位置：** `service/post_service.go`

**检测内容：**
- 帖子标题
- 帖子正文内容

**使用场景：** 论坛（场景值：3）

**检测逻辑：**
```go
// 检查标题安全性（使用论坛场景）
if req.Title != "" {
    isSafe, err := s.securityService.IsContentSafe(openid, req.Title, SceneForum)
    if err != nil {
        return nil, fmt.Errorf("标题安全检测失败: %v", err)
    }
    if !isSafe {
        return nil, fmt.Errorf("标题包含违规内容，请修改后重试")
    }
}

// 检查内容安全性（使用论坛场景）
if req.Content != "" {
    isSafe, err := s.securityService.IsContentSafe(openid, req.Content, SceneForum)
    if err != nil {
        return nil, fmt.Errorf("内容安全检测失败: %v", err)
    }
    if !isSafe {
        return nil, fmt.Errorf("内容包含违规信息，请修改后重试")
    }
}
```

### 2. 评论发布

**文件位置：** `service/comment_service.go`

**检测内容：**
- 评论内容

**使用场景：** 评论（场景值：2）

**检测逻辑：**
```go
// 内容安全校验
if openid != "" && req.Content != "" {
    isSafe, err := s.securityService.IsContentSafe(openid, req.Content, SceneComment)
    if err != nil {
        return nil, fmt.Errorf("内容安全检测失败: %v", err)
    }
    if !isSafe {
        return nil, fmt.Errorf("评论内容包含违规信息，请修改后重试")
    }
}
```

## 错误处理

当检测到违规内容时，系统会返回相应的错误信息：

- **标题违规**：`标题包含违规内容，请修改后重试`
- **内容违规**：`内容包含违规信息，请修改后重试`
- **评论违规**：`评论内容包含违规信息，请修改后重试`
- **检测失败**：`内容安全检测失败: [具体错误信息]`

## 依赖要求

### 1. OpenID获取

内容安全检测需要用户的OpenID，系统会从请求头中获取：

```go
// 从请求头获取openid
openid := r.Header.Get("x-wx-openid")
```

### 2. 微信小程序环境

需要在微信小程序环境中运行，确保能够获取到用户的OpenID。

## 配置说明

### 1. 场景值常量

在 `service/content_security_service.go` 中定义了场景值常量：

```go
const (
    SceneProfile = 1 // 资料
    SceneComment = 2 // 评论
    SceneForum   = 3 // 论坛
    SceneSocial  = 4 // 社交日志
)
```

### 2. 检测服务

内容安全检测服务位于 `service/content_security_service.go`，主要方法：

- `CheckContentSecurity(openid, content, scene)` - 检查内容安全性
- `IsContentSafe(openid, content, scene)` - 判断内容是否安全
- `GetContentSecurityResult(openid, content, scene)` - 获取详细检测结果

## 使用示例

### 前端调用

```javascript
// 发布帖子
const postData = {
    title: "帖子标题",
    content: "帖子内容",
    category: "tech",
    tags: ["技术", "分享"],
    isPublic: true
};

const response = await fetch('/api/posts', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
        'x-wx-openid': '用户的openid'
    },
    body: JSON.stringify(postData)
});

const result = await response.json();
if (result.code !== 200) {
    console.error('发布失败:', result.message);
}
```

### 错误处理

```javascript
try {
    const response = await fetch('/api/posts', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'x-wx-openid': '用户的openid'
        },
        body: JSON.stringify(postData)
    });
    
    const result = await response.json();
    if (result.code === 200) {
        console.log('发布成功');
    } else {
        // 处理内容安全检测失败
        if (result.message.includes('违规')) {
            alert('内容包含违规信息，请修改后重试');
        } else {
            alert('发布失败: ' + result.message);
        }
    }
} catch (error) {
    console.error('请求失败:', error);
}
```

## 注意事项

1. **OpenID必需**：内容安全检测需要用户的OpenID，确保请求头中包含 `x-wx-openid`
2. **网络依赖**：检测需要调用微信API，确保网络连接正常
3. **超时处理**：检测服务设置了10秒超时，避免长时间等待
4. **错误降级**：如果检测服务不可用，可以考虑降级处理（当前版本会直接返回错误）
5. **性能考虑**：检测会增加响应时间，建议在UI上提供加载提示

## 扩展建议

1. **缓存机制**：对相同内容可以添加缓存，避免重复检测
2. **异步检测**：可以考虑异步检测，不阻塞用户操作
3. **降级策略**：当检测服务不可用时，可以采取降级策略
4. **日志记录**：记录检测结果用于后续分析和优化
5. **自定义规则**：可以结合业务需求添加自定义的内容检测规则
