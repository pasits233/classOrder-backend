#!/bin/bash
# 检查并修复 uploads 目录

BACKEND_DIR="/www/wwwroot/classOrder-backend"
UPLOADS_DIR="$BACKEND_DIR/uploads"

echo "检查 uploads 目录..."

# 检查后端目录是否存在
if [ ! -d "$BACKEND_DIR" ]; then
    echo "错误: 后端目录不存在: $BACKEND_DIR"
    exit 1
fi

# 创建 uploads 目录（如果不存在）
if [ ! -d "$UPLOADS_DIR" ]; then
    echo "创建 uploads 目录: $UPLOADS_DIR"
    mkdir -p "$UPLOADS_DIR"
fi

# 设置正确的权限
chmod 755 "$UPLOADS_DIR"
echo "设置 uploads 目录权限为 755"

# 检查目录中的文件
echo ""
echo "uploads 目录内容:"
ls -lh "$UPLOADS_DIR"

echo ""
echo "检查完成！"

