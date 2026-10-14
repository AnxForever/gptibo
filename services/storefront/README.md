# Anx Go 后端

Go + Gin + PostgreSQL，API 与业务 Worker 使用同一个二进制。此版本仅面向全新安装；匿名化后的迁移与原部署不兼容，不导入生产快照。

## Docker 本地启动

需要 Docker Compose，命令在本目录执行。首先生成三个独立随机秘密（Python 3）：

```sh
python3 scripts/init-secrets.py
cp .env.example .env
docker compose config --quiet
docker compose build api
docker compose up -d --wait postgres
docker compose --profile maintenance run --rm maintenance migrate
docker compose --profile maintenance run --rm maintenance grant-runtime
docker compose up -d --wait api worker
```

也可以在 PowerShell 使用 `./scripts/init-secrets.ps1` 初始化。秘密文件位于被忽略的 `.secrets/`，已有文件不会覆盖。API 仅监听宿主机 `127.0.0.1:8093`，PostgreSQL 不映射宿主机端口；运行数据库角色与迁移所有者分离。

`compose.auth.yaml` 可选启用管理员/邮件：配置 `ADMIN_USERNAME`、`ADMIN_EMAIL`、`MAIL_FROM`，自行生成 `.secrets/admin-totp` 并保存 `.secrets/resend-key`，再叠加该 Compose 文件。不要使用测试中的固定认证向量。

其他接入参数见 `.env.integrations.example`。复制为 `.env.integrations` 后可叠加 `compose.integrations.yaml`；所有开关默认关闭，匿名适配器的保留域名无法用于真实交易。必须先按自己的上游协议完成开发和测试，不能直接填生产秘密试运行。可选域名配置在 `compose.domains.yaml`，仅使用自己新建的云资源。

`compose.turnstile.yaml` 可选启用 Cloudflare Turnstile：把 widget 注册表写入 `.secrets/turnstile`（JSON 数组，每项含 `sitekey`、`secret`、`hosts`），再叠加该 Compose 文件。未登记的主机名一律拒绝，不会静默放行。

## 部署到服务器

把前端构建产物与后端源码推送到目标机，再由目标机上的 `scripts/rollout.sh` 完成滚动更新：重建镜像、执行迁移、按 `dist/_headers` 同步 nginx 的安全头、启动并做健康检查。

```sh
npm run build
ANX_HOST=korea-seoul sh services/storefront/scripts/deploy.sh 2026.10.14
```

用 `sh` 而不是 `bash` 调用：Windows 上 `bash` 常指向 WSL，那会用到另一套 ssh/tar 而连不上目标机。脚本本身是 POSIX sh。

未改动源码时命中构建缓存，数秒完成；改动 Go 源码后需要重新编译，约 1–2 分钟。该脚本不传输 `.secrets`、`.env` 与构建产物，因此不会覆盖目标机上的配置。

## 原生开发

程序读取环境变量，不自动加载 dotenv。通过 `DATABASE_URL_FILE` 或 `DB_HOST/DB_PORT/DB_NAME/DB_USER/DB_PASSWORD_FILE` 指向独立数据库；外部连接默认要求 TLS。先用所有者身份执行 `migrate` 和 `grant-runtime`，再用受限角色执行 `api` 和 `worker`：

```sh
go run ./cmd/storefront migrate
go run ./cmd/storefront grant-runtime
go run ./cmd/storefront api
# 另一个终端使用相同运行配置：
go run ./cmd/storefront worker
```

`grant-runtime` 还需 `APP_DB_PASSWORD_FILE`。运行身份需 `AUTH_SECRET_FILE`，默认公开来源为 `http://127.0.0.1:8093`。配置管理员与邮件后才能使用完整后台登录。根目录执行 `npx vite --config services/storefront/vite.config.ts` 启动前端。

## 测试

先在仓库根目录执行 `npm ci`；前后端路由契约测试需要 TypeScript 解析器。

```sh
go test ./...
go vet ./...
# 专用本机 PostgreSQL，库名必须以 _test 结尾：
STOREFRONT_TEST_DATABASE_URL='postgres://test_user:YOUR_LOCAL_PASSWORD@127.0.0.1:5432/storefront_test?sslmode=disable' go test -count=1 ./...
```

每个集成测试建立随机 schema 和受限角色，并在结束后清理，需要测试数据库所有者拥有创建角色权限。不可指向业务数据库；没有连接配置时，集成测试明确跳过。
