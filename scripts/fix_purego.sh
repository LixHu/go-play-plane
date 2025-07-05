#!/bin/bash

# 修复purego dlopen符号重复定义问题的脚本
echo "开始修复purego依赖问题..."

# 清理go模块缓存中的purego
echo "清理go模块缓存..."
go clean -modcache

# 更新依赖
echo "更新依赖到purego v0.8.3..."
go get github.com/ebitengine/purego@v0.8.3

# 整理go.mod
echo "整理go.mod文件..."
go mod tidy

# 验证依赖版本
PUREGO_VERSION=$(go list -m github.com/ebitengine/purego | awk '{print $2}')
echo "当前purego版本: $PUREGO_VERSION"

if [ "$PUREGO_VERSION" = "v0.8.3" ]; then
    echo "✅ 成功更新到purego v0.8.3"
    echo "现在可以尝试重新构建项目"
else
    echo "❌ purego更新失败，请手动编辑go.mod文件"
    echo "将github.com/ebitengine/purego版本改为v0.8.3"
fi
