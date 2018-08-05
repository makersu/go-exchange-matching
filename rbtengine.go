package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/Shopify/sarama"
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/emirpasic/gods/utils"
)

// TODO: refactoring to move and log format
func init() {

	// log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006/01/02 - 15:04:05",
		FullTimestamp:   true,
	})

	log.SetLevel(log.InfoLevel)
	// log.SetLevel(log.DebugLevel)

}

// TODO: convention
// TODO: array for orderbook
type RBTEngine struct {
	pair    string
	askbook *rbt.Tree
	bidbook *rbt.Tree
}

// NewEngine is the Matching Engine constructor
func NewRBTEngine(pair string) RBTEngine {
	return RBTEngine{
		pair: pair,
		// askbook: rbt.NewWithIntComparator(),
		askbook: rbt.NewWith(utils.Float64Comparator),
		// bidbook: rbt.NewWithIntComparator(),
		bidbook: rbt.NewWith(utils.Float64Comparator),
		// writeLog: NewWriteLog(exchange.String()),
	}
}

// test
func (engine *RBTEngine) AddAskOrderBook(order *Order) {
	treeNode, ok := engine.askbook.Get(order.Price)
	if !ok {
		node := NewTreeNode()
		node.upsert(order)
		engine.askbook.Put(order.Price, node)

	} else {
		node := treeNode.(*TreeNode)
		node.upsert(order)
		engine.askbook.Put(order.Price, treeNode)
	}
}

// refactoring: rename for private
//addOrder adds an order that cannot be filled any further to the orderbook
// func (engine *Engine) addOrder(order *Order) *Order {
func (engine *RBTEngine) addOrder(book *rbt.Tree, order *Order) {

	treeNode, ok := book.Get(order.Price)
	if !ok {
		node := NewTreeNode()
		node.upsert(order)
		book.Put(order.Price, node)

	} else {
		node := treeNode.(*TreeNode)
		node.upsert(order)
		book.Put(order.Price, treeNode)
	}

}

func (engine *RBTEngine) addAskOrder(sellbook *rbt.Tree, order *Order) {
	addOrder(sellbook, order)
}

func (engine *RBTEngine) addBidOrder(buybook *rbt.Tree, order *Order) {
	addOrder(buybook, order)
}

func (engine *RBTEngine) Run(order *Order) {
	// log.Info("Run")

	switch order.OrderAction {
	case "SELL":
		log.Debug("*Run ASK for Highest price ", order.Price, order)

		// Ask Operations
		matchingOrder, matches := engine.matchAskOrder(order)

		log.Debug("*Run ASK matches ", len(matches), matches)

		if matchingOrder.Amount > 0 {
			addAskOrder(engine.askbook, order)
		}

	case "BUY":
		log.Debug("*Run BID for Lowest price ", order.Price, order)

		// Bid Operations
		matchingOrder, matches := engine.matchBidOrder(order)

		log.Debug("*Run BID matches ", len(matches), matches)

		if matchingOrder.Amount > 0 {
			addBidOrder(engine.bidbook, order)
		}

	// case CANCEL:
	// matchingOrder, matches := engine.cancelOrder(order)

	default:
		//Drop the message
		log.Info("Invalid Order Type")
	}

	// engine.printEngine()

}

func (engine *RBTEngine) matchBidOrder(askOrder *Order) (*Order, []Match) {
	log.Debug("matchBidOrder")
	return engine.matchAskBook(askOrder)
}

func (engine *RBTEngine) matchAskBook(bidOrder *Order) (*Order, []Match) {

	log.Debug("matchAskBook askbook.Keys() ", engine.askbook.Keys())

	var matches []Match

	askbookit := engine.askbook.Iterator()

	// get the begin node(lowest price) then next
	for askbookit.Begin(); askbookit.Next(); {
		askBookNodePrice, askBookNode := askbookit.Key().(float64), askbookit.Value().(*TreeNode)
		log.Debug("matchAskBook bidOrder    ", bidOrder.Price, bidOrder)
		log.Debug("matchAskBook aksBookNode ", askBookNodePrice, askBookNode)

		//Check price
		if bidOrder.Price >= askBookNodePrice {
			log.Debug("matchAskBook bidOrder.Price >= askBookNodePrice ", bidOrder.Price, askBookNodePrice)

			_, nodeMatches := engine.matchBookNode(engine.askbook, askBookNode, bidOrder)

			log.Debug("matchAskBook nodeMatches ", len(nodeMatches), nodeMatches)

			for _, match := range nodeMatches {
				matches = append(matches, match)
			}

		} else {
			log.Debug("matchAskBook bidOrder.Price < askBookNodePrice ", bidOrder.Price, askBookNodePrice)
			// continue?
			break

		}

		if bidOrder.Amount == 0 {
			break
		}
	}

	return bidOrder, matches
}

func (engine *RBTEngine) matchAskOrder(askOrder *Order) (*Order, []Match) {
	log.Debug("matchAskOrder")
	return engine.matchBidBook(askOrder)
}

func (engine *RBTEngine) matchBidBook(askOrder *Order) (*Order, []Match) {
	log.Debug("matchBidBook bidbook.Keys() ", engine.bidbook.Keys())

	var matches []Match

	bidbookit := engine.bidbook.Iterator()

	// get the begin node(lowest price) then next
	for bidbookit.End(); bidbookit.Prev(); {

		bidBookNodePrice, bidBookNode := bidbookit.Key().(float64), bidbookit.Value().(*TreeNode)
		log.Debug("matchBidBook askOrder    ", askOrder.Price, askOrder)
		log.Debug("matchBidBook bidBookNode ", bidBookNodePrice, bidBookNode)

		//Check price
		if bidBookNodePrice >= askOrder.Price {
			log.Debug("bidBookNodePrice >= matchBidBook askOrder.Price ", bidBookNodePrice, askOrder.Price)

			_, nodeMatches := engine.matchBookNode(engine.bidbook, bidBookNode, askOrder)

			log.Debug("matchBidBook nodeMatches ", len(nodeMatches), nodeMatches)

			for _, match := range nodeMatches {
				matches = append(matches, match)
			}

		} else {
			log.Debug("matchBidBook askOrder.Price > bidBookNodePrice ", askOrder.Price, bidBookNodePrice)
			// continue?
			break

		}

		if askOrder.Amount == 0 {
			break
		}
	}

	return askOrder, matches
}

func (engine *RBTEngine) matchBookNode(book *rbt.Tree, node *TreeNode, matchingOrder *Order) (*Order, []Match) {
	log.Debug("matchBookNode")

	var nodeMatches []Match

	for _, nodeOrder := range node.orders {

		if nodeOrder.Amount >= matchingOrder.Amount {
			log.Debug("matchBookNode nodeOrder.Amount > matchingOrder.Amount ", nodeOrder.Amount, matchingOrder.Amount)

			nodeMatch := NewMatch(*matchingOrder, *nodeOrder, matchingOrder.Amount, nodeOrder.Price)
			log.Debug("matchBookNode nodeMatch ", nodeMatch)

			nodeMatches = append(nodeMatches, nodeMatch)

			nodeOrder.Amount = nodeOrder.Amount - matchingOrder.Amount

			if nodeOrder.Amount == 0 {
				engine.removeNodeOrder(book, node, nodeOrder)
			}

			matchingOrder.Amount = matchingOrder.Amount - matchingOrder.Amount //0
			break                                                              // break if matchingOrder.Amount == 0

		} else if nodeOrder.Amount < matchingOrder.Amount {
			log.Debug("matchBookNode nodeOrder.Amount < matchingOrder.Amount ", nodeOrder.Amount, matchingOrder.Amount)

			nodeMatch := NewMatch(*matchingOrder, *nodeOrder, nodeOrder.Amount, nodeOrder.Price)
			log.Debug("matchBookNode nodeMatch ", nodeMatch)

			nodeMatches = append(nodeMatches, nodeMatch)

			matchingOrder.Amount = matchingOrder.Amount - nodeOrder.Amount

			nodeOrder.Amount = nodeOrder.Amount - nodeOrder.Amount //0
			engine.removeNodeOrder(book, node, nodeOrder)          //remove order if nodeOrder.Amount == 0

		}
		// else if nodeOrder.Amount == matchingOrder.Amount {
		// 	log.Debug("matchBookNode nodeOrder.NumberOutstanding == matchingOrder.NumberOutstanding ", nodeOrder.Amount, matchingOrder.Amount)

		// 	nodeMatch := NewMatch(*matchingOrder, *nodeOrder, matchingOrder.Amount, nodeOrder.Price)
		// 	nodeMatches = append(nodeMatches, nodeMatch)
		// 	log.Debug("matchBookNode nodeMatch ", nodeMatch)

		// 	nodeOrder.Amount = nodeOrder.Amount - matchingOrder.Amount //0
		// 	matchingOrder.Amount = matchingOrder.Amount - matchingOrder.Amount //0
		// 	engine.removeNodeOrder(book, node, nodeOrder)
		// }

		// if matchingOrder.Amount == 0 {
		// 	log.Debug("matchBookNode matchingOrder.NumberOutstanding == 0 break matchBookNode")
		// 	break
		// }

	}

	return matchingOrder, nodeMatches
}

//TODO: linked list?
//TODO: sort after add?
// func (n TreeNode) sortedOrders() []*Order {
// 	log.Debug("sortedOrders")
// 	orders := make([]*Order, 0)
// 	for _, v := range n.orders {
// 		orders = append(orders, v)
// 	}

// 	sort.Slice(orders[:], func(i, j int) bool {
// 		return orders[i].Timestamp < orders[j].Timestamp
// 	})

// 	return orders
// }

// func (n *TreeNode) delete(id uint64) {
func (engine *RBTEngine) removeNodeOrder(book *rbt.Tree, node *TreeNode, order *Order) {
	delete(node.orders, order.OrderId)
	if len(node.orders) == 0 {
		book.Remove(order.Price)
	}

}

// func (n *TreeNode) delete(id uint64) {
func (engine *RBTEngine) removeOrder(book *rbt.Tree, n *TreeNode, order *Order) {
	delete(n.orders, order.OrderId)
	if len(n.orders) == 0 {
		book.Remove(order.Price)
	}

}

//TODO: refactor to rename
func (engine *RBTEngine) printEngine() {
	log.Info("printOrderbook ASKbook")
	printOrderbook(engine.askbook) //
	log.Info("printOrderbook BIDbook")
	printOrderbook(engine.bidbook) //
}

// func (engine *Engine) consume(in <-chan *Order) {
// 	for order := range in {
// 		engine.Run(order)
// 	}
// }

func (engine *RBTEngine) rangeConsume(partitionConsumer sarama.PartitionConsumer) {

	// Trap SIGINT to trigger a shutdown.
	// signals := make(chan os.Signal, 1)
	// signal.Notify(signals, os.Interrupt)
	done := make(chan bool)

	// Count how many message processed
	msgCount := 0

	start := time.Now()

	// event := new(OrderEvent)
	// var oe OrderEvent

	for msg := range partitionConsumer.Messages() {
		var event OrderEvent

		if err := jsoniter.Unmarshal(msg.Value, &event); err == nil {
			// 	// event.Order.Amount = 1234567890123.12345678
			// log.Info(event.Order)
			event.Order.Amount = 1
			engine.Run(event.Order)
		}

		msgCount++
		if msgCount%10000 == 0 {
			// log.Info("msgCount:", msgCount)
			// log.Info("msg.offset:", msg.Offset)
			// log.Info("msg.Value:", string(msg.Value))
			// log.Info("event:    ", event)

			elapsed := time.Since(start)
			tps := int(float64(msgCount) / elapsed.Seconds())

			log.Info("range Handled ", msgCount, " events in ", elapsed, " at ", tps, " events/second.")

			// break ConsumerLoop
		}
	}
	done <- true

}

func (engine *RBTEngine) consumeEvents(partitionConsumer sarama.PartitionConsumer) {

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	msgCount := 0
	start := time.Now()

	// event := new(OrderEvent)
	// var event OrderEvent

ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var event OrderEvent

			if err := jsoniter.Unmarshal(msg.Value, &event); err == nil {
				event.Order.Amount = 1
				// engine.Run(event.Order)

				msgCount++
				printMsgCounter(start, msgCount, msg, &event)

			} else {
				log.Info("jsoniter.Unmarshal error:", err)
				break ConsumerLoop
			}

		case <-signals:
			break ConsumerLoop

			// default:
			// 	fmt.Println("no message received")

		}
	}

}

// func printMsgCounter(start time.Time, msgCount int, msg *sarama.ConsumerMessage, event *OrderEvent) {
// 	if msgCount%10000 == 0 {
// 		// log.Info("msgCount:", msgCount)
// 		// log.Info("msg.offset:", msg.Offset)
// 		// log.Info("msg.Value:", string(msg.Value))
// 		// log.Info("event:    ", event)

// 		elapsed := time.Since(start)
// 		tps := int(float64(msgCount) / elapsed.Seconds())

// 		log.Info("Handled ", msgCount, " events in ", elapsed, " at ", tps, " events/second.")
// 	}
// }
