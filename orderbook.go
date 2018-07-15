package engine

import (
	"fmt"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
	log "github.com/sirupsen/logrus"
)

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

//TODO: refactor to Json format
func (n *TreeNode) String() string {
	vals := make([]*Order, len(n.orders))
	for _, v := range n.orders {
		vals = append(vals, v)
	}
	return fmt.Sprint(vals)
}

func printOrderbook(book *rbt.Tree) {
	log.Debug("printOrderbook orderbook ", book)
	log.Debug("printOrderbook book.Keys() ", book.Keys())

	for _, k := range book.Keys() {
		treeNode, _ := book.Get(k)
		log.Debug("printOrderbook treeNode ", treeNode)
	}
	log.Debug("printOrderbook\n")

}

// func (n *TreeNode) delete(id uuid.UUID) {
// func (n *TreeNode) delete(id uint64) {
// 	delete(n.orders, id)
// 	if len(n.orders) == 0 {
// 		myPrintln("len(n.orders)==0", n)
// 	}
// 	//if(n.orders)
// }
