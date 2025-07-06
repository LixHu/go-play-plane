# 飞机大战游戏

一个使用Go语言和Ebiten游戏引擎开发的经典飞机射击游戏。

## 游戏特点

- 精美的动画启动界面
- 两种游戏模式：关卡模式和无尽模式
- 多种敌机类型
- 丰富的武器升级系统
- 流畅的游戏体验

## 操作说明

- 方向键：移动飞机
- 空格键：发射子弹
- R键：游戏结束时重新开始
- ESC键：返回菜单

## 开发环境

- Go 1.24
- Ebiten v2

## 启动方式

```bash
go run .
```

## 打包指南

> 注意：所有打包脚本已移至 `scripts` 文件夹，请使用该文件夹中的脚本进行构建。

### 统一打包脚本（推荐）

使用我们的统一打包脚本，可以轻松打包Windows、macOS和Linux版本：

```bash
# 授予脚本执行权限
chmod +x scripts/build.sh

# 运行打包脚本
./scripts/build.sh
```

脚本会提供一个交互式菜单，让您选择要打包的目标平台。

### 方法一：直接编译生成exe文件

在Windows系统上：
```bash
# 运行打包脚本
scripts\build_windows.bat
```

在Linux/macOS系统上：
```bash
# 授予脚本执行权限
chmod +x scripts/build_windows.sh
# 运行打包脚本
./scripts/build_windows.sh
```

### 方法二：使用Fyne打包(需要先安装Fyne CLI)

在Windows系统上：
```bash
# 运行Fyne打包脚本
scripts\package_with_fyne.bat
```

在Linux/macOS系统上：
```bash
# 授予脚本执行权限
chmod +x scripts/package_with_fyne.sh
# 运行Fyne打包脚本
./scripts/package_with_fyne.sh
```

## 游戏截图

（这里可以添加游戏截图）

## macOS应用说明

使用我们的统一打包脚本可以生成macOS应用（.app格式）。请注意，由于macOS的安全策略，未签名的应用可能会被系统拦截。

### macOS构建注意事项

在macOS上构建需要安装Xcode命令行工具并启用CGO：

```bash
# 安装Xcode命令行工具
xcode-select --install

# 使用专用macOS构建脚本
chmod +x scripts/build_macos.sh
./scripts/build_macos.sh
```

如果遇到`glfw.Window`未定义等错误，请参考`scripts/MACOS_BUILD_GUIDE.md`获取详细解决方案。

### 代码签名和公证

如果您需要分发应用，建议进行代码签名和公证：

1. 获取Apple开发者证书
2. 使用以下命令签名：
   ```bash
   codesign --force --deep --sign "Developer ID Application: Your Name" dist/飞机大战.app
   ```
3. 提交公证：
   ```bash
   xcrun altool --notarize-app --primary-bundle-id "com.playgame.planewar" --username "your@apple.id" --password "app-specific-password" --file dist/飞机大战.app
   ```

## 许可证

MIT License
- 可以到release 下载，有打好的包(更新慢，每次一个小版本会更新)

- 也可以把项目下载下来自己编译 (建议)

```
// 首先要装好Go环境
go mod tidy
# Windows
go build -o 飞机大战.exe
# macOS/Linux
go build -o 飞机大战
```