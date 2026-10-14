# ekit-kitex-examples

> Kitex echo 示例——client/server + Thrift IDL + 治理配置样例。

本仓库是 [ekit](https://git.enjoye.top/enjoydream/ekit) 拆分矩阵的成员（Round-3 自 ekit-kitex 迁出），
使 [ekit-kitex](https://github.com/yi-nology/ekit-kitex) 主模块不再依赖 bootstrap 全链。

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

依赖：[ekit-kitex](https://github.com/yi-nology/ekit-kitex)、[ekit-bootstrap](https://github.com/yi-nology/ekit-bootstrap)、[ekit-registry](https://github.com/yi-nology/ekit-registry)、[ekit-logx](https://github.com/yi-nology/ekit-logx)、[ekit-tracing](https://github.com/yi-nology/ekit-tracing)。

## License

见 [LICENSE](LICENSE)
