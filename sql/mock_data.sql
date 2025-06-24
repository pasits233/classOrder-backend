-- 清空表（开发环境使用，生产环境请谨慎）
DELETE FROM coaches;
DELETE FROM users;
DELETE FROM courses;

-- 用户表 mock 数据（密码为明文 password，实际应为加密哈希）
INSERT INTO users (id, username, password_hash, role) VALUES
(1, 'vincy', 'password', 'coach'),
(2, 'jj', 'password', 'coach'),
(3, 'pentium', 'password', 'coach'),
(4, 'liming', 'password', 'coach');

-- 教练表 mock 数据
INSERT INTO coaches (id, user_id, name, description, avatar_url) VALUES
(1, 1, '陈万鑫 Vincy', '奥地利二级，PSIE三级（滑行），国职5级。专业滑雪教练，拥有丰富的教学经验，擅长单板滑雪教学。', '/assets/coach1.jpg'),
(2, 2, '钟金君 JJ', '奥地利二级，国职5级。资深滑雪教练，专注于双板滑雪教学，擅长初学者指导。', '/assets/coach1.jpg'),
(3, 3, '陈化益 Pentium', '奥地利一级，加拿大三级，国职5级。国际认证滑雪教练，擅长高级技巧训练和竞技滑雪指导。', '/assets/coach1.jpg'),
(4, 4, '李明', '加拿大一级，国职4级。专业滑雪教练，擅长儿童滑雪教学和家庭滑雪指导。', '/assets/coach1.jpg');

-- 课程表 mock 数据
INSERT INTO courses (id, name, description, price) VALUES
(1, '奥地利大神营', '开启新雪季进阶之旅', 3388),
(2, '私教课', '1V1-2高效进阶必选课', 1800); 