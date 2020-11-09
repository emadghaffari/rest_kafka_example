package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/emadghaffari/rest_kafka_example/model/elasticsearch"
	"github.com/spf13/viper"
)

// Consumer struct
type Consumer struct{}

// Sarama configuration options
var (
	brokers  = []string{}
	version  = ""
	group    = ""
	topics   = []string{}
	assignor = ""
	oldest   = false
	verbose  = false
	Topic    = ""
	config   *sarama.Config
)

// Init func
func Init() {
	brokers = viper.GetStringSlice("kafka.brokers")
	version = viper.GetString("kafka.version")
	group = viper.GetString("kafka.group")
	topics = viper.GetStringSlice("kafka.topics")
	assignor = viper.GetString("kafka.assignor")
	oldest = viper.GetBool("kafka.oldest")
	verbose = viper.GetBool("kafka.verbose")
	Topic = viper.GetString("kafka.Topic")

	if len(brokers) == 0 {
		panic("no Kafka bootstrap brokers defined, please set the -brokers flag")
	}

	if len(topics) == 0 {
		panic("no topics given to be consumed, please set the -topics flag")
	}

	if len(group) == 0 {
		panic("no Kafka consumer group defined, please set the -group flag")
	}
	config = sarama.NewConfig()
	config.ClientID = "go-kafka"
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Consumer.MaxProcessingTime = time.Second
	config.Producer.Idempotent = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Net.MaxOpenRequests = 1
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Return.Successes = true
	config.Producer.Retry.Backoff = time.Duration(time.Second * 5)
	config.Producer.Retry.Max = 5
	config.Producer.Compression = sarama.CompressionLZ4
	config.Producer.Timeout = time.Duration(time.Second * 50)
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange

	// Auth
	config.Net.SASL.Enable = true
	config.Net.SASL.Handshake = true
	config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
	config.Net.SASL.User = "admin"
	config.Net.SASL.Password = "admin-secret"
	config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }

	// switch assignor {
	// case "sticky":
	// 	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	// case "roundrobin":
	// 	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	// case "range":
	// 	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	// default:
	// 	log.Panicf("Unrecognized consumer group partition assignor: %s", assignor)
	// }

}

// SyncProducer func
func SyncProducer() (sarama.SyncProducer, error) {
	syncProducer, proErr := sarama.NewSyncProducer(brokers, config)
	if proErr != nil {
		panic(proErr)
	}
	return syncProducer, nil
}

// NewConsumer func
func newConsumer() (sarama.ConsumerGroup, error) {
	group, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		panic(err)
	}
	go func() {
		for err := range group.Errors() {
			panic(err)
		}
	}()

	return group, nil
}

// Consumer func
func (c *Consumer) Consumer() {
	group, _ := newConsumer()
	defer func() {
		if err := group.Close(); err != nil {
			panic(err)
		}
	}()
	func() {
		ctx := context.Background()
		for {
			err := group.Consume(ctx, topics, c)
			if err != nil {
				fmt.Printf("kafka consume failed: %v, sleeping and retry in a moment\n", err)
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()
}

// Producer func
func Producer(items []string) {
	syncProducer, _ := SyncProducer()
	for _, item := range items {
		_, _, err := syncProducer.SendMessage(&sarama.ProducerMessage{
			Topic: Topic,
			Value: sarama.StringEncoder(item),
		})
		time.Sleep(time.Millisecond * 10)

		if err != nil {
			log.Fatalln("failed to send message to ", Topic, err)
		}
	}

	_ = syncProducer.Close()
}

// Setup meth
func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error {
	fmt.Println("*********************Setup*******************")
	return nil
}

// Cleanup meth
func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	fmt.Println("*********************Cleanup*******************")
	return nil
}

// ConsumeClaim meth
func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	fmt.Println("*********************ConsumeClaim*******************")

	for msg := range claim.Messages() {
		id := fmt.Sprintf("%v-%d-%d", msg.Topic, msg.Partition, msg.Offset)
		go elasticsearch.Save(id, string(msg.Value))
	}
	return nil
}
