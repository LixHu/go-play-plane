# macOS构建指南

本指南提供了在macOS上构建飞机大战游戏的详细步骤和解决常见问题的方法。

## 环境准备

### 1. 安装Xcode命令行工具

```bash
xcode-select --install
```

### 2. 确保CGO启用

```bash
export CGO_ENABLED=1
```

## 常见问题及解决方案

### 1. glfw.Window未定义错误

如果遇到类似以下错误：

```
undefined: glfw.Window
```

这通常是由于GLFW库无法正确初始化引起的。解决方法：

1. 确保安装了Xcode命令行工具
2. 启用CGO
3. 使用以下命令安装GLFW依赖：

```bash
brew install glfw
```

### 2. OpenGL相关错误

如果遇到OpenGL相关的错误，确保已正确设置OpenGL上下文。在代码中，确保有类似以下的设置：

```go
if runtime.GOOS == "darwin" {
    if err := glfw.WindowHint(glfw.OpenGLForwardCompat, glfw.True); err != nil {
        return err
    }
    if err := glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile); err != nil {
        return err
    }
}
```

### 3. 构建.app包时出错

如果在创建.app包时遇到问题，请尝试以下步骤：

1. 确保有正确的目录结构：
   - `飞机大战.app/Contents/MacOS/飞机大战`（可执行文件）
   - `飞机大战.app/Contents/Resources/`（资源文件）
   - `飞机大战.app/Contents/Info.plist`（应用配置）

2. 确保Info.plist文件包含必要的键值对

3. 使用以下命令检查.app包是否有效：

```bash
spctl --assess -v /path/to/飞机大战.app
```

## 代码签名和公证

如需发布应用，请参考以下步骤进行代码签名和公证：

### 1. 代码签名

```bash
codesign --force --deep --sign "Developer ID Application: Your Name" 飞机大战.app
```

### 2. 创建ZIP归档

```bash
ditto -c -k --keepParent 飞机大战.app 飞机大战.zip
```

### 3. 提交公证

```bash
xcrun altool --notarize-app --primary-bundle-id "com.playgame.planewar" --username "your@apple.id" --password "app-specific-password" --file 飞机大战.zip
```

### 4. 检查公证状态

```bash
xcrun altool --notarization-info [REQUEST_UUID] --username "your@apple.id" --password "app-specific-password"
```

### 5. 添加公证票据

```bash
xcrun stapler staple 飞机大战.app
```

## 使用专用构建脚本

为简化以上过程，建议使用我们的专用macOS构建脚本：

```bash
chmod +x scripts/build_macos.sh
./scripts/build_macos.sh
```

该脚本会自动处理大多数常见问题并生成有效的.app包。
