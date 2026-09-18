# 项目命名与文档边界

## 当前项目身份

本仓库的项目名称为 MeterForge，Git 远程地址为 `https://github.com/Pototoooo/meterforge.git`。项目身份、功能范围和代码来源是不同问题；更改名称不能证明代码的原创归属。

README、项目源码和测试用于说明本地功能。外部网站的文档快照属于学习参考，不能据此推断本地功能、维护责任或整个项目的来源。

## 外部参考资料

原 `docs/openmeter-reference-20260909/` 是 OpenMeter 官方文档的本地快照，不是 MeterForge 的产品手册。本次清理将其从项目文档树移出，并在清理交付包中以 `EXTERNAL_REFERENCE.tar.gz` 保存全部 251 个文件，未改写快照中的来源、作者说明或内容。此前的 Git 版本也保留这些文件。

需要阅读时，将归档解压到仓库之外。不要将外部文档中的项目名、服务地址和署名替换成本项目的名称后重新发布。

## 保留的名称与原因

| 位置 | 保留内容 | 原因 |
| --- | --- | --- |
| `go.mod` | `github.com/openmeterio/run` | 它是 `github.com/oklog/run` 的实际替换依赖，不是本项目模块名。移除需要验证替换版本的行为差异。 |
| `go.sum` | 该依赖的两条校验记录 | 与真实依赖版本对应，不能通过字符串替换生成新的校验记录。 |
| `tools/migrate/migrations/20250605131637_migrate-flat-fees-to-ubp-flat-fees.up.sql` | `/openmeter-line-reason` | 已有历史迁移写入的数据键；保持历史文件不变。 |
| `tools/migrate/migrations/20250731141420_billing-migrate-flat-fee-lines.up.sql` | `/openmeter-line-reason` | 同上；若迁移存量数据，另行新增并验证迁移。 |
| `tools/migrate/migrations/20260527120000_dedupe_tax_codes_by_app_mapping.up.sql` | 注释中的旧代码路径 | 历史迁移由校验清单管理，不为修改注释而改写历史文件。 |

## Git 与来源记录

清理工作使用 `codex/meterforge-naming-cleanup` 分支。现有远程地址、分支及标签名称中未发现需要替换的 OpenMeter 或 Weknora 名称。

本次不改写已有提交标题、作者、日期、标签对象或历史文件内容。旧文档中的来源声明尚未完成逐文件核实，本次既不将其作为已确认结论，也不把命名清理当作对它的否定。根目录缺少 LICENSE 的问题不通过虚构新许可或删除第三方声明解决，需单独核对适用范围。

## 检查边界

本次对 Git 跟踪的普通文本文件及文件名进行命名检查。`.work`、`node_modules` 中的依赖与构建副本不作为项目身份判据；它们已经被跟踪的文件不会因加入忽略规则自动移出索引。二进制文件不做字符串替换，未跟踪的本地缓存不删除。

本次只修改项目文档和隔离外部参考快照，不改变业务逻辑、API、配置键或数据库行为。
