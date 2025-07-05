@echo off
setlocal enabledelayedexpansion

echo 使用Fyne打包飞机大战游戏 - Windows版本
echo =======================================

:: 检查Fyne CLI是否已安装
where fyne >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo 错误: 未检测到Fyne CLI工具，请先安装:
    echo go install fyne.io/fyne/v2/cmd/fyne@latest
    exit /b 1
)

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
        go mod tidy
        echo 已更新purego到v0.8.3
    )
)

:: 创建输出目录
if not exist dist mkdir dist

:: 游戏图标
set ICON_PATH=resources\images\player.png
if not exist "%ICON_PATH%" (
    echo 警告: 未找到玩家图标，将使用默认图标
    set ICON_PATH=
)

echo.
echo 开始为Windows打包...

if not "%ICON_PATH%"=="" (
    fyne package -os windows -icon "%ICON_PATH%" -name "飞机大战" -appID "com.playgame.planewar"
) else (
    fyne package -os windows -name "飞机大战" -appID "com.playgame.planewar"
)

if %ERRORLEVEL% equ 0 (
    :: 移动到dist目录
    move 飞机大战.exe dist\
    echo Windows打包成功! 输出: dist\飞机大战.exe
) else (
    echo Windows打包失败!
    exit /b 1
)

echo.
echo 打包过程完成!
echo 按任意键退出...
pause >nul
