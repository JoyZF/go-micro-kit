package consumer

import (
	"github.com/nsqio/go-nsq"
	"go-micro.dev/v4/util/log"
	"os"
	"os/signal"
	"syscall"
)

type NSQConsumer struct {
	Topic   string `json:"topic"`
	Channel string `json:"channel"`
	Addr    string `json:"addr"`
}

func NewNSQConsumer(topic, channel, addr string) *NSQConsumer {
	return &NSQConsumer{
		Topic:   topic,
		Channel: channel,
		Addr:    addr,
	}
}

func (n *NSQConsumer) Consumer(h nsq.Handler) {
	go func() {
		config := nsq.NewConfig()
		consumer, err := nsq.NewConsumer(n.Topic, n.Channel, config)
		if err != nil {
			log.Fatal(err)
		}
		consumer.AddHandler(h)
		err = consumer.ConnectToNSQLookupd(n.Addr)
		if err != nil {
			log.Fatal(err)
		}
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		// Gracefully stop the consumer.
		consumer.Stop()
	}()
}
