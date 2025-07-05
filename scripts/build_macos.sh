#!/bin/bash

# macOS专用构建脚本

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

# 检查 macOS 构建环境
echo -e "${YELLOW}检查 macOS 构建环境...${NC}"

# 检查 Xcode 命令行工具
if ! command -v clang &> /dev/null; then
    echo -e "${RED}错误: 未检测到 Xcode 命令行工具，请运行: xcode-select --install${NC}"
    exit 1
fi

# 检查 CGO 环境
if [ "$(go env CGO_ENABLED)" != "1" ]; then
    echo -e "${YELLOW}警告: CGO 未启用，正在启用...${NC}"
    export CGO_ENABLED=1
fi

echo -e "${GREEN}macOS 构建环境检查完成${NC}"

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

# 开始构建macOS版本
echo -e "\n${BLUE}开始构建macOS版本...${NC}"

# 设置环境变量
export GOOS=darwin
export GOARCH=amd64
export CGO_ENABLED=1

# 编译
go build -o "dist/${GAME_NAME}_macos" .

if [ $? -ne 0 ]; then
    echo -e "${RED}macOS版本构建失败!${NC}"
    exit 1
fi

# 创建.app包结构
APP_DIR="dist/${GAME_NAME}.app"
CONTENTS_DIR="${APP_DIR}/Contents"
MACOS_DIR="${CONTENTS_DIR}/MacOS"
RESOURCES_DIR="${CONTENTS_DIR}/Resources"

# 创建目录结构
mkdir -p "${MACOS_DIR}"
mkdir -p "${RESOURCES_DIR}"

# 复制可执行文件
cp "dist/${GAME_NAME}_macos" "${MACOS_DIR}/${GAME_NAME}"
chmod +x "${MACOS_DIR}/${GAME_NAME}"

# 创建图标文件（如果存在player.png，使用它作为图标）
if [ -f "resources/images/player.png" ]; then
    cp "resources/images/player.png" "${RESOURCES_DIR}/icon.png"
elif [ -f "resources/images/enemy.png" ]; then
    cp "resources/images/enemy.png" "${RESOURCES_DIR}/icon.png"
fi

# 创建Info.plist文件
cat > "${CONTENTS_DIR}/Info.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>${GAME_NAME}</string>
    <key>CFBundleIconFile</key>
    <string>icon.png</string>
    <key>CFBundleIdentifier</key>
    <string>com.playgame.planewar</string>
    <key>CFBundleInfoDictionaryVersion</key>
    <string>6.0</string>
    <key>CFBundleName</key>
    <string>${GAME_NAME}</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0.0</string>
    <key>CFBundleVersion</key>
    <string>1</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.14</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>NSHumanReadableCopyright</key>
    <string>Copyright © 2025 飞机大战开发团队. All rights reserved.</string>
</dict>
</plist>

EOF

# 清理临时文件
rm "dist/${GAME_NAME}_macos"

echo -e "${GREEN}macOS版本构建成功! 输出: ${APP_DIR}${NC}"
echo -e "${YELLOW}注意: macOS应用需要代码签名才能正常使用，未签名的应用可能会被系统拦截。${NC}"
