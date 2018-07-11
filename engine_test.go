package engine

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func buildOrders(n int, priceMean, priceStd float64, maxAmount int32) []*Order {
	orders := make([]*Order, 0)
	var price uint32
	for i := 0; i < n; i++ {
		price = uint32(math.Abs(rand.NormFloat64()*priceStd + priceMean))
		orders = append(orders, &Order{
			id:     uint64(i) + 1,
			isBuy:  float64(price) >= priceMean,
			price:  int(price),
			amount: uint32(rand.Int31n(maxAmount)),
		})
	}
	return orders
}

func doPerfTest(n int, priceMean, priceStd float64, maxAmount int32) {
	fmt.Println("doPerfTest:", " n=", n, " priceMean=", priceMean, "priceStd=", priceStd, "maxAmount=", maxAmount)
	orders := buildOrders(n, priceMean, priceStd, maxAmount)

	engine := NewEngine(ADA_ETH)

	start := time.Now()
	for _, order := range orders {
		myPrintln("order", order)
		engine.addOrder(order)
	}

	elapsed := time.Since(start)

	fmt.Printf("Handled %v actions in %v at %v n/second.\n", n, elapsed, int(float64(n)/elapsed.Seconds()))

	// printResult(engine)

}

func printResult(engine Engine) {
	fmt.Println("Result orderbook", engine.book)
	// book := engine.book
	// b, _ := engine.book.ToJSON()
	// fmt.Println("error", error)
	// fmt.Println("b", string(b))

	// for _, k := range engine.book.Keys() {
	// 	// fmt.Println("k", k)
	// 	treeNode, _ := engine.book.Get(k)
	// 	// fmt.Println("ok", ok)
	// 	fmt.Println("treeNode", treeNode)
	// }
}

func TestPerf(t *testing.T) {
	doPerfTest(10000, 5000, 10, 50)
	doPerfTest(10000, 5000, 1000, 5000)
	doPerfTest(100000, 5000, 10, 50)
	doPerfTest(100000, 5000, 1000, 5000)
	doPerfTest(1000000, 5000, 10, 50)
	doPerfTest(1000000, 5000, 1000, 5000)
}
