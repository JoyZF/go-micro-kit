package handler

import (
	"fmt"
	"github.com/nsqio/go-nsq"
)

type Greeter struct {
}

func (g *Greeter) HandleMessage(message *nsq.Message) error {
	fmt.Println(string(message.Body))
	return nil
}
