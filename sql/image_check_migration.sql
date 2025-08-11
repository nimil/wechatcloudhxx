-- 图片检测功能数据库迁移

-- 1. 为posts表添加图片检测状态字段
ALTER TABLE posts ADD COLUMN image_check_status INT DEFAULT 0 COMMENT '图片检测状态：0-待检测 1-检测中 2-检测通过 3-检测失败';

-- 2. 创建图片检测记录表
CREATE TABLE image_checks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    post_id BIGINT NOT NULL COMMENT '关联的帖子ID',
    image_url VARCHAR(500) NOT NULL COMMENT '图片URL',
    trace_id VARCHAR(100) NOT NULL COMMENT '微信检测追踪ID',
    status INT DEFAULT 0 COMMENT '检测状态：0-待检测 1-检测中 2-检测通过 3-检测失败',
    suggest VARCHAR(20) COMMENT '检测建议：pass/review/risky',
    label INT DEFAULT 0 COMMENT '检测标签',
    prob DECIMAL(5,2) COMMENT '置信度',
    strategy VARCHAR(50) COMMENT '检测策略',
    errcode INT DEFAULT 0 COMMENT '错误码',
    errmsg VARCHAR(200) COMMENT '错误信息',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_post_id (post_id),
    INDEX idx_trace_id (trace_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='图片检测记录表';

-- 3. 为现有帖子设置默认状态（没有图片的帖子设为检测通过）
UPDATE posts SET image_check_status = 2 WHERE images = '[]' OR images = '' OR images IS NULL;

-- 4. 为有图片的帖子设置待检测状态
UPDATE posts SET image_check_status = 0 WHERE images != '[]' AND images != '' AND images IS NOT NULL;
