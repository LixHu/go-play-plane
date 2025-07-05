#!/bin/bash

# Windows专用构建脚本

# 颜色定义
RED="\033[0;31m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
BLUE="\033[0;34m"
NC="\033[0m" # No Color

# 游戏名称
GAME_NAME="飞机大战"

# 创建输出目录
mkdir -p dist

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo -e "${RED}错误: 未检测到Go环境，请安装Go后再运行此脚本。${NC}"
    exit 1
fi

GO_VERSION=$(go version | grep -oP "go\d+\.\d+" | cut -c 3-)
echo -e "${GREEN}检测到Go版本: ${GO_VERSION}${NC}"

# 检查purego版本
PUREGO_VERSION=$(go list -m github.com/ebitengine/purego | awk '{print $2}')
if [ "$PUREGO_VERSION" != "v0.8.3" ]; then
    echo -e "${YELLOW}检测到purego版本 $PUREGO_VERSION，建议更新到v0.8.3以解决符号重复问题${NC}"
    read -p "是否更新purego到v0.8.3? (y/n): " choice
    if [ "$choice" = "y" ] || [ "$choice" = "Y" ]; then
        go get github.com/ebitengine/purego@v0.8.3
        echo -e "${GREEN}已更新purego到v0.8.3${NC}"
    fi
fi

# 检查依赖
echo -e "${YELLOW}正在检查依赖...${NC}"
go mod tidy
if [ $? -ne 0 ]; then
    echo -e "${RED}依赖安装失败${NC}"
    exit 1
fi

# 确保依赖完整下载
go mod download
if [ $? -ne 0 ]; then
    echo -e "${RED}依赖下载失败${NC}"
    exit 1
fi

# 开始构建Windows版本
echo -e "\n${BLUE}开始构建Windows版本...${NC}"

# 设置环境变量
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=1

# 编译
go build -o "dist/${GAME_NAME}.exe" -ldflags="-H windowsgui" .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}Windows版本构建成功! 输出: dist/${GAME_NAME}.exe${NC}"
else
    echo -e "${RED}Windows版本构建失败!${NC}"
    exit 1
fi
