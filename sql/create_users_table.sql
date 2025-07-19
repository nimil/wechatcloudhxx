-- 创建用户表
CREATE TABLE IF NOT EXISTS `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码',
  `openid` varchar(100) DEFAULT NULL COMMENT '微信OpenID',
  `unionid` varchar(100) DEFAULT NULL COMMENT '微信UnionID',
  `appid` varchar(50) DEFAULT NULL COMMENT '微信AppID',
  `createdAt` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatedAt` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_openid` (`openid`),
  UNIQUE KEY `uk_unionid` (`unionid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 创建索引
CREATE INDEX `idx_username` ON `users` (`username`);
CREATE INDEX `idx_openid` ON `users` (`openid`);
CREATE INDEX `idx_unionid` ON `users` (`unionid`); 