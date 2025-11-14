-- 创建场地表
CREATE TABLE IF NOT EXISTS `venues` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL UNIQUE,
    `address` VARCHAR(255),
    `manager_name` VARCHAR(255),
    `contact` VARCHAR(100),
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 确保存在一个默认场地，方便将历史数据回填
INSERT INTO `venues` (`id`, `name`, `address`, `manager_name`, `contact`, `created_at`, `updated_at`)
VALUES (1, '默认场地', '', '', '', NOW(), NOW())
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- 为 bookings 表新增场地列，默认指向默认场地
SET @venue_column_exists := (
    SELECT COUNT(*)
    FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'bookings'
      AND COLUMN_NAME = 'venue_id'
);

SET @add_venue_column_sql := IF(@venue_column_exists = 0,
    'ALTER TABLE `bookings` ADD COLUMN `venue_id` BIGINT UNSIGNED NOT NULL DEFAULT 1 AFTER `coach_id`',
    'SELECT "venue_id already exists"');
PREPARE stmt FROM @add_venue_column_sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 将历史数据的场地ID设置为默认值
UPDATE `bookings` SET `venue_id` = 1 WHERE `venue_id` IS NULL OR `venue_id` = 0;

-- 删除旧的教练唯一索引
SET @old_index_exists := (
    SELECT COUNT(*)
    FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'bookings'
      AND INDEX_NAME = 'uniq_coach_date_slot'
);

SET @drop_old_index_sql := IF(@old_index_exists > 0,
    'ALTER TABLE `bookings` DROP INDEX `uniq_coach_date_slot`',
    'SELECT "old index not found"');
PREPARE drop_stmt FROM @drop_old_index_sql;
EXECUTE drop_stmt;
DEALLOCATE PREPARE drop_stmt;

-- 创建新的基于场地的唯一索引
SET @new_index_exists := (
    SELECT COUNT(*)
    FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'bookings'
      AND INDEX_NAME = 'uniq_venue_date_slot'
);

SET @create_new_index_sql := IF(@new_index_exists = 0,
    'ALTER TABLE `bookings` ADD UNIQUE INDEX `uniq_venue_date_slot` (`venue_id`, `booking_date`, `time_slot`)',
    'SELECT "new index already exists"');
PREPARE create_stmt FROM @create_new_index_sql;
EXECUTE create_stmt;
DEALLOCATE PREPARE create_stmt;

-- 创建教练视频表
CREATE TABLE IF NOT EXISTS `coach_videos` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `coach_id` BIGINT UNSIGNED NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `video_url` VARCHAR(255) NOT NULL,
    `cover_url` VARCHAR(255),
    `sort_order` INT DEFAULT 0,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_coach_id` (`coach_id`),
    CONSTRAINT `fk_coach_videos_coach` FOREIGN KEY (`coach_id`) REFERENCES `coaches`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

