-- Drop the existing foreign key constraint
ALTER TABLE `coaches` DROP FOREIGN KEY `coaches_ibfk_1`;

-- Modify the column
ALTER TABLE `coaches` MODIFY COLUMN `user_id` bigint unsigned NOT NULL;

-- Re-add the foreign key constraint
ALTER TABLE `coaches` ADD CONSTRAINT `coaches_ibfk_1` 
FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE; 