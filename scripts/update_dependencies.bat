@echo off
echo 正在修复purego依赖问题...

echo 清理go模块缓存...
go clean -modcache

echo 更新依赖到purego v0.8.3...
go get github.com/ebitengine/purego@v0.8.3

echo 整理go.mod文件...
go mod tidy

echo 验证purego版本...
for /f "tokens=2" %%a in ('go list -m github.com/ebitengine/purego') do set PUREGO_VERSION=%%a

echo 当前purego版本: %PUREGO_VERSION%

if "%PUREGO_VERSION%"=="v0.8.3" (
    echo ✅ 成功更新到purego v0.8.3
    echo 现在可以尝试重新构建项目
) else (
    echo ❌ purego更新失败，请手动编辑go.mod文件
    echo 将github.com/ebitengine/purego版本改为v0.8.3
)

pause
