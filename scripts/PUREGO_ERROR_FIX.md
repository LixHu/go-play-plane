# 解决 purego 符号重复定义问题

## 问题说明

在构建过程中如果遇到以下错误：

```
link: duplicated definition of symbol dlopen, from github.com/ebitengine/purego and github.com/ebitengine/purego
```

这是由于 purego 库 v0.8.0 版本存在的已知问题，需要更新到 v0.8.3 版本解决。

## 解决方法

### 方法1：使用提供的修复脚本

```bash
# 授予脚本执行权限
chmod +x scripts/fix_purego.sh

# 运行修复脚本
./scripts/fix_purego.sh
```

### 方法2：手动修复

1. 清理 Go 模块缓存
```bash
go clean -modcache
```

2. 更新 purego 依赖
```bash
go get github.com/ebitengine/purego@v0.8.3
```

3. 整理 go.mod 文件
```bash
go mod tidy
```

4. 验证 purego 版本
```bash
go list -m github.com/ebitengine/purego
```

确保输出显示的是 `v0.8.3` 版本。

## 技术背景

purego 是一个允许在不使用 CGO 的情况下调用 C 函数的库。在 v0.8.0 版本中，存在符号重复定义的问题，这会在链接阶段导致错误。此问题已在 v0.8.3 版本中修复。

## 参考链接

- [ebitengine/purego 问题修复](https://github.com/rclone/rclone/issues/8552)
- [Go 1.24.3 regression issue](https://github.com/golang/go/issues/73617)
