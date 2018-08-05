package main_test

import (
	"fmt"
	"testing"
	"time"

	. "git.cc/Core2.0/bitopro-matching"
	// "github.com/Shopify/sarama"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	jsoniter "github.com/json-iterator/go"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
)

var (
	// brokers = []string{"127.0.0.1:9092"}
	confluentBrokers      = "localhost:9092"
	confluentBrokersTopic = "inputEvent"
	// topics  = []string{topic}
	// partition = 0
	// offset    = sarama.OffsetOldest //sarama.OffsetNewest
)

var maxCount = 100000

func printConfluentCounter(start time.Time, msgCount int, msg *kafka.Message) {
	if msgCount%10000 == 0 {
		// log.Info("msgCount:", msgCount)
		// log.Info("msg.offset:", msg.Offset)
		// log.Info("msg.Value:", string(msg.Value))
		// log.Info("event:    ", event)

		elapsed := time.Since(start)
		tps := int(float64(msgCount) / elapsed.Seconds())

		log.Info("Handled ", msgCount, " events in ", elapsed, " at ", tps, " events/second.")
	}

}

func noTestConfluent(t *testing.T) {
	fmt.Println("TestConfluent")

	cm := NewConfluentConfigMap()

	consumer, err := kafka.NewConsumer(cm)
	check(err)
	defer consumer.Close()

	check(consumer.Subscribe(confluentBrokersTopic, nil))

	var start time.Time
	count := 0

loop:
	for {
		select {
		case m, ok := <-consumer.Events():
			if !ok {
				panic("unexpected eof")
			}

			switch event := m.(type) {
			case kafka.AssignedPartitions:
				consumer.Assign(event.Partitions)

			case kafka.PartitionEOF:
				// nop

			case kafka.RevokedPartitions:
				consumer.Unassign()

			case *kafka.Message:
				count++
				if count == 1 {
					start = time.Now()
				}
				printConfluentCounter(start, count, event)

				if count == maxCount {
					break loop
				}

			default:
				panic(m)
			}
		}
	}
	elapsed := time.Now().Sub(start)
	fmt.Printf("confluent: %v records, %v\n", count, elapsed)
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func NewConfluentConfigMap() *kafka.ConfigMap {
	groupID := ksuid.New().String()
	cm := &kafka.ConfigMap{
		"session.timeout.ms":              6000,
		"metadata.broker.list":            confluentBrokers,
		"enable.auto.commit":              false,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"group.id":                        groupID,
		"default.topic.config": kafka.ConfigMap{
			"auto.offset.reset": "earliest",
		},

		// "security.protocol": "ssl",
		// "ssl.ca.location":          caFile,
		// "ssl.certificate.location": certFile,
		// "ssl.key.location":         keyFile,
	}

	return cm
}
func TestConfluentUnmarshal(t *testing.T) {
	fmt.Println("TestConfluentUnmarshal")

	cm := NewConfluentConfigMap()
	consumer, err := kafka.NewConsumer(cm)
	check(err)
	defer consumer.Close()

	check(consumer.Subscribe(confluentBrokersTopic, nil))

	var start time.Time
	count := 0

loop:
	for {
		select {
		case m, ok := <-consumer.Events():
			if !ok {
				panic("unexpected eof")
			}

			switch event := m.(type) {
			case kafka.AssignedPartitions:
				consumer.Assign(event.Partitions)

			case kafka.PartitionEOF:
				// nop

			case kafka.RevokedPartitions:
				consumer.Unassign()

			case *kafka.Message:

				count++
				if count == 1 {
					start = time.Now()
				}

				var orderEvent OrderEvent

				if err := jsoniter.Unmarshal(event.Value, &orderEvent); err == nil {

					printConfluentCounter(start, count, event)

				} else {
					log.Info("jsoniter.Unmarshal error:", err)
					break loop
				}

				if count == maxCount {
					break loop
				}

			default:
				panic(m)
			}
		}
	}
	elapsed := time.Now().Sub(start)
	fmt.Printf("confluent: %v records, %v\n", count, elapsed)
}

func TestConfluentUnmarshalAddOrder(t *testing.T) {
	fmt.Println("TestConfluentUnmarshalAddOrder")

	cm := NewConfluentConfigMap()
	consumer, err := kafka.NewConsumer(cm)
	check(err)
	defer consumer.Close()

	check(consumer.Subscribe(confluentBrokersTopic, nil))

	var start time.Time
	count := 0

	engine := NewEngine("BTC_ETH")

loop:
	for {
		select {
		case m, ok := <-consumer.Events():
			if !ok {
				panic("unexpected eof")
			}

			switch event := m.(type) {
			case kafka.AssignedPartitions:
				consumer.Assign(event.Partitions)

			case kafka.PartitionEOF:
				// nop

			case kafka.RevokedPartitions:
				consumer.Unassign()

			case *kafka.Message:

				count++
				if count == 1 {
					start = time.Now()
				}

				var orderEvent OrderEvent

				if err := jsoniter.Unmarshal(event.Value, &orderEvent); err == nil {

					engine.AddAskOrderBook(orderEvent.Order)

					printConfluentCounter(start, count, event)

				} else {
					log.Info("jsoniter.Unmarshal error:", err)
					break loop
				}

				if count == maxCount {
					break loop
				}

			default:
				panic(m)
			}
		}
	}
	elapsed := time.Now().Sub(start)
	fmt.Printf("confluent: %v records, %v\n", count, elapsed)
}

func TestConfluentUnmarshalRunOrder(t *testing.T) {
	fmt.Println("TestConfluentUnmarshalRunOrder")

	cm := NewConfluentConfigMap()
	consumer, err := kafka.NewConsumer(cm)
	check(err)
	defer consumer.Close()

	check(consumer.Subscribe(confluentBrokersTopic, nil))

	var start time.Time
	count := 0

	engine := NewEngine("BTC_ETH")

loop:
	for {
		select {
		case m, ok := <-consumer.Events():
			if !ok {
				panic("unexpected eof")
			}

			switch event := m.(type) {
			case kafka.AssignedPartitions:
				consumer.Assign(event.Partitions)

			case kafka.PartitionEOF:
				// nop

			case kafka.RevokedPartitions:
				consumer.Unassign()

			case *kafka.Message:

				count++
				if count == 1 {
					start = time.Now()
				}

				var orderEvent OrderEvent

				if err := jsoniter.Unmarshal(event.Value, &orderEvent); err == nil {
					orderEvent.Order.Amount = 1
					engine.Run(orderEvent.Order)
					printConfluentCounter(start, count, event)

				} else {
					log.Info("jsoniter.Unmarshal error:", err)
					break loop
				}

				if count == maxCount {
					break loop
				}

			default:
				panic(m)
			}
		}
	}
	elapsed := time.Now().Sub(start)
	fmt.Printf("confluent: %v records, %v\n", count, elapsed)
}
