package gateway

import (
	"context"
	"go-micro.dev/v4/api"
	"go-micro.dev/v4/util/log"
)

func Gateway() {
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		srv := api.NewApi()

		if err := srv.Run(ctx); err != nil {
			log.Fatal(err)
		}
	}()
}
