-- 清空所有相关表
DELETE FROM bookings;
DELETE FROM coaches;
DELETE FROM users;

-- 创建管理员用户（密码：admin123）
INSERT INTO users (username, password_hash, role) VALUES 
('admin', '$2a$10$QOJ6qh/ZhBfGNqHlPVpE8.dw1j6qJJwNp7U8P6Ae.S0UZ3TyXqmjG', 'admin');

-- 创建4个教练用户（密码都是：coach123）
INSERT INTO users (id, username, password_hash, role) VALUES
(2, 'vincy', '$2a$10$8jw70j7J0JJvHzGxGYq8b.f/BmNqFA8UYqX3G8JmGG6R7EVhCzKJO', 'coach'),
(3, 'jj', '$2a$10$8jw70j7J0JJvHzGxGYq8b.f/BmNqFA8UYqX3G8JmGG6R7EVhCzKJO', 'coach'),
(4, 'pentium', '$2a$10$8jw70j7J0JJvHzGxGYq8b.f/BmNqFA8UYqX3G8JmGG6R7EVhCzKJO', 'coach'),
(5, 'liming', '$2a$10$8jw70j7J0JJvHzGxGYq8b.f/BmNqFA8UYqX3G8JmGG6R7EVhCzKJO', 'coach');

-- 创建教练信息
INSERT INTO coaches (id, user_id, name, description, avatar_url) VALUES
(1, 2, '陈万鑫 Vincy', '奥地利二级，PSIE三级（滑行），国职5级。专业滑雪教练，拥有丰富的教学经验，擅长单板滑雪教学。', '/assets/coach1.jpg'),
(2, 3, '钟金君 JJ', '奥地利二级，国职5级。资深滑雪教练，专注于双板滑雪教学，擅长初学者指导。', '/assets/coach1.jpg'),
(3, 4, '陈化益 Pentium', '奥地利一级，加拿大三级，国职5级。国际认证滑雪教练，擅长高级技巧训练和竞技滑雪指导。', '/assets/coach1.jpg'),
(4, 5, '李明', '加拿大一级，国职4级。专业滑雪教练，擅长儿童滑雪教学和家庭滑雪指导。', '/assets/coach1.jpg'); 