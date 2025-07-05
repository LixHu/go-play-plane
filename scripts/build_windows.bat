@echo off
setlocal enabledelayedexpansion

echo 飞机大战游戏打包工具 - Windows版本构建
echo ======================================

:: 创建输出目录
if not exist dist mkdir dist

:: 检查Go环境
where go >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo 错误: 未检测到Go环境，请安装Go后再运行此脚本。
    exit /b 1
)

for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i
echo 检测到Go版本: %GO_VERSION%

:: 检查purego版本
for /f "tokens=2" %%v in ('go list -m github.com/ebitengine/purego') do set PUREGO_VERSION=%%v

if not "%PUREGO_VERSION%"=="v0.8.3" (
    echo 检测到purego版本 %PUREGO_VERSION%，建议更新到v0.8.3以解决符号重复问题
    set /p CHOICE="是否更新purego到v0.8.3? (y/n): "
    if /i "!CHOICE!"=="y" (
        echo 正在更新purego...
        go get github.com/ebitengine/purego@v0.8.3
        echo 已更新purego到v0.8.3
    )
)

:: 检查依赖
echo 正在检查依赖...
go mod tidy
if %ERRORLEVEL% neq 0 (
    echo 依赖安装失败
    exit /b 1
)

echo 确保依赖完整下载...
go mod download
if %ERRORLEVEL% neq 0 (
    echo 依赖下载失败
    exit /b 1
)

:: 开始构建Windows版本
echo.
echo 开始构建Windows版本...

:: 设置环境变量
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1

:: 编译
go build -o "dist\飞机大战.exe" -ldflags="-H windowsgui" .

if %ERRORLEVEL% equ 0 (
    echo Windows版本构建成功! 输出: dist\飞机大战.exe
) else (
    echo Windows版本构建失败!
    exit /b 1
)

echo.
echo 按任意键退出...
pause >nul
