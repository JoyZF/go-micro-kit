package main

import (
	"fmt"
	"github.com/JoyZF/go-micro-kit/docs"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"
	"go-micro.dev/v4/util/log"
)

var console bool

func main() {
	service := initService()
	// Run the server
	if console {
		fmt.Printf(docs.Console)
	}
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

// initService init a micro.service
func initService() micro.Service {
	service := micro.NewService(
		// Add runtime flags
		// We could do this below too
		micro.Flags(
			&cli.BoolFlag{
				Name:    "console",
				Aliases: []string{"c"},
				Usage:   "if console is true, will console start screen when start",
			},
		),
		micro.Flags(
			&cli.StringFlag{
				Name:  "float_flag",
				Usage: "this is a float flag",
			}),
	)
	service.Init(
		// Add runtime action
		// We could actually do this above
		micro.Action(func(c *cli.Context) error {
			console = c.Bool("console")
			fmt.Println(fmt.Sprintf("start params is console : %t", console))
			return nil
		}),
	)
	return service
}
