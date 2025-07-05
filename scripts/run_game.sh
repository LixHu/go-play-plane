#!/bin/bash

# 颜色定义
RED="\033[0;31m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
NC="\033[0m" # No Color

echo -e "${GREEN}正在启动飞机大战游戏...${NC}"

# 检查purego版本
PUREGO_VERSION=$(go list -m github.com/ebitengine/purego 2>/dev/null | awk '{print $2}')

if [ "$PUREGO_VERSION" != "v0.8.3" ]; then
    echo -e "${YELLOW}警告: 检测到purego版本 $PUREGO_VERSION，可能会出现符号重复定义问题${NC}"
    echo -e "${YELLOW}建议先运行 scripts/fix_purego.sh 修复依赖${NC}"
    echo
    read -p "是否继续? (y/n): " choice
    if [ "$choice" != "y" ] && [ "$choice" != "Y" ]; then
        exit 0
    fi
fi

# 运行游戏
go run .

if [ $? -ne 0 ]; then
    echo
    echo -e "${RED}游戏启动失败，可能是依赖问题${NC}"
    echo -e "${YELLOW}尝试运行 scripts/fix_purego.sh 修复依赖后再试${NC}"
fi
