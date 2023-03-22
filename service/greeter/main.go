package main

import (
	"context"
	"github.com/JoyZF/go-micro-kit/proto/greeter"
	"github.com/JoyZF/go-micro-kit/service/greeter/internal/conf"
	"github.com/JoyZF/go-micro-kit/service/greeter/internal/handle"
	"github.com/JoyZF/go-micro-kit/service/greeter/internal/middleware"
	"github.com/JoyZF/go-micro-kit/service/greeter/internal/model"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/go-micro/plugins/v4/server/grpc"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"go-micro.dev/v4/util/log"
	"sync"
)

var wg sync.WaitGroup

func main() {
	var c conf.Config
	conf.InitConfig(&c)
	// etcd registry
	reg := etcd.NewRegistry(registry.Addrs(c.App.Endpoint))
	// init config
	service := micro.NewService(
		micro.Name(c.App.Name),
		micro.Registry(reg),
		micro.WrapHandler(waitgroup(&wg)),
		// waits for the waitgroup once stopped
		micro.AfterStop(func() error {
			// wait for handlers to finish
			wg.Wait()
			return nil
		}),
		micro.WrapHandler(logWrapper),
		// heartbeat
		//micro.RegisterTTL(time.Second*time.Duration(c.App.RegisterTTL)),
		//micro.RegisterInterval(time.Second*time.Duration(c.App.RegisterInterval)),
		micro.Version(c.App.Version),
		// 2 UNKNOWN: No status received
		// # https://github.com/go-micro/go-micro/issues/2534
		// # https://github.com/go-micro/plugins/tree/main/v4/server/grpc
		micro.Server(grpc.NewServer()),
	)
	model.InitMySQL(&c.MySQL)
	model.InitRedis(&c.Redis)
	// NSQ consumer
	//consumer.NewNSQConsumer(
	//	c.NSQ.Topic,
	//	c.NSQ.Channel,
	//	c.NSQ.Addr).
	//	Consumer(&nsqHandler.Greeter{})
	service.Init(
		// middleware
		micro.WrapHandler(middleware.NewAuthWrapper(service)),
	)

	if err := greeter.RegisterGreeterHandler(service.Server(), new(handle.Greeter)); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

// waitgroup is a handler wrapper which adds a handler to a sync.WaitGroup
func waitgroup(wg *sync.WaitGroup) server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			wg.Add(1)
			defer wg.Done()
			return h(ctx, req, rsp)
		}
	}
}

// logWrapper is a handler wrapper
func logWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		log.Infof("[wrapper] server request: %v", req.Endpoint())
		err := fn(ctx, req, rsp)
		return err
	}
}
