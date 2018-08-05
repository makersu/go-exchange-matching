package main_test

import (
	"fmt"
	"os/signal"
	"testing"
	"time"

	"os"

	. "git.cc/Core2.0/bitopro-matching"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	jsoniter "github.com/json-iterator/go"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
)

var (
	brokers = []string{"127.0.0.1:9092"}
	topic   = "inputEvent"
	// topics  = []string{topic}
	partition = 0
	offset    = sarama.OffsetOldest //sarama.OffsetNewest
)

var maxSamaraCount = 100000

func NewTestConfig() *sarama.Config {

	config := sarama.NewConfig()
	config.ChannelBufferSize = 1024
	// config.Version = sarama.V0_10_1_0
	config.Version = sarama.V0_11_0_2
	return config
}

func noTestSarama(t *testing.T) {
	fmt.Println("TestSarama")

	consumer := NewSaramaConsumer(brokers, NewTestConfig())
	defer consumer.Close()

	partitionConsumer := NewPartitionConsumer(consumer, topic, int32(partition), offset)
	defer partitionConsumer.Close()

	kafkaSelectOrderEvents(partitionConsumer)
}

func kafkaSelectOrderEvents(partitionConsumer sarama.PartitionConsumer) {
	fmt.Println("kafkaSelectOrderEvents")

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// event := new(OrderEvent)
	// var event OrderEvent

	msgCount := 0
	start := time.Now()

ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			// log.Info(msg)
			msgCount++
			if msgCount == 1 {
				start = time.Now()
			}

			printMsgCounter(start, msgCount, msg)

			if msgCount == maxSamaraCount {
				break ConsumerLoop
			}

		case <-signals:
			break ConsumerLoop

			// default:
			// 	if msgCount%100000 == 0 {
			// 		fmt.Println("no message received")
			// 	}
		}
	}
	elapsed := time.Now().Sub(start)
	fmt.Printf("sarama: %v records, %v\n", msgCount, elapsed)
}

func TestSaramaUnmarshal(t *testing.T) {
	fmt.Println("TestSaramaUnmarshal")

	consumer := NewSaramaConsumer(brokers, NewTestConfig())
	defer consumer.Close()

	partitionConsumer := NewPartitionConsumer(consumer, topic, int32(partition), offset)
	defer partitionConsumer.Close()

	kafkaSelectUnmarshalOrderEvents(partitionConsumer)
}

func kafkaSelectUnmarshalOrderEvents(partitionConsumer sarama.PartitionConsumer) {
	fmt.Println("kafkaSelectUnmarshalOrderEvents")

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// event := new(OrderEvent)
	// var event OrderEvent

	msgCount := 0
	start := time.Now()

ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			msgCount++

			if msgCount == 1 {
				start = time.Now()
			}

			var event OrderEvent

			if err := jsoniter.Unmarshal(msg.Value, &event); err == nil {
				// event.Order.Amount = 1
				// engine.Run(event.Order)

				printMsgCounter(start, msgCount, msg)

			} else {
				log.Info("jsoniter.Unmarshal error:", err)
				break ConsumerLoop
			}

			if msgCount == maxSamaraCount {
				break ConsumerLoop
			}

		case <-signals:
			break ConsumerLoop

			// default:
			// 	fmt.Println("no message received")

		}
	}
}

func TestSaramaUnmarshalAddOrder(t *testing.T) {
	fmt.Println("TestSaramaUnmarshalAddOrder")

	consumer := NewSaramaConsumer(brokers, NewTestConfig())
	defer consumer.Close()

	partitionConsumer := NewPartitionConsumer(consumer, topic, int32(partition), offset)
	defer partitionConsumer.Close()

	kafkaSelectAddOrderEvents(partitionConsumer)
}

func kafkaSelectAddOrderEvents(partitionConsumer sarama.PartitionConsumer) {
	fmt.Println("kafkaSelectAddOrderEvents")

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	engine := NewEngine("BTC_ETH")

	// event := new(OrderEvent)
	// var event OrderEvent

	msgCount := 0
	start := time.Now()

ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			msgCount++

			if msgCount == 1 {
				start = time.Now()
			}

			var event OrderEvent

			if err := jsoniter.Unmarshal(msg.Value, &event); err == nil {
				// event.Order.Amount = 1
				engine.AddAskOrderBook(event.Order)

				printMsgCounter(start, msgCount, msg)

			} else {
				log.Info("jsoniter.Unmarshal error:", err)
				break ConsumerLoop
			}

			if msgCount == maxSamaraCount {
				break ConsumerLoop
			}

		case <-signals:
			break ConsumerLoop

			// default:
			// 	fmt.Println("no message received")

		}
	}

}

func TestSaramaUnmarshalRunOrder(t *testing.T) {
	fmt.Println("TestSaramaUnmarshalRunOrder")

	consumer := NewSaramaConsumer(brokers, NewTestConfig())
	defer consumer.Close()

	partitionConsumer := NewPartitionConsumer(consumer, topic, int32(partition), offset)
	defer partitionConsumer.Close()

	kafkaSelectRunOrderEvents(partitionConsumer)
}

// sarama consumer example and Go by Example: Non-Blocking Channel Operations
func kafkaSelectRunOrderEvents(partitionConsumer sarama.PartitionConsumer) {
	fmt.Println("kafkaSelectRunOrderEvents")

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	engine := NewRBTEngine("BTC_ETH")

	// event := new(OrderEvent)
	// var event OrderEvent

	msgCount := 0
	start := time.Now()

ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			msgCount++

			if msgCount == 1 {
				start = time.Now()
			}

			var event OrderEvent

			if err := jsoniter.Unmarshal(msg.Value, &event); err == nil {
				// event.Order.Amount = 1 //
				// log.Debug("event.Order:", event.Order)
				// event.Order.Price += 0.123456789 //

				engine.Run(event.Order)

				printMsgCounter(start, msgCount, msg)

			} else {
				log.Info("jsoniter.Unmarshal error:", err)
				break ConsumerLoop
			}

			if msgCount == maxSamaraCount {
				break ConsumerLoop
			}

		case <-signals:
			break ConsumerLoop

			// default:
			// 	fmt.Println("no message received")

		}
	}
}

// func noTestRangeConsumeOrderEvents(t *testing.T) {
// 	fmt.Println("TestRangeConsumeOrderEvents")

// 	consumer := NewSaramaConsumer(brokers, NewTestConfig())
// 	defer consumer.Close()

// 	partitionConsumer := NewPartitionConsumer(consumer, topic, int32(partition), offset)
// 	defer partitionConsumer.Close()

// 	rangeConsumeOrderEvents(partitionConsumer)
// }

// // kafka example around25
// func rangeConsumeOrderEvents(partitionConsumer sarama.PartitionConsumer) {
// 	fmt.Println("rangeConsumeOrderEvent")

// 	done := make(chan bool)

// 	// engine := NewEngine("BTC_ETH")

// 	// event := new(OrderEvent)
// 	// var event OrderEvent

// 	msgCount := 0
// 	start := time.Now()

// 	go func() {
// 		for msg := range partitionConsumer.Messages() {
// 			var event OrderEvent

// 			if err := jsoniter.Unmarshal(msg.Value, &event); err == nil {
// 				// event.Order.Amount = 1
// 				// engine.Run(event.Order)

// 				msgCount++
// 				printMsgCounter(start, msgCount, msg)

// 			} else {
// 				log.Info("jsoniter.Unmarshal error:", err)
// 				// done <- true
// 			}

// 		}
// 		done <- true
// 	}()

// 	<-done
// }

func printMsgCounter(start time.Time, msgCount int, msg *sarama.ConsumerMessage) {
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

func noTestSaramaClusterUnmarshal(t *testing.T) {
	fmt.Println("TestSaramaClusterUnmarshal")

	config := cluster.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.ChannelBufferSize = 1024
	config.Version = sarama.V0_11_0_2

	groupID := ksuid.New().String()

	// consumer, err := cluster.NewConsumer(brokers, "banku-consumer", topics, config)
	consumer, err := cluster.NewConsumer(brokers, groupID, []string{topic}, config)

	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	// consume notifications
	go func() {
		for n := range consumer.Notifications() {
			log.Printf("Notification: %v\n", n)
		}
	}()

	var start time.Time
	var count int

	// consume messages, watch signals
loop:
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if !ok {
				panic("messages channel unexpectedly closed")
			}

			count++

			if count == 1 {
				start = time.Now()
			}

			var event OrderEvent

			if err := jsoniter.Unmarshal(msg.Value, &event); err == nil {
				// event.Order.Amount = 1
				// engine.Run(event.Order)

				printMsgCounter(start, count, msg)

			} else {
				log.Info("jsoniter.Unmarshal error:", err)
				break loop
			}

			if count == maxSamaraCount {
				break loop
			}

		case <-signals:
			return
		}
	}
	elapsed := time.Now().Sub(start)
	fmt.Printf("sarama-cluster: %v records, %v\n", count, elapsed)

}
