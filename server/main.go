package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/oekit/bootstrap"
	bootstrapnacos "github.com/oekit/bootstrap/nacos"
	bootstrapobs "github.com/oekit/bootstrap/obs"
	"github.com/oekit/core/config"
	kitexgo "github.com/oekit/kitex"
	kitexotel "github.com/oekit/kitex/observability/tracing/instrumentation/kitexotel"
	registrykitex "github.com/oekit/kitex/registry/kitex"
	"github.com/oekit/logx"
	kitexlog "github.com/oekit/logx/kitexlog"
	"github.com/oekit/registry/nacos"
	nacosimpl "github.com/oekit/registry/nacosimpl"
	"github.com/oekit/tracing"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	echo "github.com/oekit/kitex-examples/kitex_gen/kitex/echo"
	"github.com/oekit/kitex-examples/kitex_gen/kitex/echo/echoservice"
)

type AppConfig struct {
	ServiceName string          `yaml:"service_name"`
	RpcPort     int             `yaml:"rpc_port"`
	Nacos       nacos.Bootstrap `yaml:"nacos"`
	Tracing     tracing.Config  `yaml:"tracing"`
}

func main() {
	cfg, err := config.LoadBootstrap[AppConfig]("dev")
	if err != nil {
		logx.Error("load config failed", "err", err)
		os.Exit(1)
	}

	b, err := bootstrap.New(
		bootstrap.WithServiceName(cfg.ServiceName),
		bootstrapnacos.WithNNS(cfg.Nacos, cfg.RpcPort, nil), // Nacos 注册（嵌套 module 按需注入）
		bootstrapobs.WithTracing(cfg.Tracing),
	)
	if err != nil {
		logx.Error("bootstrap init failed", "err", err)
		os.Exit(1)
	}

	if err := b.Run(func(app *bootstrap.App) error {
		// 1. 日志桥接（必须传 slog 实例避免无限递归）
		// kitex.NewLogger 已弃用（v0.x 起桥接至 logx/kitexlog），示例直接用新 API
		klog.SetLogger(kitexlog.NewLogger(logx.NewSlogLogger(cfg.ServiceName)))

		// 2. 创建 Nacos naming client → ekit Registrar → Kitex Registry
		namingCli, err := nacos.NewNamingClient(cfg.Nacos)
		if err != nil {
			return err
		}
		reg, err := nacosimpl.NewRegistrar(namingCli)
		if err != nil {
			return err
		}

		// 3. 创建 Kitex server
		// 真 codegen:业务实现 EchoService 接口(见下方 echoHandler)
		handler := &echoHandler{}
		addr := &net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: cfg.RpcPort}
		svr := echoservice.NewServer(handler,
			server.WithRegistry(registrykitex.NewRegistry(reg)),
			server.WithServiceAddr(addr),
			server.WithSuite(kitexotel.NewServerTracingSuite()),
		)

		// 4. 注册到 bootstrap 生命周期
		// WithOnReady + WithReadyAddr：TCP 探活通过后触发 app.Ready()——
		// bootstrap/nacos v0.3.0 起服务注册延迟到就绪门之后（server 可建连
		// 才进注册中心），关停时反注册先于 server 关闭执行
		kitexgo.RunServer(app, "kitex-echo", svr,
			kitexgo.WithOnReady(func(ctx context.Context) { app.Ready() }),
			kitexgo.WithReadyAddr(fmt.Sprintf(":%d", cfg.RpcPort)))

		// 阻塞等待退出信号:必须阻塞在 app ctx 上——RunServer 启动失败时
		// 的 RequestShutdown 只 cancel ctx,依赖 userFn 从 <-Done() 返回才
		// 能触发优雅关停;select {} 会让进程无服务地"假活"永不退出
		<-app.Context().Done()
		return nil
	}); err != nil {
		logx.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}

// echoHandler 实现 IDL 定义的 EchoService。
type echoHandler struct{}

// Echo 回显消息并附时间戳。
func (h *echoHandler) Echo(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	logx.Info("echo request", "message", req.Message)
	return &echo.EchoResponse{
		Message:   "echo: " + req.Message,
		Timestamp: time.Now().Unix(),
	}, nil
}
