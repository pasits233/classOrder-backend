#!/bin/bash

set -e

# 设置 Go Modules 代理
export GOPROXY=https://goproxy.cn,direct

# 当前目录即为 backend
# 创建 bin 目录
mkdir -p ./bin

# 编译
echo "Building backend..."
go mod tidy
go build -o ./bin/backend main.go

echo "Build finished. Binary is at: $(pwd)/bin/backend" 