# 飞机大战游戏构建脚本

此目录包含用于构建和打包飞机大战游戏的各种脚本。

## 通用脚本

- `build.sh` - 统一打包脚本，支持Windows、macOS和Linux平台
- `fix_purego.sh` - 修复purego依赖问题的脚本
- `run_game.sh` - Linux/macOS上便捷运行游戏的脚本
- `run_game.bat` - Windows上便捷运行游戏的脚本

## Windows专用脚本

- `build_windows.bat` - Windows专用构建脚本
- `build_windows.sh` - 从Linux/macOS构建Windows版本的脚本
- `update_dependencies.bat` - Windows环境下修复依赖的批处理脚本

## macOS专用脚本

- `build_macos.sh` - macOS专用构建脚本
- `package_with_fyne_macos.sh` - 使用Fyne打包macOS应用的专用脚本

## 通用打包脚本

- `package_with_fyne.sh` - 使用Fyne进行跨平台打包的Linux/macOS脚本
- `package_with_fyne.bat` - 使用Fyne进行Windows打包的批处理脚本

## 使用说明

### 在Linux/macOS上：

```bash
# 首先赋予脚本执行权限
chmod +x *.sh

# 运行统一打包脚本
./build.sh
```

### 在Windows上：

```batch
# 运行Windows专用构建脚本
build_windows.bat
```

## 文档文件

- `PUREGO_ERROR_FIX.md` - 解决purego符号重复定义问题的详细文档
- `MACOS_BUILD_GUIDE.md` - macOS系统构建指南和常见问题解决方案

## 注意事项

- 所有脚本都会检查并提示是否更新purego依赖到v0.8.3版本以解决符号重复定义问题
- 使用Fyne打包脚本前请确保已安装Fyne CLI工具：`go install fyne.io/fyne/v2/cmd/fyne@latest`
- macOS构建脚本需要安装Xcode命令行工具
