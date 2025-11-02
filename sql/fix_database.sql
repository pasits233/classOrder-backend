-- ============================================
-- 数据库修复脚本
-- 1. 创建 wechat_users 表
-- 2. 为 bookings 表添加 student_id 字段
-- ============================================

-- 1. 创建微信用户表
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
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间（软删除）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_open_id` (`open_id`),
  KEY `idx_union_id` (`union_id`),
  KEY `idx_role` (`role`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='微信用户表';

-- 2. 为 bookings 表添加 student_id 字段
-- 注意：如果字段已存在会报错，可以忽略或先检查
ALTER TABLE `bookings` 
ADD COLUMN `student_id` int(10) unsigned DEFAULT NULL COMMENT '学员ID（微信用户ID）' AFTER `coach_id`;

-- 3. 为 student_id 添加索引
ALTER TABLE `bookings` 
ADD INDEX `idx_student_id` (`student_id`);

