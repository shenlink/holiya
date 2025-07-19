# Go 可执行文件的名称
GO = go

# 存放程序入口文件（如 main.go）的目录
SRC_DIR = ./main.go
EXEC_PREFIX = ./

# 检测操作系统类型
ifeq ($(OS),Windows_NT)
    # 输出的二进制文件名
    BINARY_NAME = holiya.exe
    # Windows 环境使用 del 命令删除文件（通过 cmd 执行）
    REMOVE_CMD = cmd /c del /f /q
else
    # 输出的二进制文件名
    BINARY_NAME = holiya
    # 类 Unix 环境使用 rm 命令删除文件
    REMOVE_CMD = rm -f
endif

# 默认目标，执行 go build 命令
all: build

# 编译目标
build:
	$(GO) build -o $(BINARY_NAME) $(SRC_DIR)

# 清理目标，删除生成的二进制文件
clean:
	$(REMOVE_CMD) $(BINARY_NAME)

# 运行目标，编译并执行程序
run: build
	$(EXEC_PREFIX)$(BINARY_NAME)

# 安装依赖目标
deps:
	$(GO) mod tidy

# 运行 token 目录下的测试
test_token:
	$(GO) test -timeout 30s ./token

# 运行 lexer 目录下的测试
test_lexer:
	$(GO) test -timeout 30s ./lexer

# 运行所有测试
test:
	$(GO) test -timeout 30s ./...

# 声明伪目标
.PHONY: all build clean run test deps test_token test_lexer
