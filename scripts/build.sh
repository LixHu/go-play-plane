#!/bin/bash

# 飞机大战游戏打包脚本
# 支持 Windows、macOS 和 Linux

# 颜色定义
RED="\033[0;31m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
BLUE="\033[0;34m"
NC="\033[0m" # No Color

# 创建输出目录
mkdir -p dist

# 游戏名称
GAME_NAME="飞机大战"

# 检测当前操作系统
OS="$(uname -s)"
case "${OS}" in
    Linux*)
        CURRENT_OS="linux"
        ;;
    Darwin*)
        CURRENT_OS="darwin"
        ;;
    CYGWIN*|MINGW*|MSYS*)
        CURRENT_OS="windows"
        ;;
    *)
        echo -e "${RED}不支持的操作系统: ${OS}${NC}"
        exit 1
        ;;
esac

# 显示菜单
show_menu() {
    clear
    echo -e "${BLUE}=======================================${NC}"
    echo -e "${BLUE}          飞机大战游戏打包工具         ${NC}"
    echo -e "${BLUE}=======================================${NC}"
    echo -e "当前系统: ${GREEN}${CURRENT_OS}${NC}"
    echo -e "\n请选择打包目标平台:"
    echo -e "${YELLOW}1. Windows (exe)${NC}"
    echo -e "${YELLOW}2. macOS (.app)${NC}"
    echo -e "${YELLOW}3. Linux${NC}"
    echo -e "${YELLOW}4. 全部平台${NC}"
    echo -e "${YELLOW}0. 退出${NC}"
    echo -e "\n请输入选项 [0-4]: "
}

# 检查Go环境
check_go() {
    if ! command -v go &> /dev/null; then
        echo -e "${RED}错误: 未检测到Go环境，请安装Go后再运行此脚本。${NC}"
        exit 1
    fi

    GO_VERSION=$(go version | grep -oP "go\d+\.\d+" | cut -c 3-)
    echo -e "${GREEN}检测到Go版本: ${GO_VERSION}${NC}"
}

# 检查 macOS 构建环境
check_macos_build_env() {
    echo -e "${YELLOW}检查 macOS 构建环境...${NC}"

    # 检查 Xcode 命令行工具
    if ! command -v clang &> /dev/null; then
        echo -e "${RED}错误: 未检测到 Xcode 命令行工具，请运行: xcode-select --install${NC}"
        return 1
    fi

    # 检查 CGO 环境
    if [ "$(go env CGO_ENABLED)" != "1" ]; then
        echo -e "${YELLOW}警告: CGO 未启用，正在启用...${NC}"
        export CGO_ENABLED=1
    fi

    echo -e "${GREEN}macOS 构建环境检查完成${NC}"
    return 0
}

# 安装依赖
install_dependencies() {
    echo -e "${YELLOW}正在检查依赖...${NC}"

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
}

# 编译Windows版本
build_windows() {
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
        return 1
    fi

    return 0
}

# 编译macOS版本
build_macos() {
    echo -e "\n${BLUE}开始构建macOS版本...${NC}"

    # 检查 macOS 构建环境
    if ! check_macos_build_env; then
        return 1
    fi

    # 设置环境变量
    export GOOS=darwin
    export GOARCH=amd64
    export CGO_ENABLED=1

    # 编译
    go build -o "dist/${GAME_NAME}_macos" .

    if [ $? -ne 0 ]; then
        echo -e "${RED}macOS版本构建失败!${NC}"
        return 1
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

    return 0
}

# 编译Linux版本
build_linux() {
    echo -e "\n${BLUE}开始构建Linux版本...${NC}"

    # 设置环境变量
    export GOOS=linux
    export GOARCH=amd64
    export CGO_ENABLED=1

    # 编译
    go build -o "dist/${GAME_NAME}_linux" .

    if [ $? -eq 0 ]; then
        # 设置可执行权限
        chmod +x "dist/${GAME_NAME}_linux"
        echo -e "${GREEN}Linux版本构建成功! 输出: dist/${GAME_NAME}_linux${NC}"
    else
        echo -e "${RED}Linux版本构建失败!${NC}"
        return 1
    fi

    return 0
}

# 编译所有平台
build_all() {
    local success=true

    build_windows
    if [ $? -ne 0 ]; then success=false; fi

    build_macos
    if [ $? -ne 0 ]; then success=false; fi

    build_linux
    if [ $? -ne 0 ]; then success=false; fi

    if $success; then
        echo -e "\n${GREEN}所有平台构建完成!${NC}"
    else
        echo -e "\n${YELLOW}部分平台构建失败，请检查错误信息。${NC}"
    fi
}

# 主函数
main() {
    check_go
    install_dependencies

    while true; do
        show_menu
        read -r choice

        case $choice in
            0)
                echo -e "${BLUE}谢谢使用，再见!${NC}"
                exit 0
                ;;
            1)
                build_windows
                ;;
            2)
                build_macos
                ;;
            3)
                build_linux
                ;;
            4)
                build_all
                ;;
            *)
                echo -e "${RED}无效选项，请重新选择${NC}"
                ;;
        esac

        echo -e "\n按Enter键继续..."
        read -r
    done
}

# 运行主函数
main
