-- 创建微信用户表
CREATE TABLE IF NOT EXISTS `wechat_users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `open_id` varchar(100) NOT NULL COMMENT '微信OpenID',
  `union_id` varchar(100) DEFAULT NULL COMMENT '微信UnionID',
  `nick_name` varchar(100) DEFAULT NULL COMMENT '微信昵称',
  `avatar_url` text DEFAULT NULL COMMENT '微信头像URL',
  `gender` tinyint(1) DEFAULT 0 COMMENT '性别 0-未知 1-男 2-女',
  `country` varchar(50) DEFAULT NULL COMMENT '国家',
  `province` varchar(50) DEFAULT NULL COMMENT '省份',
  `city` varchar(50) DEFAULT NULL COMMENT '城市',
  `language` varchar(20) DEFAULT NULL COMMENT '语言',
  `role` varchar(20) DEFAULT 'student' COMMENT '用户角色 student-学员',
  `status` tinyint(1) DEFAULT 1 COMMENT '状态 0-禁用 1-正常',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_open_id` (`open_id`),
  KEY `idx_union_id` (`union_id`),
  KEY `idx_role` (`role`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='微信用户表';

-- 插入测试数据（可选）
-- INSERT INTO `wechat_users` (`open_id`, `union_id`, `nick_name`, `avatar_url`, `gender`, `country`, `province`, `city`, `language`, `role`, `status`) VALUES
-- ('test_openid_001', 'test_unionid_001', '测试用户1', 'https://example.com/avatar1.jpg', 1, '中国', '广东省', '深圳市', 'zh_CN', 'student', 1),
-- ('test_openid_002', 'test_unionid_002', '测试用户2', 'https://example.com/avatar2.jpg', 2, '中国', '北京市', '北京市', 'zh_CN', 'student', 1); 