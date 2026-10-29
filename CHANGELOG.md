# Changelog — oekit/kitex-examples

由 ekit v0.34.0（git.enjoye.top/enjoydream/ekit，冻结维护）拆分而来。所有 notable 变更记录如下。

## [v0.1.22] — 2026-10-30

- build: 矩阵内依赖地板 `tracing` v0.3.1 → v0.3.2（维护轮全量对齐；本版不含 API 变更）

## [v0.1.21] — 2026-10-24

- 跟随 kitex v0.7.0——tracing suite 改由嵌套模块 kitexotel 提供（kitex 根包 NewServerTracingSuite / NewClientTracingSuite 已迁入 kitexotel；#37）

## [v0.1.20] — 2026-10-22

- build(deps): 嵌套模块地板对齐——bootstrap/nacos v0.3.1 → v0.3.5、bootstrap/obs v0.2.5 → v0.2.9、logx/kitexlog v0.1.1 → v0.1.4（传递地板随 MVS 上抬：bootstrap v0.6.4、tracing v0.2.3、concurrency v0.17.1、pkg v0.15.1、metrics v0.1.8、nns v0.4.3；零代码改动）

## [v0.1.19] — 2026-10-21

- 改用嵌套模块 registry/nacos 与 registry/nacosimpl——registry v0.5.0 起不再提供这两个子包（import 路径不变，由独立 module 提供）

## [v0.1.18] — 2026-10-19

- build: 矩阵内依赖地板对齐（align_floors 不动点）：bootstrap v0.6.1 → v0.6.2、core v0.7.0 → v0.9.2、logx v0.6.0 → v0.7.0

## [v0.1.17] — 2026-10-18

- fix(deps): kitex v0.6.2、tracing v0.2.2、metrics v0.1.7（indirect）地板对齐（零代码改动）

## [v0.1.16] — 2026-10-18

- 示例 import 迁移至扁平化路径（kitex/registry/nns 包提升至模块根）+ ready-gate 延迟注册演示更新；fix(deps): 地板跟随——bootstrap v0.6.1（nacos v0.3.1）、kitex v0.6.1、core v0.7.0、logx v0.6.0、registry v0.4.1、nns v0.4.1（indirect）

## [v0.1.15] — 2026-10-17

- bootstrap v0.5.3（nacos/obs v0.2.5）、kitex v0.5.6、nns v0.3.5 地板跟随

## [v0.1.14] — 2026-10-17

- pkg 地板 v0.11.0；kitex v0.5.5、concurrency v0.5.4 地板跟随

## [v0.1.13] — 2026-10-17

- fix(deps): 地板跟随——kitex v0.5.4、concurrency v0.5.3（pkg 随 kitex 链至 v0.10.0）

## [v0.1.12] — 2026-10-17

- fix(deps): nacos/obs 地板 v0.2.2 → v0.2.3（obs core 直接依赖归位 + bootstrap v0.5.2 地板）

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