---
name: "laravel-skel-cli"
description: "通过 CLI 调用 Laravel 后端管理 API（52 个接口，涵盖管理员/角色/菜单/地区/配置/通知/聊天等）。Invoke when user asks to manage backend data, call admin APIs, or operate backend resources via CLI."
---

# laravel-skel-cli

将 Laravel 后端管理 API 封装为命令行工具。通过 `call <slug> [--参数 值]` 调用接口，返回 JSON。

## 二进制位置

```bash
# 项目构建产物
./dist/laravel-skel-cli

# 或全局安装后
laravel-skel-cli
```

后续命令示例统一用 `laravel-skel-cli`，若未安装到 PATH 请用 `./dist/laravel-skel-cli` 或 `go run .` 替代。

## 首次配置

```bash
# 1. 设置后端地址
laravel-skel-cli config set-base-url http://your-backend.test

# 2. 登录（保存 token）
laravel-skel-cli login --account admin --password 12345678

# 3. 验证
laravel-skel-cli whoami
```

配置文件位于 `~/.config/laravel-skel-cli/config.json`（权限 0600）。

也可手动设置 token（跳过登录）：

```bash
laravel-skel-cli config set-token <token>
```

## 命令总览

| 命令 | 说明 |
|------|------|
| `login --account <账号> --password <密码> [--base-url <地址>]` | 登录并保存 token |
| `logout` | 退出并清除 token |
| `whoami` | 获取当前登录管理员信息（等价 `call auth.info`） |
| `config set-base-url <地址>` | 设置后端 API 地址 |
| `config set-token <token>` | 手动设置 token |
| `config get` | 查看当前配置（token 打码显示） |
| `call <slug> [--参数 值]` | 调用指定接口，输出 JSON |
| `list` | 列出全部可调用接口 |
| `version` | 显示版本号 |

## call 参数规则

- **路径参数**（path）：如 `admins.show --admin_id 10000000`，替换 URL 中的 `{admin_id}`
- **查询参数**（query）：如 `admins.list --keyword abc --page 1`
- **请求体参数**（body）：字段名作为 flag，如 `admins.create --username test --name 测试 --password 12345678 --status 1`
- **原始 JSON 请求体**：`--body '[1,2,3]'`（用于数组类接口如 `notifications.mark-read`）
- **数组类型**：逗号分隔，如 `--roles super_admin,editor`
- **布尔值**：`1/true/yes` 或 `0/false/no`
- **覆盖地址**：所有命令支持 `--base-url` 临时覆盖配置

## 全部接口清单（52 个）

### 认证（auth）

| Slug | 方法 | 路径 | 参数 |
|------|------|------|------|
| `auth.login` | POST | /admin/auth/login | body: `account`(必填,string), `password`(必填,string) |
| `auth.info` | GET | /admin/auth/info | 无 |
| `auth.logout` | POST | /admin/auth/logout | 无 |

### 管理员（admins）

| Slug | 方法 | 路径 | 参数 |
|------|------|------|------|
| `admins.list` | GET | /admin/admins | query: `keyword`(string), `status`(integer: 0=禁用,1=正常), `page`, `per_page` |
| `admins.create` | POST | /admin/admins | body: `username`(必填,string), `email`(string), `phone`(string), `name`(必填,string), `password`(必填,string,至少8位含字母数字), `status`(必填,integer: 0/1), `roles`(array: 角色名称数组) |
| `admins.show` | GET | /admin/admins/{admin_id} | path: `admin_id`(必填,integer) |
| `admins.update` | PUT | /admin/admins/{admin_id} | path: `admin_id`(必填); body: `email`, `phone`, `name`(必填), `password`(留空不修改), `status`(必填,integer), `roles`(array) |
| `admins.delete` | DELETE | /admin/admins/{admin_id} | path: `admin_id`(必填,integer) |
| `admins.reset-password` | PUT | /admin/admins/{admin_id}/reset-password | path: `admin_id`(必填); body: `password`(必填), `password_confirmation`(必填) |
| `admins.roles` | GET | /admin/admins/{admin_id}/roles | path: `admin_id`(必填,integer) |
| `admins.roles-put` | PUT | /admin/admins/{admin_id}/roles | path: `admin_id`(必填); body: `roles`(必填,array: 角色名称数组) |
| `admins.toggle-status` | PUT | /admin/admins/{admin_id}/toggle-status | path: `admin_id`(必填,integer) |
| `admins.login-histories` | GET | /admin/admins/{id}/login-histories | path: `id`(必填); query: `keyword`, `page`, `per_page` |

### 角色（roles）

| Slug | 方法 | 路径 | 参数 |
|------|------|------|------|
| `roles.list` | GET | /admin/roles | query: `role_name`(string), `page`, `per_page` |
| `roles.create` | POST | /admin/roles | body: `name`(必填,string: 2-50字符,唯一), `permissions`(array: 权限ID数组) |
| `roles.show` | GET | /admin/roles/{role_id} | path: `role_id`(必填,integer) |
| `roles.update` | PUT | /admin/roles/{role_id} | path: `role_id`(必填); body: `name`(必填,string), `permissions`(array) |
| `roles.delete` | DELETE | /admin/roles/{role_id} | path: `role_id`(必填,integer) |
| `roles.permissions` | GET | /admin/roles/permissions | 无 |
| `roles.permissions-get` | GET | /admin/roles/{role_id}/permissions | path: `role_id`(必填,integer) |
| `roles.permissions-put` | PUT | /admin/roles/{role_id}/permissions | path: `role_id`(必填); body: `permissions`(必填,array: 权限ID数组) |

### 菜单（menus）

| Slug | 方法 | 路径 | 参数 |
|------|------|------|------|
| `menus.list` | GET | /admin/menus | 无（返回菜单树） |
| `menus.create` | POST | /admin/menus | body: `parent_id`(必填,integer: 0=顶级), `title`(必填,string), `type`(必填,integer: 0=目录,1=菜单,2=按钮,3=iframe,4=外链), `sort`(必填,integer: 0-9999), `is_enable`(必填,boolean), `is_hide`(必填,boolean), `is_hide_tab`(必填,boolean), `is_iframe`(必填,boolean), `keep_alive`(必填,boolean), `is_full_page`(必填,boolean), `fixed_tab`(必填,boolean), `show_badge`(必填,boolean), `path`, `name`, `component`, `redirect`, `icon`, `link`, `show_text_badge`, `active_path`, `permission`, `roles`(array) |
| `menus.show` | GET | /admin/menus/{menu_id} | path: `menu_id`(必填,integer) |
| `menus.update` | PUT | /admin/menus/{menu_id} | path: `menu_id`(必填); body 同 create |
| `menus.delete` | DELETE | /admin/menus/{menu_id} | path: `menu_id`(必填,integer) |

### 地区（areas）

| Slug | 方法 | 路径 | 参数 |
|------|------|------|------|
| `areas.list` | GET | /admin/areas | 无 |
| `areas.create` | POST | /admin/areas | body: `name`(必填,string), `parent_id`(integer), `area_code`(integer), `lat`(number), `lng`(number), `city_code`(string), `sort`(integer) |
| `areas.show` | GET | /admin/areas/{area} | path: `area`(必填,integer: 地区ID) |
| `areas.update` | PUT | /admin/areas/{area} | path: `area`(必填); body 同 create |
| `areas.delete` | DELETE | /admin/areas/{area} | path: `area`(必填,integer) |

### 配置（settings）

| Slug | 方法 | 路径 | 参数 |
|------|------|------|------|
| `settings.list` | GET | /admin/settings | query: `keyword`(string), `cast_type`(string: string/int/bool/float/json), `page`, `per_page` |
| `settings.create` | POST | /admin/settings | body: `name`(必填,string), `key`(必填,string: 唯一), `cast_type`(必填,string), `input_type`(必填,string), `value`(string), `param`(string: JSON), `sort`(integer), `remark`(string) |
| `settings.show` | GET | /admin/settings/{setting} | path: `setting`(必填,integer: 配置ID) |
| `settings.update` | PUT | /admin/settings/{setting} | path: `setting`(必填); body 同 create |
| `settings.delete` | DELETE | /admin/settings/{setting} | path: `setting`(必填,integer) |

### 通知（notifications）

| Slug | 方法 | 路径 | 参数 |
|------|------|------|------|
| `notifications.list` | GET | /admin/notifications | query: `type`(string: 通知类名), `page`, `per_page` |
| `notifications.unread` | GET | /admin/notifications/unread | query: `type`, `page`, `per_page` |
| `notifications.mark-all-read` | PUT | /admin/notifications/mark-all-read | 无 |
| `notifications.mark-read` | PUT | /admin/notifications/mark-read | body: `--body '[1,2,3]'`(必填,原始JSON: 通知ID数组) |
| `notifications.clear-read` | DELETE | /admin/notifications/clear-read | 无 |

### AI 聊天（chat）

| Slug | 方法 | 路径 | 参数 |
|------|------|------|------|
| `chat.create` | POST | /admin/chat | body: `prompt`(必填,string), `conversation_id`(string: 继续对话时传入) |
| `chat.stream` | POST | /admin/chat/stream | body: 同 chat.create（SSE 流式响应） |
| `chat.conversations` | GET | /admin/chat/conversations | query: `page`, `per_page` |
| `chat.conversations-get` | GET | /admin/chat/conversations/{conversationId} | path: `conversationId`(必填,string: UUID) |
| `chat.conversations-delete` | DELETE | /admin/chat/conversations/{conversationId} | path: `conversationId`(必填,string: UUID) |

### 验证码（phone-codes / mail-codes）

| Slug | 方法 | 路径 | 参数 |
|------|------|------|------|
| `phone-codes.list` | GET | /admin/phone-codes | query: `phone`(string), `scene`(string), `state`(integer: 0=未使用,1=已使用) |
| `phone-codes.show` | GET | /admin/phone-codes/{id} | path: `id`(必填,integer) |
| `mail-codes.list` | GET | /admin/mail-codes | query: `email`(string), `state`(integer) |
| `mail-codes.show` | GET | /admin/mail-codes/{id} | path: `id`(必填,integer) |

### 其他

| Slug | 方法 | 路径 | 参数 |
|------|------|------|------|
| `routes.list` | GET | /admin/routes | 无（返回前端路由配置） |
| `uploader.token` | POST | /admin/uploader/token | body: `filename`(必填,string) |

## 常用示例

```bash
# 查看当前管理员
laravel-skel-cli whoami

# 管理员列表（带搜索和分页）
laravel-skel-cli call admins.list --keyword admin --status 1 --page 1 --per_page 20

# 创建管理员
laravel-skel-cli call admins.create --username test --name 测试 --password Abc12345 --status 1 --roles super_admin

# 修改管理员
laravel-skel-cli call admins.update --admin_id 10000001 --name 新名字 --status 1

# 查看角色权限
laravel-skel-cli call roles.permissions-get --role_id 1

# 创建菜单
laravel-skel-cli call menus.create --parent_id 0 --title 仪表盘 --type 1 --sort 1 --is_enable 1 --is_hide 0 --is_hide_tab 0 --is_iframe 0 --keep_alive 1 --is_full_page 0 --fixed_tab 0 --show_badge 0

# 发送聊天消息
laravel-skel-cli call chat.create --prompt 你好

# 查看会话列表
laravel-skel-cli call chat.conversations --page 1

# 标记通知已读
laravel-skel-cli call notifications.mark-read --body '[1,2,3]'
```

## 错误处理

- **401 Unauthenticated**：token 过期或未登录，需重新 `login`
- **422 Validation**：参数校验失败，错误信息在响应 `message` 或 `errors` 字段
- **404 Not Found**：路径参数对应的资源不存在
- **500 Server Error**：后端异常，CLI 会透出完整错误信息

## 更新接口清单

后端接口变更时，从 Apifox 重新生成注册表：

```bash
# 1. 导出 OpenAPI（需要 apifox-cli）
apifox export --project 8653021 --format openapi --oas-version 3.0 --output openapi.json

# 2. 重新生成
python3 tools/gen_registry.py openapi.json internal/apidefs/registry.go

# 3. 重新构建
make build
```
