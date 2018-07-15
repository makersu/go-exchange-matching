package engine

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
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

	fmt.Printf("Handled %v orders in %v at %v orders/second.\n", n, elapsed, int(float64(n)/elapsed.Seconds()))

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

func init() {
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// ... rest of the program ...

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
}
func TestPerf(t *testing.T) {

	// doPerfTest(10000, 5000, 10, 50)
	// doPerfTest(10000, 5000, 1000, 5000)
	// doPerfTest(100000, 5000, 10, 50)
	// doPerfTest(100000, 5000, 1000, 5000)
	// doPerfTest(1000000, 5000, 10, 50)
	// doPerfTest(1000000, 5000, 1000, 5000)

	// doPerfTest(100, 5000, 10, 50)
	// doPerfTest(1000, 5000, 10, 50)
	// doPerfTest(10000, 5000, 10, 50)
	// doPerfTest(100000, 5000, 10, 50)
	// doPerfTest(200000, 5000, 10, 50)
	// doPerfTest(300000, 5000, 10, 50)
	// doPerfTest(400000, 5000, 10, 50)
	// doPerfTest(500000, 5000, 10, 50)
	// doPerfTest(600000, 5000, 10, 50)
	// doPerfTest(700000, 5000, 10, 50)
	// doPerfTest(800000, 5000, 10, 50)
	// doPerfTest(900000, 5000, 10, 50)
	// doPerfTest(1000000, 5000, 10, 50)

	// doPerfTest(1000, 5000, 1000, 5000)
	// doPerfTest(100000, 5000, 1000, 5000)
	// doPerfTest(200000, 5000, 1000, 5000)
	// doPerfTest(300000, 5000, 1000, 5000)
	// doPerfTest(400000, 5000, 1000, 5000)
	// doPerfTest(500000, 5000, 1000, 5000)
	// doPerfTest(600000, 5000, 1000, 5000)
	// doPerfTest(700000, 5000, 1000, 5000)
	// doPerfTest(800000, 5000, 1000, 5000)
	// doPerfTest(900000, 5000, 1000, 5000)
	// doPerfTest(1000000, 5000, 1000, 5000)

}
