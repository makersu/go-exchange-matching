package engine

import (
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
	pair ExchangePair
	book *rbt.Tree
}

// NewEngine is the Matching Engine constructor
func NewEngine(pair ExchangePair) Engine {
	return Engine{
		pair: pair,
		book: rbt.NewWithIntComparator(),
		// writeLog: NewWriteLog(exchange.String()),
	}
}

//addOrder adds an order that cannot be filled any further to the orderbook
func (engine *Engine) addOrder(order *Order) *Order {
	book := engine.book

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

	return order
}
