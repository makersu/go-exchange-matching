package engine

import (
	"time"
)

type Match struct {
	// ID           uuid.UUID    `json:"id"`
	ID           string       `json:"id"`
	pair         ExchangePair `json:"pair"`
	Number       uint32       `json:"number"`
	Price        int          `json:"price"`
	Timestamp    int          `json:"timestamp"`
	Participants []*Order     `json:"participants"`
	Closed       []*Order     `json:"closed"`
}

func NewMatch(pair ExchangePair, num uint32, price int, part []*Order, closed []*Order) Match {
	// uid, _ := uuid.NewV4()

	return Match{
		// ID:           uid.String(),
		ID:           string(time.Now().Nanosecond()),
		pair:         pair,
		Number:       num,
		Price:        price,
		Timestamp:    time.Now().Nanosecond(),
		Participants: part,
		Closed:       closed}
}

// type Message struct {
// 	Kind     MsgKind `json:"kind"`
// 	Price    uint64  `json:"price"`
// 	Amount   uint64  `json:"amount"`
// 	StockId  uint64  `json:"stockId"`
// 	TraderId uint32  `json:"traderId"`
// 	TradeId  uint32  `json:"tradeId"`
// }

// func (m *Message) String() string {
// 	return fmt.Sprintf("\nMessage{Kind:%v,Price:%v,Amount:%v,StockId:%v, TraderId:%v, TradeId:%v}",
// 		m.Kind, m.Price, m.Amount, m.StockId, m.TraderId, m.TradeId)
// }

// type MsgKind uint64

// const (
// 	NO_KIND       = MsgKind(iota)
// 	BUY           = MsgKind(iota)
// 	SELL          = MsgKind(iota)
// 	CANCEL        = MsgKind(iota)
// 	PARTIAL       = MsgKind(iota)
// 	FULL          = MsgKind(iota)
// 	CANCELLED     = MsgKind(iota)
// 	NOT_CANCELLED = MsgKind(iota)
// 	REJECTED      = MsgKind(iota)
// 	SHUTDOWN      = MsgKind(iota)
// 	NEW_TRADER    = MsgKind(iota)
// 	NUM_OF_KIND   = int(iota)
// )

// func (m *Message) WriteCancelFor(om *Message) {
// 	*m = *om
// 	m.Kind = CANCEL
// }
