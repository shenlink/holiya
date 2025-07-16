# Holiya 编程语言  

Holiya 是一个基于 Go 语言实现的解释型编程语言项目，参考《用Go语言自制解释器》一书的代码，在原书代码的基础上进行了一些修改和增强。

## 🚀 快速开始
### 1. 安装 Go 语言环境
请参考 [Go 语言安装指南](https://golang.org/doc/install)

### 2. 克隆项目
- [Gitee 地址](https://gitee.com/shenlink/holiya)
```shell
git clone https://gitee.com/shenlink/holiya.git
```

- [GitHub 地址](https://github.com/shenlink/holiya)
```shell
git clone https://github.com/shenlink/holiya.git
```

### 3. 编译项目
windows:
```shell
go build -o holiya.exe main.go
```
linux:
```shell
go build -o holiya main.go
```

### 4. 使用项目
- 可以直接运行编译后的 holiya.exe 或者 holiya，然后在 repl 中输入文本，会即时执行输入的文本
- 可以使用 holiya.exe filename.holiya 或 holiya filename.holiya，holiya 会自动执行该文件

### 5. 测试项目
```shell
go test ./...
```

## 📄 许可证
本项目采用 [MIT 许可证](LICENSE)