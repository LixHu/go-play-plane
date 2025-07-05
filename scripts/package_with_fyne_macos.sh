#!/bin/bash

# macOS专用Fyne打包脚本

# 颜色定义
RED="\033[0;31m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
BLUE="\033[0;34m"
NC="\033[0m" # No Color

# 检查是否在macOS上运行
if [ "$(uname)" != "Darwin" ]; then
    echo -e "${RED}错误: 此脚本只能在macOS上运行${NC}"
    exit 1
fi

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

# 检查Xcode命令行工具
if ! command -v xcode-select &> /dev/null; then
    echo -e "${RED}错误: 未检测到Xcode命令行工具${NC}"
    echo -e "${YELLOW}请运行: xcode-select --install${NC}"
    exit 1
fi

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

# 游戏图标
ICON_PATH="resources/images/player.png"
if [ ! -f "$ICON_PATH" ]; then
    echo -e "${YELLOW}警告: 未找到玩家图标，将使用默认图标${NC}"
    ICON_PATH=""
fi

echo -e "\n${BLUE}开始为macOS打包...${NC}"

# 确保CGO启用
export CGO_ENABLED=1

if [ -n "$ICON_PATH" ]; then
    fyne package -os darwin -icon "$ICON_PATH" -name "飞机大战" -appID "com.playgame.planewar"
else
    fyne package -os darwin -name "飞机大战" -appID "com.playgame.planewar"
fi

if [ $? -eq 0 ]; then
    # 移动到dist目录
    mv 飞机大战.app dist/
    echo -e "${GREEN}macOS打包成功! 输出: dist/飞机大战.app${NC}"

    # 是否需要签名
    echo -e "\n${YELLOW}注意: 未签名的应用在macOS上可能会被拦截${NC}"
    read -p "是否要对应用进行签名? (y/n): " sign_choice
    if [ "$sign_choice" = "y" ] || [ "$sign_choice" = "Y" ]; then
        # 列出可用的签名证书
        echo -e "\n${BLUE}可用的开发者证书:${NC}"
        security find-identity -v -p codesigning

        read -p "请输入要使用的证书名称 (例如 'Developer ID Application: Your Name'): " cert_name

        if [ -n "$cert_name" ]; then
            echo -e "\n${BLUE}正在签名应用...${NC}"
            codesign --force --deep --sign "$cert_name" dist/飞机大战.app

            if [ $? -eq 0 ]; then
                echo -e "${GREEN}应用签名成功!${NC}"
            else
                echo -e "${RED}应用签名失败!${NC}"
            fi
        fi
    fi
else
    echo -e "${RED}macOS打包失败!${NC}"
    exit 1
fi

echo -e "\n${GREEN}打包过程完成!${NC}"
