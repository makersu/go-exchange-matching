package engine

//TODO: refactor type of
type TreeNode struct {
	depth int //
	// orders map[uuid.UUID]Order
	orders map[uint64]*Order
}

func NewTreeNode() *TreeNode {
	return &TreeNode{
		depth: 0,
		// orders: make(map[uuid.UUID]Order),
		orders: make(map[uint64]*Order),
	}
}

//TODO: refactor for rename
func (node *TreeNode) upsert(order *Order) {
	node.orders[order.id] = order
}

// func (node TreeNode) String() string {
// 	vals := make([]*Order, len(node.orders))
// 	for _, v := range node.orders {
// 		// fmt.Println("v", v)
// 		vals = append(vals, v)
// 	}
// 	return fmt.Sprint(vals)
// }
