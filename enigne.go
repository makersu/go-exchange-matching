package engine

import (
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
)

type ExchangePair int

const (
	ADA_ETH  ExchangePair = iota // value: 1, type: ExchangePair
	BTC_ETH                      // value: 2, type: ExchangePair
	CVC_ETH                      // value: 3, type: ExchangePair
	DASH_ETH                     // value: 4, type: ExchangePair
	INVALID_EXCHANGE_PAIR
)

type Engine struct {
	pair    ExchangePair
	askbook *rbt.Tree
	bidbook *rbt.Tree
}

// NewEngine is the Matching Engine constructor
func NewEngine(pair ExchangePair) Engine {
	return Engine{
		pair:    pair,
		askbook: rbt.NewWithIntComparator(),
		bidbook: rbt.NewWithIntComparator(),
		// writeLog: NewWriteLog(exchange.String()),
	}
}

//addOrder adds an order that cannot be filled any further to the orderbook
// func (engine *Engine) addOrder(order *Order) *Order {
func addOrder(book *rbt.Tree, order *Order) {

	treeNode, ok := book.Get(order.price)
	if !ok {
		node := NewTreeNode()
		node.upsert(order)
		book.Put(order.price, node)

	} else {
		node := treeNode.(*TreeNode)
		node.upsert(order)
		book.Put(order.price, treeNode)
	}

}

func addAskOrder(sellbook *rbt.Tree, order *Order) {
	addOrder(sellbook, order)
}

func addBidOrder(buybook *rbt.Tree, order *Order) {
	addOrder(buybook, order)
}

//executeOrder walk the orderbook and match asks and bids that can fill
// func (engine *Engine) executeOrder(book *rbt.Tree, executingOrder *Order) (*Order, []Match) {
func (engine *Engine) executeOrder(executingOrder *Order) (*Order, []Match) {

	var matches []Match

	// bid for lowest price
	if executingOrder.operation == BID {
		log.Debug("executeOrder BID")

		askbookit := engine.askbook.Iterator()

		// get the begin node(lowest price) then next
		for askbookit.Begin(); askbookit.Next(); {
			nodePrice, node := askbookit.Key().(int), askbookit.Value().(*TreeNode)
			log.Debug("executeOrder executingOrder price ", executingOrder.price, executingOrder) //
			log.Debug("executeOrder BID askbook.Iterator()  ", nodePrice, node)                   //

			//Check price
			if nodePrice <= executingOrder.price {

				log.Debug("executeOrder BID askbook.Iterator() nodePrice <= executingOrder.price ", nodePrice, executingOrder.price) //
				// nodeOrderResult, nodeFills := matchNode(node, ord)
				_, nodeMatches := engine.matchNode(engine.askbook, node, executingOrder)

				//TODO nodeMatch.Number = 0
				for _, nodeMatch := range nodeMatches {
					if nodeMatch.Number > 0 {
						matches = append(matches, nodeMatch)
					}
				}

			} else {
				log.Debug("executeOrder BID askbook.Iterator() nodePrice > executingOrder.price ", nodePrice, executingOrder.price) //
				//skip this node, too expensive (The cheapest ask could be higher than this bid)
				// continue
				break

			}

			if executingOrder.NumberOutstanding == 0 {
				// if we have 0 outstanding we can quit
				break
			}

		}
		return executingOrder, matches
	} else if executingOrder.operation == ASK {
		log.Debug("executeOrder ASK")

		bidbookit := engine.bidbook.Iterator()

		// get the end element(highest) then previous
		for bidbookit.End(); bidbookit.Prev(); {
			nodePrice, node := bidbookit.Key().(int), bidbookit.Value().(*TreeNode)

			log.Debug("executeOrder executingOrder price ", executingOrder.price, executingOrder) //
			log.Debug("executeOrder ASK book.Iterator()  ", nodePrice, node)                      //

			//Check price to sell high?
			if nodePrice >= executingOrder.price {
				log.Debug("executeOrder ASK book.Iterator() nodePrice >= executingOrder.price ", nodePrice, executingOrder.price) //

				// nodeOrderResult, nodeFills := matchNode(node, ord)
				_, nodeMatches := engine.matchNode(engine.bidbook, node, executingOrder)

				log.Debug("executeOrder matchNode nodeMatches", nodeMatches)

				//TODO nodeMatch.Number = 0
				for _, nodeMatche := range nodeMatches {
					if nodeMatche.Number > 0 {
						matches = append(matches, nodeMatche)
					}
				}

			} else {
				log.Debug("executeOrder ASK book.Iterator() nodePrice < executingOrder.price ", nodePrice, executingOrder.price) //
				//skip this node, too expensive (The cheapest ask could be higher than this bid)
				// continue
				break
			}

			if executingOrder.NumberOutstanding == 0 {
				// if we have 0 outstanding we can quit
				break
			}

		}
		return executingOrder, matches
	} else {
		// Not a valid bid/ask
	}

	return &Order{}, nil
}

func (engine *Engine) Run(order *Order) {
	switch order.operation {
	case ASK:
		log.Debug("*Run ASK for Highest price ", order.price, order)

		// Ask Operations
		executedOrder, matches := engine.executeOrder(order)
		log.Debug("*Run ASK matches", matches)

		if executedOrder.NumberOutstanding > 0 {
			addOrder(engine.askbook, order)
		}

	case BID:
		log.Debug("*Run BID for Lowest price ", order.price, order)

		// Bid Operations
		executedOrder, matches := engine.executeOrder(order)
		log.Debug("*Run BID matches", matches)

		if executedOrder.NumberOutstanding > 0 {
			addBidOrder(engine.bidbook, order)
		}

	case CANCEL:
		//Cancel an order
		// fill := cancelOrder(d.book, order.ID)

		// fmt.Println("CANCEL fill", fill)
		// d.writeLog.logFill(fill)
		// out <- fill

	default:
		//Drop the message
		fmt.Println("Invalid Order Type")
	}
	// }
	// log.Debug("printOrderbook ASKbook")
	// printOrderbook(engine.askbook) //
	// log.Debug("printOrderbook BIDbook")
	// printOrderbook(engine.bidbook) //

}

//matchNode takes an order and fills it against a node, NOT IDEMPOTENT
func (engine *Engine) matchNode(book *rbt.Tree, node *TreeNode, matchingOrder *Order) (*Order, []Match) {

	//TODO
	//We only deal with ask and bid
	if matchingOrder.operation == CANCEL || matchingOrder.operation == INVALID_OPERATION {
		return matchingOrder, []Match{}
	}

	// orders := node.sortedOrders()
	orders := node.orders

	activeOrder := matchingOrder //?
	var matches []Match

	for _, oldOrder := range orders {
		// if activeOrder.operation != oldOrder.operation {
		// log.Debug("matchNode activeOrder.operation != oldOrder.operation")
		// If the current order can fill new order
		if oldOrder.NumberOutstanding >= matchingOrder.NumberOutstanding {
			log.Debug("matchNode oldOrder.NumberOutstanding >= matchingOrder.NumberOutstanding ", oldOrder.NumberOutstanding, matchingOrder.NumberOutstanding)
			partialFill := []*Order{activeOrder, oldOrder}
			closed := []*Order{activeOrder}

			if oldOrder.NumberOutstanding-matchingOrder.NumberOutstanding == 0 {
				log.Debug("matchNode oldOrder.NumberOutstanding - matchingOrder.NumberOutstanding == 0")
				closed = append(closed, oldOrder)
				// node.delete(oldOrder.id)
				engine.removeOrder(book, node, oldOrder)
				nodeMatch := NewMatch(activeOrder.pair, activeOrder.NumberOutstanding, oldOrder.price, partialFill, closed)

				//Order is filled
				activeOrder.NumberOutstanding = 0
				matches = append(matches, nodeMatch)

			} else { // Update old order
				log.Debug("matchNode oldOrder.NumberOutstanding - matchingOrder.NumberOutstanding != 0")

				oldRemaining := oldOrder.NumberOutstanding - activeOrder.NumberOutstanding
				oldOrder.NumberOutstanding = oldRemaining

				nodeMatch := NewMatch(activeOrder.pair, activeOrder.NumberOutstanding, oldOrder.price, partialFill, closed)

				//Order is matched
				activeOrder.NumberOutstanding = 0
				matches = append(matches, nodeMatch)

				node.upsert(oldOrder)
			}

		} else { // If the current order is too small to fill the new order

			log.Debug("matchNode oldOrder.NumberOutstanding < matchingOrder.NumberOutstanding")

			// node.delete(oldOrder.id)
			engine.removeOrder(book, node, oldOrder)

			partialFill := []*Order{activeOrder, oldOrder}
			closed := []*Order{oldOrder}
			nodeMatch := NewMatch(activeOrder.pair, oldOrder.NumberOutstanding, oldOrder.price, partialFill, closed)

			activeOrder.NumberOutstanding = activeOrder.NumberOutstanding - oldOrder.NumberOutstanding
			matches = append(matches, nodeMatch)

		}

		// } else {
		// 	log.Debug("matchNode activeOrder.operation == oldOrder.operation")
		// }
	}

	return activeOrder, matches
}

//TODO: linked list?
//TODO: sort after add?
func (n TreeNode) sortedOrders() []*Order {
	log.Debug("sortedOrders")
	orders := make([]*Order, 0)
	for _, v := range n.orders {
		orders = append(orders, v)
	}

	sort.Slice(orders[:], func(i, j int) bool {
		return orders[i].Timestamp < orders[j].Timestamp
	})

	return orders
}

// func (n *TreeNode) delete(id uint64) {
func (engine *Engine) removeOrder(book *rbt.Tree, n *TreeNode, order *Order) {
	delete(n.orders, order.id)
	if len(n.orders) == 0 {
		// myPrintln("len(n.orders)==0", n)
		book.Remove(order.price)
	}

}

func init() {
	// log.SetFormatter(&log.JSONFormatter{})
	// log.SetFormatter(&log.TextFormatter{})
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006/01/02 - 15:04:05",
		FullTimestamp:   true,
	})

	log.SetLevel(log.InfoLevel)

}
