package handle

import (
	"context"
	"github.com/JoyZF/go-micro-kit/proto/greeter"
	"github.com/JoyZF/go-micro-kit/service/greeter/internal/logic"
)

type Greeter struct {
}

func (g *Greeter) Hello(ctx context.Context, req *greeter.Request, rsp *greeter.Response) error {
	return logic.NewGreeterLogic().Hello(ctx, req, rsp)
}

func (g *Greeter) Set(ctx context.Context, req *greeter.SetRequest, rsp *greeter.ComRsp) error {
	return logic.NewGreeterLogic().Set(ctx, req, rsp)
}
