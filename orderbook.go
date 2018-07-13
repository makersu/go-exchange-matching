package engine

import (
	"fmt"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
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

// //TODO: refactor for Json
// func (node TreeNode) String() string {
// 	fmt.Println("len(node.orders)", len(node.orders))

// 	vals := make([]*Order, len(node.orders))

// 	// fmt.Println("vals", vals)

// 	for _, v := range node.orders {
// 		// fmt.Println("v", v)
// 		vals = append(vals, v)
// 	}
// 	return fmt.Sprint(vals)
// }

func (n TreeNode) String() string {
	vals := make([]*Order, len(n.orders))
	for _, v := range n.orders {
		vals = append(vals, v)
	}
	return fmt.Sprint(vals)
}

func printOrderbook(book *rbt.Tree) {
	fmt.Println("\nprintOrderbook orderbook", book)
	fmt.Println("book.Keys()", book.Keys())

	for _, k := range book.Keys() {
		// fmt.Println("k", k)
		treeNode, _ := book.Get(k)
		// fmt.Println("ok", ok)
		fmt.Println("treeNode", treeNode) //nil??
	}
}
