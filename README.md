# go_sleep_admin

睡眠内容管理后台。

- 后端目录：`backend/`
- 前端目录：`frontend/`
- 当前后端架构已经按 `go-gin-api` 的运行时骨架重构：
  `main -> cmd -> bootstrap -> platform/config|data -> router.Module -> handler/service/repository`
- 当前后端核心技术栈：
  `Gin + GORM + Casbin + gtkit/logger + gtkit/orm + gtkit/redis + gtkit/encry/jwt`

## 项目说明

这个仓库现在分成两部分：

- `backend/`：管理后台 API、鉴权、权限、上传、日志、配置、运行时装配。
- `frontend/`：Vue 3 + Vite + TypeScript + Arco Design 管理后台页面。

已实现的主要后端能力：

- 管理员登录、JWT 鉴权、在线会话记录。
- Casbin RBAC 权限控制。
- 后台菜单、角色、管理员、字典、系统参数、操作日志、登录日志、在线用户。
- 文章管理、内容菜单、上传文件中心。
- 部门、岗位、数据范围。
- 本地上传、阿里云 OSS、腾讯云 COS、华为云 OBS。

当前前端已经同步落地的重点功能：

- 首页工作台：
  - 常用功能分组卡片。
  - 第三版指标模块，集中展示内容量、资源量、管理员量、角色量、在线会话量等核心指标。
  - 财务指标模块，展示订单与流水的聚合结果。
  - 销售额趋势模块，包含柱状图和渠道分布饼图。
  - 三级内容菜单模块，直接读取 `menus` 表，并按 `level` 展示一级、二级、三级菜单。
- 文章管理：
  - 支持单篇新增、编辑、状态批量处理、删除。
  - 新增文章页支持富文本正文编辑、封面上传、资源上传、资源库选择。
  - 新增文章页的资源库弹窗默认不过滤 `scene`，可以直接选择在“资源上传”页上传的历史资源。
  - 新增文章页支持按文章类型上传视频、音频、图片正文资源，并自动写入正文内容。
- 批量新增文章：
  - 新增了独立的“批量新增”子菜单。
  - 支持“公共配置 + 多条文章明细”一次性提交。
  - 纯文本类型下，每一条正文都支持富文本编辑。
  - 非纯文本类型下，每一条正文都支持单独上传对应资源。
  - 每一条文章都支持单独上传封面覆盖；如果不上传，则继承公共封面。
- 资源上传与资源库：
  - 统一支持本地、阿里云 OSS、腾讯云 COS、华为云 OBS。
  - 上传记录支持按文件类型、场景、云存储类型筛选。
  - 文章页、批量新增页、菜单图标上传都可以复用资源库。
- 全局默认上传类型：
  - 前端增加了“全局默认上传类型”状态。
  - 只要在“资源上传”页或文章页设置一次，后续文章封面上传、正文资源上传、批量新增上传、菜单图标上传都会默认跟随该值。
  - 该默认值会持久化到浏览器本地存储，刷新页面后仍然生效。

默认内置账号：

- `admin / admin`
- `operator / 123456`
- `reviewer / 123456`

说明：

- `dev/test` 环境启动时会同步内置角色、账号、菜单与权限。
- `prod` 环境启动时只补齐缺失的内置数据，不覆盖现有生产数据。

## 页面与功能说明

### 首页工作台

- 首页不是静态展示页，而是通过后端聚合接口实时读取数据。
- “第三版指标”模块用于快速查看当前内容运营的核心状态。
- “财务指标”模块用于查看支付、退款、订单数量等财务数据。
- “销售额趋势”模块用于查看按月份聚合的订单金额走势，以及按渠道聚合的金额结构。
- “三级内容菜单”模块直接读取 `menus` 表，前端不会写死层级。

### 新增文章

- 纯文本类型：
  - 正文使用富文本编辑器。
  - 适合普通文章、说明文、图文混排内容。
- 非纯文本类型：
  - 正文通过上传视频、音频、图片资源生成。
  - 上传成功后，会自动把资源地址写入正文内容。
- 从资源库选择：
  - 打开弹窗后可以直接选择“资源上传”页里上传过的历史资源。
  - 支持按文件类型、云存储类型筛选。
- 上传类型：
  - 页面上的上传类型已经改成全局默认值。
  - 当前页面修改后，其他上传入口也会默认使用同一类型。

### 批量新增文章

- 入口位置：
  - 左侧导航“内容管理 -> 批量新增”
  - 文章列表页工具栏中的“批量新增”
- 适用场景：
  - 同一菜单、同一渠道下批量创建多篇内容。
  - 适合同步录入一组专题、栏目、课程类文章。
- 录入方式：
  - 可以手工逐条填写。
  - 也可以先用“文本导入”快速生成多条记录，再逐条补充。
- 正文规则：
  - 纯文本时，每条正文都使用富文本编辑器。
  - 非纯文本时，每条正文都需要单独上传资源。
- 封面规则：
  - 可以设置公共封面图，全部文章继承。
  - 也可以为单条文章上传自己的封面覆盖。

### 资源上传

- 资源上传页除了上传文件，还承担“全局默认上传类型设置”的作用。
- 当前支持：
  - 本地存储
  - 阿里云 OSS
  - 腾讯云 COS
  - 华为云 OBS
- 上传完成后的记录会进入资源库，可供文章正文、封面、菜单图标继续复用。

## 运行时调用链

```text
浏览器 / 前端页面
        │
        ▼
Nginx / 反向代理
        │
        ▼
Gin Router
        │
        ├── 公共中间件：CORS / ErrorHandler / RequestID / Recovery
        ├── 鉴权中间件：JWT
        ├── 权限中间件：Casbin
        └── 操作日志中间件：OperationLogger
        │
        ▼
router.Module
        │
        ▼
handler
        │
        ▼
service
        │
        ▼
repository
        │
        ├── MySQL  : gtkit/orm + GORM
        ├── Redis  : gtkit/redis
        └── Storage: local / OSS / COS / OBS
```

## 目录架构图

### 仓库总览

```text
go_sleep_admin/
├── README.md                           # 项目总说明，包含架构、开发、部署文档
├── backend/                            # 后端工程根目录
│   ├── main.go                         # 程序主入口，只负责执行 Cobra 命令
│   ├── cmd/                            # 命令行入口层
│   ├── config/                         # 配置资源目录，保存 yml 配置和 JWT PEM
│   ├── database/                       # 数据库迁移目录，预留给 gtkit/migrate
│   ├── internal/                       # 后端核心业务代码
│   ├── go.mod                          # 后端 Go Module 定义
│   ├── go.sum                          # 后端依赖锁定文件
│   ├── logs/                           # 运行日志输出目录
│   └── public/uploads/                 # 本地上传目录
├── frontend/                           # 前端管理后台
├── deploy/nginx/                       # Nginx 配置示例
└── deploy/systemd/                     # systemd 服务配置示例
```

### 后端详细架构

```text
backend/
├── main.go                             # 主入口：调用 cmd.Execute()
├── cmd/
│   ├── root.go                         # 根命令：统一处理 `server -c dev|test|prod`
│   ├── server.go                       # 服务启动命令：创建 App 并运行 HTTP 服务
│   └── registry.go                     # 命令注册中心：便于后续扩展 migrate/db 等命令
├── config/
│   ├── config.go                       # embed 配置入口，提供内嵌 yml fallback
│   ├── path.go                         # 配置路径解析辅助函数
│   ├── README.md                       # 配置子目录说明
│   ├── env/
│   │   ├── dev.yml                     # 开发环境配置
│   │   ├── test.yml                    # 测试环境配置
│   │   ├── prod.yml                    # 生产环境配置
│   └── pem/
│       ├── jwtpri.pem                  # JWT Ed25519 私钥
│       └── jwtpub.pem                  # JWT Ed25519 公钥
├── database/
│   └── migrations/                     # 迁移脚本目录，当前保留为空目录
├── internal/
│   ├── appmeta/                        # 应用元信息：命令名、API 前缀、迁移锁名
│   │   ├── meta.go                     # 项目名、命令名、缓存键等全局常量
│   │   └── http.go                     # `/api`、`/api/v1` 等 HTTP 根路径定义
│   ├── bootstrap/                      # 运行时装配层，负责把所有基础设施和业务对象接起来
│   │   ├── app.go                      # 应用生命周期管理：启动、信号退出、资源关闭
│   │   ├── runtime.go                  # Runtime 构建：加载配置并初始化底层依赖
│   │   ├── providers.go                # Provider 装配：DB/Redis/JWT/Casbin/Uploader/Module
│   │   ├── http.go                     # HTTP Server 配置、启动、优雅关闭
│   │   ├── worker.go                   # 后台 worker 管理器，当前为可扩展骨架
│   │   ├── container.go                # 旧业务三层的集中装配入口
│   │   └── seed.go                     # 内置角色/账号/菜单/权限/基础数据初始化
│   ├── dto/                            # 请求与响应 DTO
│   ├── handler/                        # HTTP Handler 层
│   │                                    # 说明：这是本次平滑迁移阶段保留的 legacy handler
│   │                                    # 当前已经不由 router 直接 new，而是统一交给 bootstrap 装配
│   ├── service/                        # 业务服务层
│   ├── repository/                     # 数据访问层
│   ├── models/                         # GORM 模型定义
│   ├── middleware/                     # Gin 中间件
│   │   ├── init.go                     # 中间件总入口
│   │   ├── auth.go                     # JWT 认证中间件
│   │   ├── permission.go               # Casbin 权限校验中间件
│   │   ├── operation_log.go            # 操作日志中间件
│   │   ├── error.go                    # 统一错误处理
│   │   └── cors.go                     # CORS 处理
│   ├── module/                         # 新路由模块层，对齐 go-gin-api 的 Module 注册模式
│   │   ├── auth/transport/http/        # 认证模块路由注册
│   │   ├── common/transport/http/      # 健康检查、模块元数据路由注册
│   │   ├── content/transport/http/     # 文章、菜单、上传路由注册
│   │   ├── rbac/transport/http/        # 管理员、角色、后台菜单路由注册
│   │   └── system/transport/http/      # 字典、系统参数、日志、在线用户、岗位、部门、任务路由注册
│   ├── pkg/                            # 项目内公共包
│   │   ├── response/                   # 统一响应结构
│   │   ├── apperror/                   # 业务错误定义
│   │   ├── jwtauth/                    # 基于 gtkit/encry 的 JWT Ed25519 初始化
│   │   ├── log/                        # 基于 gtkit/logger 的日志初始化
│   │   ├── news/                       # 告警/通知包装
│   │   └── file/                       # 文件运行时初始化
│   ├── platform/                       # 基础设施平台层
│   │   ├── config/                     # typed config 加载、校验、schema 管理
│   │   ├── data/                       # MySQL / Redis 初始化
│   │   ├── auth/                       # JWTManager 兼容层
│   │   └── storage/                    # local / OSS / COS / OBS 上传驱动
│   ├── router/
│   │   ├── init.go                     # Gin Engine 创建与 Module 路由注册
│   │   └── module.go                   # `router.Module` 接口定义
│   └── runtime/
│       └── resource/path.go            # 运行时资源路径解析
├── logs/                               # 默认日志目录
├── public/uploads/                     # 本地上传目录
├── go.mod                              # Go 依赖声明
└── go.sum                              # Go 依赖锁定
```

## 配置说明

### 后端配置加载方式

后端不再使用 `.env` 作为运行时主配置。

当前配置链路：

1. 启动命令通过 `-c` 指定环境名，例如 `dev`、`test`、`prod`
2. 运行时读取 `backend/config/env/<env>.yml`
3. 如果磁盘上没有该文件，则回退到 `config/config.go` 中 embed 的同名 yml
4. `internal/platform/config/load.go` 使用 `viper` 解码并做 schema 校验

关键配置文件：

- `backend/config/env/dev.yml`
- `backend/config/env/test.yml`
- `backend/config/env/prod.yml`

关键基础设施配置项：

- `application.*`：监听地址、端口、超时、时区
- `database.*`：MySQL 地址、账号、连接池
- `redis.*`：Redis 地址、密码、逻辑库列表
- `log.*`：日志级别、日志目录、SQL 日志开关
- `jwt.*`：JWT 过期时间
- `upload.*`：上传驱动、本地目录、云存储参数

### JWT 密钥

当前 JWT 使用 Ed25519 PEM 文件：

- `backend/config/pem/jwtpri.pem`
- `backend/config/pem/jwtpub.pem`

开发环境可以直接使用仓库内文件。

生产环境必须替换为你自己的密钥。推荐命令：

```bash
mkdir -p backend/config/pem
openssl genpkey -algorithm ED25519 -out backend/config/pem/jwtpri.pem
openssl pkey -in backend/config/pem/jwtpri.pem -pubout -out backend/config/pem/jwtpub.pem
chmod 600 backend/config/pem/jwtpri.pem
chmod 644 backend/config/pem/jwtpub.pem
```

## 开发环境部署

### 1. 环境要求

建议版本：

- Go `1.26.1+`
- Node.js `>= 18.12`，推荐 `24.x`
- pnpm `>= 8.7`
- MySQL `8.x`
- Redis `6.x / 7.x`
- Linux / macOS 均可

### 2. 初始化数据库与 Redis

创建数据库：

```sql
CREATE DATABASE sleep_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
```

启动 Redis，确保默认地址可连接：

```bash
127.0.0.1:6379
```

### 3. 配置后端开发环境

编辑：

```bash
backend/config/env/dev.yml
```

至少确认这些字段正确：

```yaml
application:
  host: "0.0.0.0"
  port: "8080"

database:
  host: "127.0.0.1"
  port: "3306"
  name: "sleep_admin"
  username: "root"
  password: ""

redis:
  addr: "127.0.0.1:6379"
  dbs: [0, 1, 2]

upload:
  driver: "local"
  local_dir: "public/uploads"
  public_path: "/uploads"
```

说明：

- 本地上传模式下，上传文件会落到 `backend/public/uploads/`
- `dev` 环境每次启动会同步内置账号、角色、菜单和权限

### 4. 启动后端

如果 `go` 已在 `PATH` 中：

```bash
cd backend
go run . server -c dev
```

如果当前机器的 Go 没加到 `PATH`：

```bash
cd backend
/usr/local/go/bin/go run . server -c dev
```

启动成功后检查：

```bash
curl http://127.0.0.1:8080/api/v1/healthz
```

预期返回：

```json
{"code":0,"message":"ok"}
```

### 5. 启动前端

安装依赖：

```bash
cd frontend
pnpm install
```

确认开发环境 API 地址：

```bash
frontend/.env.development
```

默认值：

```env
VITE_APP_BASE_URL = 'http://127.0.0.1:8080'
```

启动开发服务器：

```bash
pnpm dev
```

默认访问地址：

- 前端：`http://127.0.0.1:5173`
- 后端：`http://127.0.0.1:8080`

### 6. 开发环境验证

登录验证：

- 用户名：`admin`
- 密码：`admin`

建议额外验证以下页面与功能：

- 首页：
  - 检查常用功能、第三版指标、财务指标、销售额趋势、三级内容菜单是否正常展示。
- 新增文章：
  - 检查纯文本类型下富文本编辑器是否正常。
  - 检查非纯文本类型下正文资源上传、资源库选择、封面上传是否正常。
- 批量新增文章：
  - 检查文本导入、逐条富文本编辑、逐条资源上传、逐条封面覆盖是否正常。
- 资源上传：
  - 检查上传后是否能在资源库弹窗中被看到。
  - 检查切换“全局默认上传类型”后，文章页和批量新增页是否同步变化。

建议最少验证以下接口：

```bash
curl http://127.0.0.1:8080/api/v1/healthz
curl -X POST http://127.0.0.1:8080/api/v1/auth/login
curl -H "Authorization: Bearer <token>" http://127.0.0.1:8080/api/v1/auth/me
curl -H "Authorization: Bearer <token>" http://127.0.0.1:8080/api/v1/admins
```

## 生产环境部署

下面给出推荐的单机部署方案：

- 后端：`systemd` 托管二进制
- 前端：`Nginx` 静态资源
- API：`Nginx` 反向代理到后端 `127.0.0.1:8080`

### 1. 服务器准备

建议准备：

- Linux 服务器
- 已安装 MySQL、Redis、Nginx
- 已创建部署目录 `/opt/go_sleep_admin`
- 运行用户建议使用 `www-data` 或专用低权限用户

建议目录结构：

```text
/opt/go_sleep_admin/
├── backend/
│   ├── bin/sleep-admin
│   ├── config/
│   │   ├── env/prod.yml
│   │   └── pem/*.pem
│   ├── logs/
│   └── public/uploads/
└── frontend/
    └── dist/
```

### 2. 后端生产配置

编辑：

```bash
backend/config/env/prod.yml
```

重点检查：

```yaml
application:
  domain: "https://admin.example.com"
  host: "127.0.0.1"
  port: "8080"
  mode: "release"

database:
  host: "127.0.0.1"
  port: "3306"
  name: "sleep_admin"
  username: "your_user"
  password: "your_password"

redis:
  addr: "127.0.0.1:6379"
  password: "your_password"

log:
  level: "info"
  path: "./logs/sleep_admin"

upload:
  driver: "local"   # 或 oss / cos / obs
```

生产环境重要说明：

- `prod` 启动时不会重置已有内置账号、角色、菜单，只补缺失数据。
- 但仍会执行 `autoMigrate` 和权限规则同步，升级前请先备份数据库。
- JWT PEM 必须替换为你自己的密钥。

### 3. 编译后端

如果在服务器本机编译：

```bash
cd /opt/go_sleep_admin/backend
/usr/local/go/bin/go build -o bin/sleep-admin .
```

如果在 CI 或本地交叉编译：

```bash
cd backend
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 /usr/local/go/bin/go build -o sleep-admin .
```

然后将二进制上传到：

```text
/opt/go_sleep_admin/backend/bin/sleep-admin
```

### 4. 构建前端

如果前端部署在域名根路径 `/`，推荐把：

```bash
frontend/.env.production
```

中的：

```env
VITE_PUBLIC_PATH = '/system/'
```

改成：

```env
VITE_PUBLIC_PATH = '/'
```

再执行构建：

```bash
cd frontend
pnpm install
pnpm build:prod
```

将产物上传到：

```text
/opt/go_sleep_admin/frontend/dist
```

如果你确实想部署到子路径 `/system/`，那就保留 `VITE_PUBLIC_PATH=/system/`，并同步修改 Nginx location。

### 5. 配置 systemd

当前仓库示例文件：

```text
deploy/systemd/sleep-admin.service
```

安装方式：

```bash
sudo cp deploy/systemd/sleep-admin.service /etc/systemd/system/sleep-admin.service
sudo systemctl daemon-reload
sudo systemctl enable sleep-admin
sudo systemctl start sleep-admin
```

常用命令：

```bash
sudo systemctl status sleep-admin
sudo systemctl restart sleep-admin
sudo journalctl -u sleep-admin -f
```

### 6. 配置 Nginx

当前仓库示例文件：

```text
deploy/nginx/sleep-admin.conf
```

安装方式：

```bash
sudo cp deploy/nginx/sleep-admin.conf /etc/nginx/conf.d/sleep-admin.conf
sudo nginx -t
sudo systemctl reload nginx
```

示例配置做了三件事：

- `/`：前端静态页面
- `/api/`：反向代理到 `127.0.0.1:8080`
- `/uploads/`：本地上传目录映射

如果生产环境使用 OSS/COS/OBS：

- 可以保留 `/uploads/` location 但不会成为主要访问路径
- 更推荐直接通过云存储公网地址访问资源

### 7. 生产环境验收

后端健康检查：

```bash
curl http://127.0.0.1:8080/api/v1/healthz
```

Nginx 外部检查：

```bash
curl https://admin.example.com/api/v1/healthz
```

登录检查：

```bash
curl -X POST https://admin.example.com/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin"}'
```

上线前务必处理：

- 修改默认管理员密码
- 替换 JWT PEM 密钥
- 使用生产数据库账号与 Redis 密码
- 调整 `prod.yml` 的 `domain`、`log.path`、`upload.*`
- 做数据库备份

### 8. 升级与回滚建议

建议升级顺序：

1. 备份数据库
2. 上传新后端二进制
3. 上传新前端 `dist`
4. 检查 `prod.yml`
5. `systemctl restart sleep-admin`
6. 做健康检查和登录检查

建议回滚策略：

1. 保留上一版后端二进制
2. 保留上一版前端 `dist`
3. 升级前先做数据库备份
4. 如果新版本异常，先回滚二进制和前端资源，再评估数据库是否需要回滚

## 常用命令

### 后端

```bash
cd backend
/usr/local/go/bin/go vet ./...
/usr/local/go/bin/go test ./...
/usr/local/go/bin/go test -race ./...
/usr/local/go/bin/go build ./...
```

### 前端

```bash
cd frontend
pnpm install
pnpm dev
pnpm build:prod
```

## 当前架构状态说明

当前后端已经完成运行时骨架重构，但为了平滑迁移，仍保留了部分 legacy 三层代码：

- `internal/handler`
- `internal/service`
- `internal/repository`

这些层现在已经不再由 `router` 直接创建，而是统一由 `bootstrap/providers.go` 进行装配，再通过 `internal/module/*` 以 `router.Module` 方式挂到路由层。

这意味着：

- 目录层级已经完成新的运行时架构切换。
- 业务代码还能继续稳定运行。
- 后续如果要进一步纯化架构，可以继续把 legacy handler/service/repository 下沉到各自 `module/*` 内部。 
