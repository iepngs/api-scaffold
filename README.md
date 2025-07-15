

### 文档生成

```shell
make doc
```

### 运行
```shell
make run
# 指定配置文件 (未手动指定时，使用项目根目录下的 config.yaml 文件)
go run cmd/main.go -c config.local.yaml
```

### 编译打包
```shell
make build
# 交叉编译
make build-cross

# 指定版本号（优先读取指定传参版本号，没有指定则降级取当前tag的版本，默认： dev ）
make VERSION=v1.2.3 build
```