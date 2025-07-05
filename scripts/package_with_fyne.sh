#!/bin/bash

# 使用Fyne打包飞机大战游戏

# 颜色定义
RED="\033[0;31m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
BLUE="\033[0;34m"
NC="\033[0m" # No Color

# 检查Fyne CLI是否已安装
if ! command -v fyne &> /dev/null; then
    echo -e "${RED}错误: 未检测到Fyne CLI工具，请先安装:${NC}"
    echo -e "${YELLOW}go install fyne.io/fyne/v2/cmd/fyne@latest${NC}"
    exit 1
fi

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
        go mod tidy
    fi
fi

# 创建输出目录
mkdir -p dist

# 检测当前操作系统
OS="$(uname -s)"
case "${OS}" in
    Linux*)
        TARGET_OS="linux"
        ;;
    Darwin*)
        TARGET_OS="darwin"
        ;;
    CYGWIN*|MINGW*|MSYS*)
        TARGET_OS="windows"
        ;;
    *)
        echo -e "${RED}不支持的操作系统: ${OS}${NC}"
        exit 1
        ;;
esac

# 显示菜单
echo -e "${BLUE}=======================================${NC}"
echo -e "${BLUE}     使用Fyne打包飞机大战游戏        ${NC}"
echo -e "${BLUE}=======================================${NC}"
echo -e "当前系统: ${GREEN}${TARGET_OS}${NC}"
echo -e "\n请选择打包目标平台:"
echo -e "${YELLOW}1. Windows安装包${NC}"
echo -e "${YELLOW}2. macOS应用${NC}"
echo -e "${YELLOW}3. Linux安装包${NC}"
echo -e "${YELLOW}0. 退出${NC}"
echo -e "\n请输入选项 [0-3]: "
read -r choice

# 游戏图标
ICON_PATH="resources/images/player.png"
if [ ! -f "$ICON_PATH" ]; then
    echo -e "${YELLOW}警告: 未找到玩家图标，将使用默认图标${NC}"
    ICON_PATH=""
fi

# 根据选择打包
case $choice in
    0)
        echo -e "${BLUE}退出打包工具${NC}"
        exit 0
        ;;
    1)
        echo -e "\n${BLUE}开始为Windows打包...${NC}"
        if [ -n "$ICON_PATH" ]; then
            fyne package -os windows -icon "$ICON_PATH" -name "飞机大战" -appID "com.playgame.planewar"
        else
            fyne package -os windows -name "飞机大战" -appID "com.playgame.planewar"
        fi

        if [ $? -eq 0 ]; then
            # 移动到dist目录
            mv 飞机大战.exe dist/
            echo -e "${GREEN}Windows打包成功! 输出: dist/飞机大战.exe${NC}"
        else
            echo -e "${RED}Windows打包失败!${NC}"
            exit 1
        fi
        ;;
    2)
        echo -e "\n${BLUE}开始为macOS打包...${NC}"
        if [ -n "$ICON_PATH" ]; then
            fyne package -os darwin -icon "$ICON_PATH" -name "飞机大战" -appID "com.playgame.planewar"
        else
            fyne package -os darwin -name "飞机大战" -appID "com.playgame.planewar"
        fi

        if [ $? -eq 0 ]; then
            # 移动到dist目录
            mv 飞机大战.app dist/
            echo -e "${GREEN}macOS打包成功! 输出: dist/飞机大战.app${NC}"
        else
            echo -e "${RED}macOS打包失败!${NC}"
            exit 1
        fi
        ;;
    3)
        echo -e "\n${BLUE}开始为Linux打包...${NC}"
        if [ -n "$ICON_PATH" ]; then
            fyne package -os linux -icon "$ICON_PATH" -name "飞机大战" -appID "com.playgame.planewar"
        else
            fyne package -os linux -name "飞机大战" -appID "com.playgame.planewar"
        fi

        if [ $? -eq 0 ]; then
            # 移动到dist目录
            mv 飞机大战.tar.xz dist/
            echo -e "${GREEN}Linux打包成功! 输出: dist/飞机大战.tar.xz${NC}"
        else
            echo -e "${RED}Linux打包失败!${NC}"
            exit 1
        fi
        ;;
    *)
        echo -e "${RED}无效选项${NC}"
        exit 1
        ;;
esac

echo -e "\n${GREEN}打包过程完成!${NC}"
