#!/bin/bash

# 确保脚本在错误时停止
set -e

echo "开始构建前端项目..."

# 检查是否安装了 pnpm
if ! command -v pnpm &> /dev/null; then
    echo "正在安装 pnpm..."
    npm install -g pnpm
fi

# 安装依赖
echo "安装项目依赖..."
pnpm install

# 构建项目
echo "构建项目..."
pnpm run build

echo "构建完成！"
echo "dist 目录已生成，现在可以重启后端服务了。" 