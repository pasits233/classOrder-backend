# 微信小程序后端部署指南

## 概述
本文档介绍如何部署支持微信小程序的后端服务，域名：`snowowski.site`

## 环境要求
- Go 1.19+
- MySQL 8.0+
- 已备案的域名（snowowski.site）
- HTTPS 证书

## 1. 环境变量配置

创建 `.env` 文件（基于 `config.env.example`）：

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_actual_password
DB_NAME=classorder

# JWT 配置
JWT_SECRET=your_actual_jwt_secret_key

# 微信小程序配置
WECHAT_APP_ID=your_actual_wechat_app_id
WECHAT_APP_SECRET=your_actual_wechat_app_secret

# 服务器配置
SERVER_PORT=9528
SERVER_MODE=release

# 域名配置
DOMAIN=snowowski.site
```

## 2. 数据库迁移

### 2.1 创建微信用户表
```sql
-- 执行 backend/sql/wechat_users.sql
source backend/sql/wechat_users.sql;
```

### 2.2 为预约表添加学员ID字段
```sql
-- 执行 backend/sql/add_student_id_to_bookings.sql
source backend/sql/add_student_id_to_bookings.sql;
```

## 3. 微信小程序配置

### 3.1 获取 AppID 和 AppSecret
1. 登录微信公众平台：https://mp.weixin.qq.com/
2. 选择你的小程序
3. 进入"开发" -> "开发管理" -> "开发设置"
4. 复制 AppID 和 AppSecret

### 3.2 配置服务器域名
在微信公众平台配置以下域名：
- request 合法域名：`https://snowowski.site`
- uploadFile 合法域名：`https://snowowski.site`
- downloadFile 合法域名：`https://snowowski.site`

## 4. 编译和部署

### 4.1 编译后端
```bash
cd backend
go mod tidy
go build -o classorder-backend main.go
```

### 4.2 启动服务
```bash
# 开发环境
go run main.go

# 生产环境
./classorder-backend
```

### 4.3 使用 systemd 管理服务（推荐）
```bash
# 创建服务文件
sudo nano /etc/systemd/system/classorder-backend.service

[Unit]
Description=ClassOrder Backend Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/your/backend
ExecStart=/path/to/your/backend/classorder-backend
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target

# 启用服务
sudo systemctl enable classorder-backend
sudo systemctl start classorder-backend
sudo systemctl status classorder-backend
```

## 5. Nginx 配置

确保 Nginx 配置正确转发 API 请求：

```nginx
# 微信小程序 API 路由
location /api/wechat/ {
    proxy_pass http://127.0.0.1:9528;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}

# 其他 API 路由
location /api/ {
    proxy_pass http://127.0.0.1:9528;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}
```

## 6. 测试接口

### 6.1 微信登录测试
```bash
curl -X POST https://snowowski.site/api/wechat/login \
  -H "Content-Type: application/json" \
  -d '{"code":"test_code"}'
```

### 6.2 获取教练列表测试
```bash
curl https://snowowski.site/api/coaches
```

## 7. 常见问题

### 7.1 微信登录失败
- 检查 AppID 和 AppSecret 是否正确
- 确认域名是否已在微信公众平台配置
- 检查后端日志

### 7.2 数据库连接失败
- 检查数据库服务是否运行
- 确认数据库连接参数
- 检查防火墙设置

### 7.3 接口 404 错误
- 确认后端服务是否启动
- 检查 Nginx 配置是否正确
- 验证路由是否正确注册

## 8. 监控和日志

### 8.1 查看服务状态
```bash
sudo systemctl status classorder-backend
```

### 8.2 查看日志
```bash
# 实时日志
sudo journalctl -u classorder-backend -f

# 查看最近的日志
sudo journalctl -u classorder-backend -n 100
```

## 9. 安全建议

1. **JWT Secret**：使用强随机字符串
2. **数据库密码**：使用强密码
3. **HTTPS**：必须启用 HTTPS
4. **防火墙**：只开放必要端口
5. **定期备份**：定期备份数据库

## 10. 联系支持

如遇到问题，请检查：
1. 后端服务日志
2. Nginx 错误日志
3. 数据库连接状态
4. 微信公众平台配置 