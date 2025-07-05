@echo off
echo 正在启动飞机大战游戏...

:: 检查purego版本
for /f "tokens=2" %%v in ('go list -m github.com/ebitengine/purego 2^>nul') do set PUREGO_VERSION=%%v

if not "%PUREGO_VERSION%"=="v0.8.3" (
    echo 警告: 检测到purego版本 %PUREGO_VERSION%，可能会出现符号重复定义问题
    echo 建议先运行 scripts\update_dependencies.bat 修复依赖
    echo.
    set /p CHOICE="是否继续? (y/n): "
    if /i not "%CHOICE%"=="y" exit /b
)

:: 运行游戏
go run .

if %ERRORLEVEL% neq 0 (
    echo.
    echo 游戏启动失败，可能是依赖问题
    echo 尝试运行 scripts\update_dependencies.bat 修复依赖后再试
)

pause
