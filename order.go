package engine

import "fmt"

/*
 * bexchange
 */

type OrderStatus int

type Operation int

const (
	ASK    Operation = iota // value: 1, type: BookOperation
	BID                     // value: 2, type: BookOperation
	CANCEL                  // value: 3, type: BookOperation
	STATUS
	INVALID_OPERATION
)

// TODO: refactor types of id, price and amount
type Order struct {
	id   uint64
	pair ExchangePair `json:"exchange"` // The exchange either BTC/USD, BTC/LTC, BTC/Doge, BTC/XMR(Monero)
	// isBuy bool
	operation Operation
	// price  uint32
	price             int
	amount            uint32
	NumberOutstanding uint32
	status            OrderStatus
	// next   *Order //
	Timestamp int `json:"timestamp"` // timestamp in nanoseconds
}

//TODO: refactor to Json
func (o *Order) String() string {
	return fmt.Sprintf("Order{id:%v,operation:%v,price:%v,amount:%v,NumberOutstanding:%v,Timestamp:%v}",
		o.id, o.operation, o.price, o.amount, o.NumberOutstanding, o.Timestamp)
}
