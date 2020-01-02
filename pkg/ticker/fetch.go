package ticker

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/fetcher/pkg/store"
	kiteconnect "github.com/zerodhatech/gokiteconnect"
	kiteticker "github.com/zerodhatech/gokiteconnect/ticker"
)

var (
	ticker        *kiteticker.Ticker
	subscriptions []uint32
)

// Triggered when any error is raised
func onError(err error) {
	fmt.Println("Error: ", err)
}

// Triggered when websocket connection is closed
func onClose(code int, reason string) {
	fmt.Println("Close: ", code, reason)
}

// Triggered when connection is established and ready to send and accept data
func onConnect() {
	fmt.Println("Connection Status: Connected")
	err := ticker.Subscribe(subscriptions)
	if err != nil {
		fmt.Println("err: ", err)
	}
	log.Println("Mode: ", kiteticker.ModeFull)
	ticker.SetMode(kiteticker.ModeFull, subscriptions)
}

// Triggered when tick is recevived
func onTick(tick kiteticker.Tick) {
	log.Printf("Tick: %+v", tick)
	StoreTickInDB(&tick)
}

// Triggered when reconnection is attempted which is enabled by default
func onReconnect(attempt int, delay time.Duration) {
	fmt.Printf("Reconnect attempt %d in %fs\n", attempt, delay.Seconds())
}

// Triggered when maximum number of reconnect attempt is made and the program is terminated
func onNoReconnect(attempt int) {
	fmt.Printf("Maximum no of reconnect attempt reached: %d", attempt)
}

// Triggered when order update is received
func onOrderUpdate(order kiteconnect.Order) {
	fmt.Printf("Order: %+v", order.OrderID)
}

//Start creates a tick connection to retrieve tick data from Zerodha Kite API
func (t Ticker) Start() {
	// Create new Kite ticker instance
	ticker = kiteticker.New(t.APIKey, t.APIAccesToken)
	subscriptions = t.Subscrptions

	// Assign callbacks
	ticker.OnError(onError)
	ticker.OnClose(onClose)
	ticker.OnConnect(onConnect)
	ticker.OnReconnect(onReconnect)
	ticker.OnNoReconnect(onNoReconnect)
	ticker.OnTick(onTick)
	ticker.OnOrderUpdate(onOrderUpdate)

	// Start the connection
	ticker.Serve()
}

//StoreTickInDB stors the tick in influx db
func StoreTickInDB(tick *kiteticker.Tick) {
	db := store.NewDB(DBUrl, DBName, "")
	db.Measurement = fmt.Sprintf("%s_%s", "ticks", strconv.FormatUint(uint64(tick.InstrumentToken), 10))
	db.StoreTick(tick)

}
