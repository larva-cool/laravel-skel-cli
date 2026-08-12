# laravel-skel-cli

将后端 API 接口封装成可在命令行执行的工具。接口清单从 Apifox 项目自动生成，采用数据驱动方式，新增/变更接口无需改代码。

## 特性

- 🔐 **Bearer Token 认证**：登录后自动保存并携带 token
- 📦 **数据驱动**：52 个接口由 OpenAPI 规范自动生成，无需手写重复命令
- 🔄 **跨平台**：单文件可执行二进制，无需运行时依赖
- 📋 **类型安全参数**：自动将命令行参数转换为 string / integer / number / boolean / array
- ⚙️ **本地配置**：base_url 与 token 持久化于用户配置目录（权限 0600）

## 快速开始

### 1. 下载二进制

从 [Releases](https://github.com/larva-cool/laravel-skel-cli/releases) 下载对应平台的二进制（macOS / Linux / Windows），赋予执行权限：

```bash
chmod +x laravel-skel-cli
```

> 或将 `laravel-skel-cli` 放入 `$PATH`（如 `/usr/local/bin`）以便全局调用。

### 2. 配置后端地址

```bash
laravel-skel-cli config set-base-url http://你的后端地址
```

### 3. 登录

```bash
laravel-skel-cli login --account admin --password 你的密码
```

### 4. 调用接口

```bash
laravel-skel-cli list                       # 查看全部可调用接口
laravel-skel-cli whoami                     # 当前登录管理员信息
laravel-skel-cli call admins.list           # 管理员列表
laravel-skel-cli call admins.update --admin_id 10000000 --name 新名字
```

## 命令说明

```text
认证:
  login --account <账号> --password <密码> [--base-url <地址>]   登录并保存 token
  logout                                                       退出并清除 token
  whoami                                                       获取当前登录管理员信息
  config set-base-url <地址>   设置后端 API 地址到配置
  config set-token <token>     设置 token 到配置
  config get                  查看当前配置

接口调用（数据驱动）:
  call <slug> [--参数 值]   调用指定接口并输出 JSON
  list                      列出全部可调用接口

其他:
  version     显示版本号
  help        显示帮助信息
```

### call 参数规则

- **路径参数**：如 `admins.show --admin_id 10000000` 会替换路径中的 `{admin_id}`
- **查询参数**：如 `admins.list --keyword admin --status 1`
- **请求体**：字段作为 `--字段名` 传入（仅发送显式提供的字段，便于部分更新）
- **原始 JSON 请求体**：`--body '[1,2,3]'` 直接透传（用于数组等特殊接口）
- **数组类型**：逗号分隔，如 `--roles 1,2,3`
- **覆盖地址**：所有命令支持 `--base-url` 临时覆盖配置

## 从 Apifox 更新接口清单

接口注册表由 Apifox 导出的 OpenAPI 规范生成：

```bash
# 1. 用 apifox-cli 导出 OpenAPI
apifox export --project <项目ID> --format openapi --oas-version 3.0 --output openapi.json

# 2. 重新生成注册表
python3 tools/gen_registry.py openapi.json internal/apidefs/registry.go
```

## 开发与构建

```bash
make build                # 构建当前平台二进制 → dist/
make release              # 构建全部跨平台二进制
make clean                # 清理 dist
make test                 # 静态检查 + 构建

# 指定版本号（注入 version 命令）
make VERSION=v1.1.0 build
```

## 测试

CI 自动执行 `go vet` + `go build` + `go test -race -cover`。本地运行：

```bash
go test -race -cover ./...
```

## 项目结构

```
laravel-skel-cli/
├── main.go                  # 入口
├── Makefile                 # 构建脚本
├── tools/gen_registry.py    # 从 OpenAPI 生成接口注册表
├── internal/
│   ├── apidefs/             # 接口定义 + 生成的 52 接口注册表
│   ├── client/              # HTTP 客户端（Bearer/JSON/错误解析）
│   └── config/              # 配置持久化
├── cmd/                     # 命令层（root/login/logout/config/call/list）
└── .github/workflows/       # CI 与 Release 工作流
```

## License

[MIT](LICENSE)
