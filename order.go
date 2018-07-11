package engine

/*
 * bexchange
 */

type OrderStatus int

// TODO: refactor types of id, price and amount
type Order struct {
	id    uint64
	isBuy bool
	// price  uint32
	price  int
	amount uint32
	status OrderStatus
	// next   *Order //
}

// func (o *Order) String() string {
// 	return fmt.Sprintf("\nOrder{id:%v,isBuy:%v,price:%v,amount:%v}",
// 		o.id, o.isBuy, o.price, o.amount)
// }
