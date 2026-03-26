# Config Layout

项目配置代码与示例统一放在 `config/` 目录：

- `config/config.go`
  嵌入式配置入口，负责提供 `config/env/*.yml`。
- `config/env/dev.yml`
  开发环境配置。
- `config/env/test.yml`
  测试环境配置。
- `config/env/prod.yml`
  生产环境配置。
- `internal/platform/config/load.go`
  运行时 typed config 加载与校验。

运行时加载顺序：

1. 显式传入的 `config/env/${ENV}.yml`
2. 若磁盘文件不存在，则回退到 `config/config.go` 中 embed 的同名 yml

通常建议：

- 直接维护 `config/env/*.yml`
- 通过 `server -c dev|test|prod` 选择环境
- 不再使用 `config/env/.env` 这类旧 env 文件
