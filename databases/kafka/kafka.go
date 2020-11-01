package kafka

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/emadghaffari/rest_kafka_example/model/elasticsearch"
)

// Consumer struct
type Consumer struct {}

// Sarama configuration options
var (
	brokers  = []string{"kafka1:9092","kafka2:9092"}
	version  = ""
	group    = "go-kafka"
	topics   = []string{"first_topic"}
	assignor = "roundrobin"
	oldest   = true
	verbose  = false
	Topic     = flag.String("topic", "first_topic", "The Kafka topic to use")
	config *sarama.Config
)

func init() {
	// flag.StringVar(&group, "group", "go-kafka", "Kafka consumer group definition")
	// flag.StringVar(&version, "version", "2.1.1", "Kafka cluster version")
	// flag.StringVar(&topics, "topics", "first_topic", "Kafka topics to be consumed, as a comma separated list")
	// flag.StringVar(&assignor, "assignor", "roundrobin", "Consumer group partition assignment strategy (range, roundrobin, sticky)")
	// flag.BoolVar(&oldest, "oldest", true, "Kafka consumer consume initial offset from oldest")
	// flag.BoolVar(&verbose, "verbose", false, "Sarama logging")
	// flag.Parse()

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
	// config.Producer.Idempotent = true
	// config.Producer.Retry.Max = int(^uint(0)  >> 1)  
	// config.Producer.Flush.MaxMessages = 5
	// config.Producer.Flush.Bytes = int(32*1024)
	config.Producer.Compression = sarama.CompressionSnappy
	switch assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", assignor)
	}
	
}

// SyncProducer func
func SyncProducer() (sarama.SyncProducer, error) {
	syncProducer,proErr := sarama.NewSyncProducer(brokers,config)
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

	return group,nil
}

// Consumer func
func (c *Consumer) Consumer()  {
	group,_ := newConsumer()

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
			  time.Sleep(time.Second)
		   }
		}
	 }()
}
// Producer func
func Producer(items []string)  {
	syncProducer,_ := SyncProducer()
	for _, item := range items {
		_, _, err := syncProducer.SendMessage(&sarama.ProducerMessage{
			Topic: *Topic,
			Value: sarama.StringEncoder(item),
		})
		time.Sleep(time.Millisecond * 10)
	
		if err != nil {
			log.Fatalln("failed to send message to ", *Topic, err)
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
		id := fmt.Sprintf("%v-%d-%d",msg.Topic,msg.Partition,msg.Offset)
		go elasticsearch.Save(id, string(msg.Value))
	}
	return nil
 }