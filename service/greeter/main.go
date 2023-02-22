package main

import (
	"github.com/JoyZF/go-micro-kit/proto/greeter"
	"github.com/JoyZF/go-micro-kit/service/greeter/internal/conf"
	"github.com/JoyZF/go-micro-kit/service/greeter/internal/handle"
	"github.com/JoyZF/go-micro-kit/service/greeter/internal/model"
	"go-micro.dev/v4"
	"go-micro.dev/v4/util/log"
)

func main() {
	var c conf.Config
	conf.InitConfig(&c)
	// init config
	service := micro.NewService(
		micro.Name(c.App.Name),
	)
	model.InitMySQL(&c.MySQL)
	model.InitRedis(&c.Redis)
	service.Init()

	if err := greeter.RegisterGreeterHandler(service.Server(), new(handle.Greeter)); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
