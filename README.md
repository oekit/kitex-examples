# kitex-examples

> Kitex echo 示例——client/server + Thrift IDL + 治理配置样例。

> 项目主页：**https://oekit.github.io/** · 名称由来：**O**pen **E**asy **Kit**（开源 · 易用 · 工具集合）
> 🤖 AI Agent：本仓规范与硬规则见 [AGENTS.md](AGENTS.md)；机器可读索引见 [llms.txt](https://oekit.github.io/llms.txt)（接口签名以源码为准）。

本仓库是 [ekit](https://git.enjoye.top/enjoydream/ekit) 拆分矩阵的成员（Round-3 自 ekit-kitex 迁出），
使 [ekit-kitex](https://github.com/oekit/kitex) 主模块不再依赖 bootstrap 全链。

## 结构

- `kitex_gen/` — thrift 生成的 echo 服务代码
- `client/`、`server/` — 示例客户端与服务端（接入 bootstrap + Nacos + 治理配置）
- `idl/echo.thrift` — IDL
- `configs/` — retry/rpc_timeout/limit 治理配置样例

## 运行

```bash
go run ./server    # 先起 server（注册到 Nacos / mock 模式）
go run ./client
```

依赖：[ekit-kitex](https://github.com/oekit/kitex)、[ekit-bootstrap](https://github.com/oekit/bootstrap)、[ekit-registry](https://github.com/oekit/registry)、[ekit-logx](https://github.com/oekit/logx)、[ekit-tracing](https://github.com/oekit/tracing)。

## License

见 [LICENSE](LICENSE)
