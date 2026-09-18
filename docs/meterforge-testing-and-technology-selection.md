# MeterForge 独立测试与技术选型实验报告

## 1. 结论

- 被测源码：`main@2439188b3429020810132314a2399fb641ada78c`，未接入 LoreLattice。
- Go 1.25.6 下全量构建通过。
- 根模块完成 6,229 项测试，9 项跳过，4 个失败记录均属于同一个账单迁移套件；失败源于运行时代码调用 `mf_func_migrate_customer_invoices_to_schema_level_2_bulk`，而迁移 SQL 实际创建的是 `om_func_migrate_customer_invoices_to_schema_level_2_bulk`。
- 从当前提交构建的独立 Docker 栈通过默认 E2E（44.701s）和 credits-disabled E2E（4.213s）。
- 不同技术对照改为 Redis 7.4.7 `SET NX` 与 PostgreSQL 14.20 `UNLOGGED` 表 + 唯一主键。15 组成对实验中 Redis 吞吐提升中位数为 **22.76%**，P95/P99 降幅中位数分别为 **14.22%/50.05%**，所有 15 组方向一致，Bootstrap 95% 区间均未跨过 0。
- 两种共享存储都通过双实例重复事件验证；据此选择 Redis 承载高频事件幂等热路径，PostgreSQL 继续承载账单等需要关系约束与持久化的业务数据。

## 2. 功能与工程基线

| 项目 | 结果 |
|---|---|
| `go build -p 4 -tags=dynamic ./...` | 通过，exit 0 |
| 根模块测试 | 6,229 tests，9 skipped，4 failure records，exit 1 |
| 通知子系统（补齐 Svix/Redis） | 通过，exit 0 |
| 默认 E2E | 通过，`ok .../e2e 44.701s` |
| credits-disabled E2E | 通过，`ok .../e2e/creditsdisabled 4.213s` |

### 已定位缺陷

`meterforge/billing/adapter/schemamigration.go` 调用 `mf_func_...`，但 `tools/migrate/migrations/20260717060208_add_bulk_invoice_schema_level_2_migration_function.up.sql` 创建 `om_func_...`。因此根模块暂时不能描述为“全量测试全部通过”。

## 3. 不同技术的控制变量方法

对照对象是两种不同的共享去重技术：

- **Redis 7.4.7**：直接调用 MeterForge 当前 `redisdedupe.Deduplicator` 的 `SET NX` + TTL 路径，使用 XXH3-128 + Base64 键；
- **PostgreSQL 14.20**：隔离实验夹具实现 `UNLOGGED` 表、`event_key TEXT PRIMARY KEY` 与 `INSERT ... ON CONFLICT DO NOTHING`，模拟可替代的共享原子去重后端；没有修改生产源码。

固定条件如下：

- 每轮 100,000 次调用，90,000 个唯一事件和 10,000 个重复事件；
- 并发 32，随机种子 20260902；
- 两种容器都限制为 2 CPU、512 MiB，并使用 tmpfs；
- Redis 关闭 RDB/AOF；PostgreSQL 使用 `UNLOGGED` 表，避免把持久化能力差异混入高频临时去重热路径；
- Go 1.25.6；
- 15 组成对实验，交替执行先后顺序；
- 唯一变化为去重存储与原子写入机制。

统计采用每组 Redis 相对同组 PostgreSQL 的百分比差异，并以固定种子进行 20,000 次成对 Bootstrap，报告中位数 95% 区间。

## 4. Redis 与 PostgreSQL：主要选型证据

| 指标（15 轮中位数） | Redis 7.4.7 | PostgreSQL 14.20 | Redis 成对差异中位数 | Bootstrap 95% 区间 |
|---|---:|---:|---:|---:|
| 吞吐量 | 26,415 ops/s | 21,302 ops/s | **+22.76%** | `[+20.50%, +27.93%]` |
| P50 | 1,067.667 μs | 1,276.125 μs | **-16.75%** | `[-17.49%, -14.40%]` |
| P95 | 2,140.291 μs | 2,486.584 μs | **-14.22%** | `[-18.33%, -11.84%]` |
| P99 | 2,800.292 μs | 6,100.667 μs | **-50.05%** | `[-55.96%, -43.06%]` |
| 正确运行 | 15/15 | 15/15 | 相同 | — |

这里的“降幅区间”在统计脚本中按正值记录，表中为便于阅读展示为负延迟差异。吞吐、P50、P95、P99 的 15/15 组成对结果都支持 Redis 方向，并且 95% 区间没有跨过 0。

两种引擎的空间数字不能直接用于选型百分比：Redis 记录的是进程已分配内存变化，PostgreSQL 记录的是表及索引关系大小；口径不同，因此仅保留原始数据，不写成“节省存储”。

## 5. 双实例正确性

使用两个独立客户端实例并发提交 20,000 次操作，其中包含 10,000 个跨实例重复事件：

| 指标 | Redis | PostgreSQL |
|---|---:|---:|
| 唯一事件 | 10,000 | 10,000 |
| 识别重复事件 | 10,000/10,000 | 10,000/10,000 |
| 错误 | 0 | 0 |
| 吞吐量 | 25,878 ops/s | 21,481 ops/s |
| P95 | 2,138.459 μs | 2,444.667 μs |

两者都满足生产多实例所需的全局原子去重正确性；最终选择依据不是“只有 Redis 正确”，而是在相同正确性下 Redis 的热路径吞吐和尾延迟更优。

## 6. 内存 LRU 与 Redis：部署边界补充

此前同一负载下，内存 LRU 单机中位吞吐为 967,147 ops/s，Redis 为 24,882 ops/s；但两个独立内存 LRU 回放 10,000 个跨实例重复事件时，识别 0/10,000，而 Redis 识别 10,000/10,000。

因此生产多实例选择共享 Redis；本地开发或明确单实例的环境保留内存 LRU。这个对照解释部署边界，Redis 与 PostgreSQL 对照解释共享后端的技术选型。

## 7. Redis 键结构优化

固定 Redis 7.4.7，采用 15 组成对实验比较 raw key 与 XXH3-128 + Base64 哈希键。90,000 个唯一事件下，Redis 增量内存由 14,334,912 B 降至 11,454,912 B，节省 2,880,000 B，即 **20.09%**；15/15 组同方向。吞吐和延迟置信区间跨过 0，因此只陈述内存收益，不声称性能提升。

## 8. 推荐的简历描述

> 针对 AI 用量事件的高频幂等写入，在固定 10 万次请求、32 并发和 2 CPU/512 MiB 资源下，对 Redis `SET NX` 与 PostgreSQL `UNLOGGED` + 唯一索引完成 15 组交替控制变量压测；两者均通过双实例 1 万个重复事件的 100% 去重验证，Redis 吞吐提升中位数 22.76%（95% CI 20.50%–27.93%），P95/P99 分别降低 14.22%/50.05%，据此选择 Redis 承载事件幂等热路径、PostgreSQL 承载持久化业务数据；进一步通过 XXH3-128 哈希键将 9 万事件的 Redis 增量内存降低 20.09%。

## 9. 证据路径

- `root-tests.log`、`root-tests.json`、`root-tests.junit.xml`：根模块测试原始证据；
- `e2e-clean.log`：从当前提交构建的默认和 credits-disabled E2E；
- `technology-selection/raw/`：Redis/PostgreSQL 每轮 JSON 与双实例原始数据；
- `technology-selection/summary.json`：15 组成对统计与 Bootstrap 区间；
- `technology-selection/analyze.py`：可复算统计脚本；
- `technology-selection/VERIFICATION.txt`：夹具差异、精确命令和回滚验证。
