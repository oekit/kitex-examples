package main

import (
	"context"
	"os"

	"github.com/oekit/core/config"
	kitexgo "github.com/oekit/kitex"
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
	Nacos        nacos.Bootstrap `yaml:"nacos"`
	Tracing      tracing.Config  `yaml:"tracing"`
}

func main() {
	cfg, err := config.LoadBootstrap[AppConfig]("dev")
	if err != nil {
		logx.Error("load config failed", "err", err)
		os.Exit(1)
	}

	// 1. 日志桥接（必须传 slog 实例避免无限递归）
	// kitex.NewLogger 已弃用（v0.x 起桥接至 logx/kitexlog），示例直接用新 API
	klog.SetLogger(kitexlog.NewLogger(logx.NewSlogLogger(cfg.ServiceName)))

	// 2. 创建 Nacos naming client → ekit Discoverer → Kitex Resolver
	namingCli, err := nacos.NewNamingClient(cfg.Nacos)
	if err != nil {
		logx.Error("create naming client failed", "err", err)
		os.Exit(1)
	}
	disc, err := nacosimpl.NewDiscoverer(namingCli)
	if err != nil {
		logx.Error("create discoverer failed", "err", err)
		os.Exit(1)
	}

	// 3. 创建 Kitex client（带服务发现 + 可观测性）
	cli, err := echoservice.NewClient(cfg.TargetServer,
		client.WithResolver(registrykitex.NewResolver(disc)),
		client.WithSuite(kitexgo.NewClientTracingSuite()),
	)
	if err != nil {
		logx.Error("create client failed", "err", err)
		os.Exit(1)
	}

	// 4. 调用
	resp, err := cli.Echo(context.Background(), &echo.EchoRequest{
		Message: "hello from ekit",
	})
	if err != nil {
		logx.Error("echo call failed", "err", err)
		os.Exit(1)
	}

	logx.Info("response received", "message", resp.Message, "timestamp", resp.Timestamp)
}
