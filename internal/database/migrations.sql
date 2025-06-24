-- 确保users表的id列是正确的类型
ALTER TABLE `users` MODIFY COLUMN `id` bigint unsigned NOT NULL AUTO_INCREMENT;

-- 处理coaches表的外键约束
SET @constraint_exists = (
    SELECT COUNT(*)
    FROM information_schema.TABLE_CONSTRAINTS 
    WHERE CONSTRAINT_SCHEMA = DATABASE()
    AND TABLE_NAME = 'coaches' 
    AND CONSTRAINT_NAME = 'coaches_ibfk_1'
);

SET @sql = IF(@constraint_exists > 0,
    'ALTER TABLE `coaches` DROP FOREIGN KEY `coaches_ibfk_1`',
    'SELECT "No constraint exists"');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 处理bookings表的外键约束
SET @booking_constraint_exists = (
    SELECT COUNT(*)
    FROM information_schema.TABLE_CONSTRAINTS 
    WHERE CONSTRAINT_SCHEMA = DATABASE()
    AND TABLE_NAME = 'bookings' 
    AND CONSTRAINT_NAME = 'bookings_ibfk_1'
);

SET @sql = IF(@booking_constraint_exists > 0,
    'ALTER TABLE `bookings` DROP FOREIGN KEY `bookings_ibfk_1`',
    'SELECT "No booking constraint exists"');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 修改coaches表的user_id列类型
ALTER TABLE `coaches` MODIFY COLUMN `user_id` bigint unsigned NOT NULL;

-- 修改coaches表的id列类型
ALTER TABLE `coaches` MODIFY COLUMN `id` bigint unsigned NOT NULL AUTO_INCREMENT;

-- 修改bookings表的coach_id列类型
ALTER TABLE `bookings` MODIFY COLUMN `coach_id` bigint unsigned NOT NULL;

-- 重新添加外键约束
ALTER TABLE `coaches` ADD CONSTRAINT `coaches_ibfk_1` 
FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

ALTER TABLE `bookings` ADD CONSTRAINT `bookings_ibfk_1`
FOREIGN KEY (`coach_id`) REFERENCES `coaches` (`id`) ON DELETE CASCADE; 