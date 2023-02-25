package producer

import (
	"github.com/nsqio/go-nsq"
	"go-micro.dev/v4/util/log"
)

type NSQProducer struct {
	Addr        string `json:"addr"`
	MessageBody []byte `json:"messageBody"`
	TopicName   string `json:"topicName"`
}

func NewNSQProducer() *NSQProducer {
	return &NSQProducer{}
}

func (n *NSQProducer) Producer() {
	// Instantiate a producer.
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(n.Addr, config)
	defer func() {
		// Gracefully stop the producer when appropriate (e.g. before shutting down the service)
		producer.Stop()
	}()
	if err != nil {
		log.Fatal(err)
	}
	// Synchronously publish a single message to the specified topic.
	// Messages can also be sent asynchronously and/or in batches.
	err = producer.Publish(n.TopicName, n.MessageBody)
	if err != nil {
		log.Fatal(err)
	}
}
