# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 仓库性质

这是商城前台 + 管理后台 + 代理工作台 + 主站后端的**脱敏开源版本**（MIT）。上游供应商和聚合支付通道已匿名化，接入地址统一用保留域名 `example.invalid`，商户身份与密钥不随源码提供。真实支付、采购、邮件、机器人默认关闭。

因此：**不要**把匿名域名、演示商品映射当成可用配置，也不要"补上密钥就能上线"。新增代码要维持匿名化边界，不能引入真实商户标识、生产域名或云资源实值。`npm run security:scan` 会拒绝未知外联地址和云资源实值。

项目协作规则（含 Git 提交要求）见 [AGENTS.md](AGENTS.md)；发布检查清单见 [SECURITY.md](SECURITY.md)；"商品为什么显示暂停售卖"见 [services/storefront/docs/INTEGRATION.md](services/storefront/docs/INTEGRATION.md)。

### 真实渠道档案在哪

源码里搜不到真实供应商名、真实支付主机名，**这是设计，不是遗漏**。真实身份、地址与密钥不随仓库提供：

- **本机私有清单** `SUPPLIERS.local.md`：真实供应商身份、部署节点路径、真实接入地址、渠道开关现状。该文件已由 `.gitignore` 排除，**绝不提交**。
- **部署节点**：上游供应商接口契约、商品成本映射、渠道开关与密钥路径、运维流程，都在部署主机上（具体路径见上面那份私有清单）。
- **聚合支付 A / B 的真实主机名**：本仓库里已匿名化且**无法从本仓库恢复**，只能从上游原始项目取回，方法见 [services/storefront/docs/INTEGRATION.md](services/storefront/docs/INTEGRATION.md)。


## 两套后端，一套前端

同一份 React 前端，有两个**功能平行但互不共享代码/表结构**的后端：

| | Go 后端（主） | Cloudflare Worker（备用） |
|---|---|---|
| 入口 | `services/storefront/cmd/storefront/main.go` | `worker/index.ts` |
| 存储 | PostgreSQL | D1 |
| 迁移 | `services/storefront/migrations/`（`NNN_name.sql`，嵌入式） | `migrations/`（`NNNN_name.sql`） |
| 路由 | `internal/httpapi/*.go`（Gin） | `worker/commerce.ts` |

两者都实现了支付、充值、采购、X 赠送、对账，路由形状与适配器协议概念上对应，但**没有任何代码或 DDL 复用**。改动一个后端的接口时，另一个通常也要改；前端契约由 `shared/` 保证。README 把 Worker 称作"备用实现"。

## 常用命令

### 前端（仓库根目录）

```sh
npm ci
npm run dev              # Vite 5173，/api 代理到 Worker 8787
npm test                 # 先 emails:build，再 node --test tests/*.test.ts（Node 24 原生 TS + node:sqlite）
npm run lint             # oxlint
npm run build            # 见下方"构建流水线"
npm run security:scan    # python3 scripts/scan-public-source.py
```

接 Go 后端时换用 `npx vite --config services/storefront/vite.config.ts`（5175 → 8093，标题加 `[测试站]` 前缀）。管理员入口 `/admin`，代理入口 `/partner` 与 `/join`。

单测：`node --test tests/commerce.test.ts`（可再叠加 `--test-name-pattern`）。

### Go 后端（`cd services/storefront`）

```sh
go test ./...
go vet ./...
# 集成测试必须显式给库，否则整包 Skip：
STOREFRONT_TEST_DATABASE_URL='postgres://test_user:PASS@127.0.0.1:5432/storefront_test?sslmode=disable' go test -count=1 ./...
go run ./cmd/storefront migrate && go run ./cmd/storefront grant-runtime
go run ./cmd/storefront api      # 另一个终端：worker
```

单测：`go test ./integration/ -run TestXxx -count=1`。

`main.go` 是子命令分发（`api` / `worker` / `migrate` / `grant-runtime` / `check-runtime` / `health-*` / `import-d1-*` / `telegram-configure`）。**先以所有者身份 `migrate` + `grant-runtime`，再用受限角色跑 `api`/`worker`**；`api`/`worker` 启动时会 `VerifyRuntime` 拒绝超级用户和迁移所有者连接。

### Cloudflare 本地

```sh
cp .dev.vars.example .dev.vars
npx wrangler d1 migrations apply gptibo-local --local
npm run preview:worker
```

## 构建流水线

`npm run build` 依次执行：`emails:build` → `tsc -b` → `tsc -p tsconfig.worker.json` → `vite build` → `node scripts/stamp-csp-hash.mjs`。

两个容易踩的点：

- **邮件模板是生成物**。`emails/*.tsx`（react-email）由 `scripts/build-emails.tsx` 渲染，产物写入 `worker/emails.generated.ts`、`services/storefront/internal/mailer/templates.json`，并断言模板 token 不变量。改邮件要改 `.tsx` 源，不要手改生成文件。
- **`stamp-csp-hash.mjs` 不能跳过**。`index.html` 里有一段内联首屏 bootstrap 脚本，部署 CSP 没有 `'unsafe-inline'`，因此每次构建都把该脚本的 SHA-256 重新写进 `dist/_headers` 的 `script-src`。改了 `index.html` 的内联脚本而没跑这一步，线上会被 CSP 拦掉。

### 前端契约测试

Go 集成测试 `integration/frontend_test.go` 会调 `services/storefront/scripts/check-frontend-contracts.mjs`，用 TypeScript 编译器扫 `src/` 里 `api()` / `post()` / `useResource()` / `fetch()` 的字面量路径，和 Go 的 `RouteInventory()` 对比。**所以改 Go 路由或前端请求路径时要同时更新另一边**，否则这条测试会挂。跑它需要先在根目录 `npm ci`。

## 架构要点

### 前端（`src/`）

`src/Root.tsx` 按 pathname 手工分流三个界面，无路由库：`/admin*` → `admin/Admin.tsx`，`/join`、`/partner` → `partner/Partner.tsx`，其余 → `App.tsx`（包在 `StorefrontI18nProvider` 里）。`Admin.tsx` / `Partner.tsx` 是 38KB / 48KB 的单文件工作台，内部用 tab 数组组织，面板拆分到同目录文件。

`src/lib/` 是前端共享层：`api.ts`（fetch 封装、超时、401/403 事件）、`use-store.ts` / `use-session.ts` / `use-resource.ts`（状态）、`storefront-*.ts`（i18n 与路由解析）、`request-lane.ts`、`visible-polling.ts`。

样式用 Tailwind 4 + shadcn（`components.json`）。**不要引入 shadcn CLI**：仓库只保留 `src/styles/vendor/shadcn-tailwind.css` 静态样式，因为生成 CLI 的依赖链带未修补的 `braces` 拒绝服务漏洞（见该目录 README）。新增组件请沿用 `src/components/ui/`。

### `shared/` —— 跨运行时契约层

仓库根的 `shared/*.ts` 是纯类型、常量和无副作用的纯函数，**被前端（`src/`）、Worker（`worker/`）和测试同时 import**。这是保持 Worker DTO 与 UI 类型对齐的唯一机制。涉及跨端字段（订单状态、商品映射、分账、发票、公告、支付窗口、xgift 签名）的改动，先看这里有没有对应模块，并同步两端调用点。

### Worker（`worker/`）

`worker/index.ts` 导出 `fetch` + `scheduled`。`fetch` 固定顺序：`requireSecureTransport` → `limitApi` → 路由分发（webhook 优先，然后交给 `commerce.ts`）→ 全部响应过 `secureResponse`（CSP、nosniff、HSTS、`/api/` 强制 `no-store`）。`scheduled` 跑各类对账与重试。

- `worker/platform.ts` 定义 `Env`（所有绑定与密钥）、`Statement`/`Database` 抽象、`ApiError`、`json`、`hash`、cookie 工具——Worker 运行时地基。
- `worker/commerce.ts` 是 API 总路由，按前缀分发到 `partners.ts` / `invoices.ts` / `customers.ts` / `accounting.ts` / `refunds.ts` / 各供应商服务。
- 适配器目录 `account_supplier/`、`delivery_supplier/`、`payments/`、`recharge/`、`xgift/` 统一是 `client.ts`（传输与签名）+ `service.ts`（业务与对账）。Go 侧有对应的 `internal/` 包。

`tests/` 用 `DatabaseSync(':memory:')` 建 D1，按文件名顺序跑完所有 `migrations/*.sql`，然后直接调 `worker.fetch(request, env)` 打真实 Worker。夹具见 `tests/customer-fixture.ts`、`tests/admin-email-fixture.ts`。测试直接以 `.ts` 扩展名 import 源码。

### Go 后端（`services/storefront/internal/`）

- `store/`：只有 `Store{DB *pgxpool.Pool}` 和连接/迁移/心跳，**没有 repository 抽象**——各业务包直接用 `s.DB` 写 SQL。
- `httpapi/`：`server.go` 的 `Router(...)` 装全局中间件链（请求 ID、安全头、panic 恢复、`trafficProtection`、`domainIngress`、`securityAudit`、管理员 Access 网关、`apiLimit`、`rehearsalGuard`），再由 `identity_management.go` 建 `/api` 分组和各功能 `registerXxx(...)`。鉴权用 cookie（`__Host-sf_<role>`），CSRF 用 Origin 相等判断。
- `worker/`：`Job{Name, Interval, Timeout, Run, Wake}`，每任务一个 goroutine，`Wake` 走 PostgreSQL LISTEN/NOTIFY。
- `config/`：**只读环境变量，不加载 dotenv**；`Secret(name)` 读 `NAME` 或 `NAME_FILE`。`LoadPayments` / `LoadRecharge` / `LoadProcurement` / `LoadXGift` / `LoadDomains` 等按集成分别加载，未配置即返回禁用。
- `fault/`：统一错误。**不要把请求体、cookie、路径或数据库错误写进日志。**

### 迁移

两套迁移都**只面向全新安装**，不兼容任何既有生产库（`store/migration_compat.go` 明确说明不带任何私有部署的哈希例外）。

Go 侧：迁移文件通过 `migrations/embed.go` 嵌入，按 `NNN_name.sql` 排序在**单个事务**内执行，带 `pg_advisory_xact_lock` 防并发，校验和存 `storefront_migrations`。**已应用的迁移文件内容一旦改动就会导致启动失败**；数据库比发布版本新也直接报错；首次迁移要求 schema 为空。改迁移请新增文件，不要改旧文件。

D1 侧：纯 SQL，文件名零填充顺序执行，无校验和。

### 领域约定

- **金额一律整数分**（`price_cents bigint`，带 CHECK 约束）；费率用 `payments.FeeCents(base, bps)`。不要引入浮点金额。
- **幂等**：`orders.idempotency_key` 按 `owner_hash` 唯一；调用供应商带 `Idempotency-Key`；提交结果不明时按同一 key 对账。
- **订单状态机**：`status` / `payment_status` + 乐观锁 `version` + `order_events` 审计行。
- **集成测试夹具**（`integration/storefront_test.go` 的 `fixture`）要求 `STOREFRONT_TEST_DATABASE_URL` 指向 **loopback 且库名以 `_test` 结尾**的库，否则直接 `t.Fatal`；它会建随机 schema 与受限角色并在结束时清理，所以测试库所有者需要有建角色权限。**绝不能指向业务库。**

### 商品可售性

后台把套餐设为"启用"**不等于**前台可下单。真正决定的是 `internal/procurement/service.go` 的 `Enabled(plan)`：`EmailEnabled`（`EMAIL_PRODUCTS_ENABLED`）、`PaymentReady`（通道可用 + `PAYMENT_CHECKOUT_ENABLED` + `PAYMENT_NOTIFY_ORIGIN`）、`Delivery.Sender` 三道公共门槛，再加各商品自己的条件。"后台显示启用、前台显示暂停售卖"是预期行为，不是 bug。

## 部署

```sh
npm run build
ANX_HOST=<ssh 别名或 user@host> sh services/storefront/scripts/deploy.sh <release>
```

**用 `sh` 调用而不是 `bash`**：Windows 上 `bash` 常指向 WSL，会换一套 ssh/tar 而连不上目标机；脚本本身是 POSIX sh。脚本只传源码和 `dist`，排除 `.secrets`、`.env` 与构建产物，因此不会覆盖目标机配置。目标机上的 `scripts/rollout.sh` 负责重建镜像、迁移、同步 nginx 安全头、启动并健康检查。

`wrangler.jsonc` 不绑定任何生产账号或域名，D1 的 `database_id` 是全零占位值；仓库没有一键生产部署命令。
