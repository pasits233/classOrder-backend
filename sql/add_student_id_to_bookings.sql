-- 为 bookings 表添加 student_id 字段
ALTER TABLE `bookings` 
ADD COLUMN `student_id` int(11) unsigned DEFAULT NULL COMMENT '学员ID（微信用户ID）' AFTER `coach_id`,
ADD INDEX `idx_student_id` (`student_id`);

-- 更新现有预约记录的 student_id（如果有需要的话）
-- UPDATE `bookings` SET `student_id` = 1 WHERE `student_id` IS NULL; 