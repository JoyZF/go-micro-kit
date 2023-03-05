package producer

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"go-micro.dev/v4/util/log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type KafkaProducer struct {
	Version     string
	Brokers     string
	Producers   int
	KeepRunning bool
}

// pool of producers that ensure transactional-id is unique.
type ProducerProvider struct {
	transactionIdGenerator int32

	producersLock sync.Mutex
	producers     []sarama.AsyncProducer

	producerProvider func() sarama.AsyncProducer
}

type fn func(producerProvider *ProducerProvider, recordsNumber int64, topic string, value string)

func NewKafkaProducer() *KafkaProducer {
	return &KafkaProducer{}
}

func (k *KafkaProducer) Producer(topic string, message string, num int64, f fn) {
	version, err := sarama.ParseKafkaVersion(k.Version)
	if err != nil {
		log.Errorf("sarama ParseKafkaVersion err %+v", err)
	}

	producerProvider := newProducerProvider(strings.Split(k.Brokers, ","), func() *sarama.Config {
		config := sarama.NewConfig()
		config.Version = version
		config.Producer.Idempotent = true
		config.Producer.Return.Errors = false
		config.Producer.RequiredAcks = sarama.WaitForAll
		config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
		config.Producer.Transaction.Retry.Backoff = 10
		config.Producer.Transaction.ID = "txn_producer"
		config.Net.MaxOpenRequests = 1
		return config
	})

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	for i := 0; i < k.Producers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					f(producerProvider, num, topic, message)
				}
			}
		}()
	}
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for k.KeepRunning {
		<-sigterm
		log.Infof("terminating: via signal")
		k.KeepRunning = false
	}
	cancel()
	wg.Wait()

	producerProvider.clear()
}

func newProducerProvider(brokers []string, producerConfigurationProvider func() *sarama.Config) *ProducerProvider {
	provider := &ProducerProvider{}
	provider.producerProvider = func() sarama.AsyncProducer {
		config := producerConfigurationProvider()
		suffix := provider.transactionIdGenerator
		// Append transactionIdGenerator to current config.Producer.Transaction.ID to ensure transaction-id uniqueness.
		if config.Producer.Transaction.ID != "" {
			provider.transactionIdGenerator++
			config.Producer.Transaction.ID = config.Producer.Transaction.ID + "-" + fmt.Sprint(suffix)
		}
		producer, err := sarama.NewAsyncProducer(brokers, config)
		if err != nil {
			return nil
		}
		return producer
	}
	return provider
}

func (p *ProducerProvider) Borrow() (producer sarama.AsyncProducer) {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	if len(p.producers) == 0 {
		for {
			producer = p.producerProvider()
			if producer != nil {
				return
			}
		}
	}

	index := len(p.producers) - 1
	producer = p.producers[index]
	p.producers = p.producers[:index]
	return
}

func (p *ProducerProvider) Release(producer sarama.AsyncProducer) {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	// If released producer is erroneous close it and don't return it to the producer pool.
	if producer.TxnStatus()&sarama.ProducerTxnFlagInError != 0 {
		// Try to close it
		_ = producer.Close()
		return
	}
	p.producers = append(p.producers, producer)
}

func (p *ProducerProvider) clear() {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	for _, producer := range p.producers {
		producer.Close()
	}
	p.producers = p.producers[:0]
}
