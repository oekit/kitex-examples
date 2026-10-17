# Changelog — oekit/kitex-examples

由 ekit v0.34.0（git.enjoye.top/enjoydream/ekit，冻结维护）拆分而来。所有 notable 变更记录如下。

## [v0.1.10] — 2026-10-17

- fix(deps): 内部依赖地板全量对齐最新——bootstrap v0.5.1（nacos/obs v0.2.1）、kitex v0.5.2、core v0.6.0、logx v0.5.0（kitexlog v0.1.1）、registry v0.3.1、tracing v0.2.1；go 维持 1.26.0（kitex/OTel 生态传导）

## [v0.1.9] — 2026-10-17

- bootstrap 结构重置对齐（v0.5.0 起包提升至模块根）：server 示例 import 由旧路径 `bootstrap/bootstrap` 改为模块根；依赖地板 bootstrap v0.5.0、bootstrap/nacos v0.2.0、bootstrap/obs v0.2.0、kitex v0.5.1（旧版本 tag 已随结构重置清空）

## [v0.1.8] — 2026-10-16

- 依赖对齐 kitex v0.5.0（含 istio xDS 动态熔断子包）；示例与矩阵文档同步

## [v0.1.4] — 2026-10-15

- bootstrap v0.4.0 对齐（注入模式）；v0.1.3 为 bootstrap v0.3.0 对齐
## [v0.1.0] — 2026-10-15

自 ekit v0.34.0（git.enjoye.top/enjoydream/ekit，冻结维护）拆分而来的初始公开版本，完成拆分稳定化：模块路径统一为 `github.com/oekit/kitex-examples`，多轮拆分/合并收敛至当前模块边界，组织基建（共享 CI workflow、golangci v2、Dependabot）就位。