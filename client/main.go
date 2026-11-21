package main

import (
	"context"
	"fmt"
	"os"

	"github.com/oekit/core/config"
	kitexotel "github.com/oekit/kitex/observability/tracing/instrumentation/kitexotel"
	registrykitex "github.com/oekit/kitex/registry/kitex"
	"github.com/oekit/logx"
	kitexlog "github.com/oekit/logx/kitexlog"
	"github.com/oekit/registry/nacos"
	nacosimpl "github.com/oekit/registry/nacosimpl"
	"github.com/oekit/tracing"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	echo "github.com/oekit/kitex-examples/kitex_gen/kitex/echo"
	"github.com/oekit/kitex-examples/kitex_gen/kitex/echo/echoservice"
)

type AppConfig struct {
	ServiceName  string          `yaml:"service_name"`
	TargetServer string          `yaml:"target_server"`
	TargetAddr   string          `yaml:"target_addr"` // mock 模式直连地址
	Nacos        nacos.Bootstrap `yaml:"nacos"`
	Tracing      tracing.Config  `yaml:"tracing"`
}

func main() {
	// run() 返回错误再统一退出：tracing provider 的 defer Shutdown 必须
	// 真正执行（gocritic exitAfterDefer），错误路径也不在深栈里散落 os.Exit。
	if err := run(); err != nil {
		logx.Error("client exited with error", logx.ErrKey, err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadBootstrap[AppConfig]("dev")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// AGENTS 日志节硬规则（#112）：业务服务的包级 logx.* 全部经 Default()
	// 输出——默认 logger 必须挂 service 字段，否则每条日志都无 service 归属。
	// （下方 klog 桥接用的是独立 slog 实例，两条接线在此各自示范。）
	logx.UseSlog(cfg.ServiceName)

	// 1. 日志桥接（必须传 slog 实例避免无限递归）
	// kitex.NewLogger 已弃用（v0.x 起桥接至 logx/kitexlog），示例直接用新 API
	klog.SetLogger(kitexlog.NewLogger(logx.NewSlogLogger(cfg.ServiceName)))

	// 2. 初始化 tracing（#175：此前挂 ClientTracingSuite 但从不 Init——OTel
	// 全局 provider 为 noop，client span 静默丢弃）。
	prov, err := tracing.Init(cfg.Tracing)
	if err != nil {
		return fmt.Errorf("tracing init: %w", err)
	}
	defer func() { _ = prov.Shutdown(context.Background()) }()

	// 3. 创建 Kitex client（mock 直连 / NNS 发现二选一）
	co := []client.Option{client.WithSuite(kitexotel.NewClientTracingSuite())}
	if cfg.Nacos.Addr == "" {
		logx.Info("mock 模式：nacos.addr 为空——直连 target_addr", "addr", cfg.TargetAddr)
		co = append(co, client.WithHostPorts(cfg.TargetAddr))
	} else {
		namingCli, nerr := nacos.NewNamingClient(cfg.Nacos)
		if nerr != nil {
			return fmt.Errorf("create naming client: %w", nerr)
		}
		disc, derr := nacosimpl.NewDiscoverer(namingCli)
		if derr != nil {
			return fmt.Errorf("create discoverer: %w", derr)
		}
		co = append(co, client.WithResolver(registrykitex.NewResolver(disc)))
	}

	cli, err := echoservice.NewClient(cfg.TargetServer, co...)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}

	// 4. 调用
	resp, err := cli.Echo(context.Background(), &echo.EchoRequest{
		Message: "hello from ekit",
	})
	if err != nil {
		return fmt.Errorf("echo call: %w", err)
	}

	logx.Info("response received", "message", resp.Message, "timestamp", resp.Timestamp)
	return nil
}
