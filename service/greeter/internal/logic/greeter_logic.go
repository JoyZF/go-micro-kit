package logic

import (
	"context"
	"errors"
	"github.com/JoyZF/go-micro-kit/proto/greeter"
	"github.com/JoyZF/go-micro-kit/service/greeter/internal/model"
	"time"
)

type GreeterLogic struct {
}

func NewGreeterLogic() *GreeterLogic {
	return &GreeterLogic{}
}

func (g *GreeterLogic) Hello(ctx context.Context, req *greeter.Request, rsp *greeter.Response) error {
	if req.GetName() == "" {
		return errors.New("name is empty")
	}
	rsp.Greeting = req.GetName()
	return nil
}

func (g *GreeterLogic) Set(ctx context.Context, req *greeter.SetRequest, rsp *greeter.ComRsp) error {
	err := model.GetRdb().Set(ctx, req.GetKey(), req.GetValue(), 86400*time.Second).Err()
	if err != nil {
		return err
	}
	rsp.Code = 1
	rsp.Msg = "success"
	rsp.Data = ""
	return nil
}
