-- 纯净社区小程序数据库表结构
-- 数据库名称: golang_demo
-- 字符集: utf8mb4
-- 排序规则: utf8mb4_unicode_ci

-- 用户表
CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `nickname` varchar(50) DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(500) DEFAULT NULL COMMENT '头像URL',
  `bio` varchar(200) DEFAULT NULL COMMENT '个人简介',
  `level` int DEFAULT 1 COMMENT '用户等级',
  `is_verified` tinyint(1) DEFAULT 0 COMMENT '是否认证用户',
  `password` varchar(255) NOT NULL COMMENT '密码',
  `openid` varchar(100) DEFAULT NULL COMMENT '微信OpenID',
  `unionid` varchar(100) DEFAULT NULL COMMENT '微信UnionID',
  `appid` varchar(100) DEFAULT NULL COMMENT '应用ID',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_openid` (`openid`),
  KEY `idx_unionid` (`unionid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 分类表
CREATE TABLE `categories` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `name` varchar(50) NOT NULL COMMENT '分类名称',
  `code` varchar(20) NOT NULL COMMENT '分类代码',
  `icon` varchar(50) DEFAULT NULL COMMENT '分类图标',
  `description` varchar(200) DEFAULT NULL COMMENT '分类描述',
  `post_count` int DEFAULT 0 COMMENT '帖子数量',
  `sort` int DEFAULT 0 COMMENT '排序',
  `is_active` tinyint(1) DEFAULT 1 COMMENT '是否启用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分类表';

-- 帖子表
CREATE TABLE `posts` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '帖子ID',
  `title` varchar(100) NOT NULL COMMENT '帖子标题',
  `content` text NOT NULL COMMENT '帖子内容',
  `excerpt` varchar(500) DEFAULT NULL COMMENT '帖子摘要',
  `author_id` bigint NOT NULL COMMENT '作者ID',
  `category` varchar(20) NOT NULL COMMENT '分类代码',
  `category_name` varchar(50) NOT NULL COMMENT '分类名称',
  `tags` text DEFAULT NULL COMMENT '标签(JSON格式)',
  `images` text DEFAULT NULL COMMENT '图片URL列表(JSON格式)',
  `is_public` tinyint(1) DEFAULT 1 COMMENT '是否公开',
  `is_deleted` tinyint(1) DEFAULT 0 COMMENT '是否已删除',
  `likes` int DEFAULT 0 COMMENT '点赞数',
  `comments` int DEFAULT 0 COMMENT '评论数',
  `views` int DEFAULT 0 COMMENT '浏览量',
  `shares` int DEFAULT 0 COMMENT '分享数',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_author_id` (`author_id`),
  KEY `idx_category` (`category`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_likes` (`likes`),
  KEY `idx_views` (`views`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='帖子表';

-- 评论表
CREATE TABLE `comments` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '评论ID',
  `content` varchar(500) NOT NULL COMMENT '评论内容',
  `author_id` bigint NOT NULL COMMENT '评论者ID',
  `post_id` bigint NOT NULL COMMENT '帖子ID',
  `parent_id` bigint DEFAULT NULL COMMENT '父评论ID',
  `likes` int DEFAULT 0 COMMENT '点赞数',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_author_id` (`author_id`),
  KEY `idx_post_id` (`post_id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论表';

-- 用户点赞关系表
CREATE TABLE `user_likes` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `post_id` bigint NOT NULL COMMENT '帖子ID',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_post` (`user_id`,`post_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_post_id` (`post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户点赞关系表';

-- 初始化默认分类数据
INSERT INTO `categories` (`id`, `name`, `code`, `icon`, `description`, `sort`) VALUES
(1, '闲置', 'idle', '📦', '闲置物品交易、二手市场、物品分享', 1),
(2, '美食', 'food', '🍜', '美食制作、餐厅推荐、食谱分享', 2),
(3, '生活', 'life', '🏠', '日常生活、家居装饰、生活技巧', 3),
(4, '交流', 'chat', '💬', '社区交流、话题讨论、经验分享', 4),
(5, '其他', 'other', '📝', '其他话题、杂谈、综合内容', 5);

-- 创建索引优化查询性能
CREATE INDEX idx_posts_category_created_at ON posts(category, created_at DESC);
CREATE INDEX idx_posts_likes_created_at ON posts(likes DESC, created_at DESC);
CREATE INDEX idx_comments_post_created_at ON comments(post_id, created_at DESC); 