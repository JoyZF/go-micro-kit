package handler

import (
	"github.com/JoyZF/go-micro-kit/common/kafka/producer"
	"github.com/Shopify/sarama"
	"go-micro.dev/v4/util/log"
)

func greeter(producerProvider *producer.ProducerProvider, recordsNumber int64, topic string, value string) {
	producer := producerProvider.Borrow()
	defer producerProvider.Release(producer)

	// Start kafka transaction
	err := producer.BeginTxn()
	if err != nil {
		log.Infof("unable to start txn %s\n", err)
		return
	}

	// Produce some records in transaction
	var i int64
	for i = 0; i < recordsNumber; i++ {
		producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.StringEncoder(value)}
	}

	// commit transaction
	err = producer.CommitTxn()
	if err != nil {
		log.Infof("Producer: unable to commit txn %s\n", err)
		for {
			if producer.TxnStatus()&sarama.ProducerTxnFlagFatalError != 0 {
				// fatal error. need to recreate producer.
				log.Infof("Producer: producer is in a fatal state, need to recreate it")
				break
			}
			// If producer is in abortable state, try to abort current transaction.
			if producer.TxnStatus()&sarama.ProducerTxnFlagAbortableError != 0 {
				err = producer.AbortTxn()
				if err != nil {
					// If an error occured just retry it.
					log.Infof("Producer: unable to abort transaction: %+v", err)
					continue
				}
				break
			}
			// if not you can retry
			err = producer.CommitTxn()
			if err != nil {
				log.Infof("Producer: unable to commit txn %s\n", err)
				continue
			}
		}
		return
	}
}
