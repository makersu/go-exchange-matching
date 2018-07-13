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
		// var amount = uint32(rand.Int31n(maxAmount)) + 1
		var amount uint32 = 1
		orders = append(orders, &Order{
			id:   uint64(i) + 1,
			pair: BTC_ETH,
			// isBuy:  float64(price) >= priceMean,
			operation: randomOperation(price, priceMean),
			price:     int(price),
			// amount:    uint32(rand.Int31n(maxAmount)),
			amount:            amount,
			NumberOutstanding: amount,
			Timestamp:         time.Now().Nanosecond(),
		})
	}
	return orders
}

func randomOperation(price uint32, priceMean float64) Operation {
	if float64(price) >= priceMean {
		return BID
	} else {
		return ASK
	}
	// return ASK

}

func doPerfTest(n int, priceMean, priceStd float64, maxAmount int32) {
	fmt.Println("doPerfTest:", " n=", n, " priceMean=", priceMean, "priceStd=", priceStd, "maxAmount=", maxAmount)
	orders := buildOrders(n, priceMean, priceStd, maxAmount)

	engine := NewEngine(ADA_ETH)

	start := time.Now()
	for _, order := range orders {
		// myPrintln("order", order)
		// engine.addOrder(order)
		engine.Run(order)
	}

	elapsed := time.Since(start)

	fmt.Printf("\n\nHandled %v actions in %v at %v n/second.\n", n, elapsed, int(float64(n)/elapsed.Seconds()))

	// printResult(&engine)
	// printOrderbook(engine.book) //

}

// func printResult(engine *Engine) {
// 	fmt.Println("Result orderbook", engine.book)
// 	// book := engine.book
// 	// bookJson, _ := engine.book.ToJSON()
// 	// fmt.Println("error", error)
// 	// fmt.Println("b", string(bookJson))

// 	for _, k := range engine.book.Keys() {
// 		// fmt.Println("k", k)
// 		treeNode, _ := engine.book.Get(k)
// 		// fmt.Println("ok", ok)
// 		fmt.Println("\ntreeNode", treeNode) //nil??
// 	}

// }

func TestPerf(t *testing.T) {
	doPerfTest(10, 5000, 10, 50)
	// doPerfTest(10000, 5000, 1000, 5000)
	// doPerfTest(100000, 5000, 10, 50)
	// doPerfTest(100000, 5000, 1000, 5000)
	// doPerfTest(1000000, 5000, 10, 50)
	// doPerfTest(1000000, 5000, 1000, 5000)
}
