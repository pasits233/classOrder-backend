-- 确保users表的id列是正确的类型
ALTER TABLE `users` MODIFY COLUMN `id` bigint unsigned NOT NULL AUTO_INCREMENT;

-- 删除现有的外键约束（如果存在）
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

-- 修改coaches表的user_id列类型
ALTER TABLE `coaches` MODIFY COLUMN `user_id` bigint unsigned NOT NULL;

-- 添加新的外键约束
ALTER TABLE `coaches` ADD CONSTRAINT `coaches_ibfk_1` 
FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE; 