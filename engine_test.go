package engine

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"math"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"testing"
	"time"
)

func InitLog() {
	// log.SetFormatter(&log.JSONFormatter{})
	// log.SetFormatter(&log.TextFormatter{})
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006/01/02 - 15:04:05",
		FullTimestamp:   true,
	})

	log.SetLevel(log.DebugLevel)
}

func InitProfile() {
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

func buildOrders(totalAmount int, priceMean, priceStd float64, maxAmount int32) []*Order {
	orders := make([]*Order, 0)
	var price uint32
	for i := 0; i < totalAmount; i++ {
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
}

func doPerfTest(totalAmount int, priceMean, priceStd float64, maxOrderAmount int32) {
	log.Info("doPerfTest:", " totalAmount=", totalAmount, " priceMean=", priceMean, " priceStd=", priceStd, " maxOrderAmount=", maxOrderAmount)

	orders := buildOrders(totalAmount, priceMean, priceStd, maxOrderAmount)

	engine := NewEngine(ADA_ETH)

	start := time.Now()

	for _, order := range orders {
		engine.Run(order)
	}

	elapsed := time.Since(start)
	tps := int(float64(totalAmount) / elapsed.Seconds())

	log.Info("Handled ", totalAmount, " orders in ", elapsed, " at ", tps, " orders/second.")
	engine.printEngine() //refactor
}
func TestPerf(t *testing.T) {

	doPerfTest(10000, 5000, 10, 50)
	doPerfTest(10000, 5000, 1000, 5000)
	doPerfTest(100000, 5000, 10, 50)
	doPerfTest(100000, 5000, 1000, 5000)
	doPerfTest(1000000, 5000, 10, 50)
	doPerfTest(1000000, 5000, 1000, 5000)

	// doPerfTest(100, 5000, 10, 50)
	// doPerfTest(1000, 5000, 10, 50)
	// doPerfTest(10000, 5000, 10, 50)
	// doPerfTest(100000, 5000, 10, 50)
	// doPerfTest(500000, 5000, 10, 50)
	// doPerfTest(1000000, 5000, 10, 50)

	// doPerfTest(100, 5000, 1000, 5000)
	// doPerfTest(1000, 5000, 1000, 5000)
	// doPerfTest(10000, 5000, 1000, 5000)
	// doPerfTest(100000, 5000, 1000, 5000)
	// doPerfTest(500000, 5000, 1000, 5000)
	// doPerfTest(1000000, 5000, 1000, 5000)

	// doPerfTest(1000000, 5000, 100, 10000)

}
