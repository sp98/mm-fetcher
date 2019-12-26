package main

import (
	"log"
	"os"
	"time"

	"github.com/fetcher/pkg/ticker"
	"github.com/fetcher/pkg/utility"
)

const (
	//APIKey is the Zeoradha Kite connect API Key
	APIKey = "API_KEY"
	//APISecret is the Zeoradha Kite connect API Secret
	APISecret = "API_SECRET"
	//APIRequestToken is the Zeoradha Kite connect API Request Token
	APIRequestToken = "API_REQUESTTOKEN"
	//Stocks are the share markets stocks to fetch ticks from
	Stocks = "STOCKS"
	//IsTest runs fetcher with fake data for testing
	IsTest = "TEST_WITH_FAKE_DATA"
)

func main() {
	log.Println("Starting Fetcher")
	isTest := os.Getenv(IsTest)
	if len(isTest) > 0 {
		dummySetup()
	} else {
		setup()
	}
}

func setup() {

	tkr := ticker.NewTicker(os.Getenv("API_KEY"), os.Getenv("API_SECRET"), os.Getenv("API_REQUESTTOKEN"),
		utility.GetSubscriptions(os.Getenv("Stocks")), utility.GetStocks(os.Getenv("Stocks")))
	err := tkr.Connect()
	if err != nil {
		log.Fatalf("failed connecting to Kite API. %+v", err)
		return
	}
	err = tkr.InitDB()
	if err != nil {
		log.Fatalf("failed intializding Database to Kite API. %+v", err)
		return
	}
	log.Printf("Ticker - %+v", tkr)
	tkr.Start()

}

func dummySetup() {
	//Run with dummy data when market is closed!
	testTicker := ticker.NewTicker(os.Getenv(APIKey), os.Getenv(APISecret), os.Getenv(APIRequestToken),
		utility.GetSubscriptions(os.Getenv(Stocks)), utility.GetStocks(os.Getenv(Stocks)))
	log.Println("Ticker object - ", testTicker)
	err := testTicker.InitDB()
	if err != nil {
		log.Println("error- ", err)
		return
	}
	for i := 0; i < 1000; i++ {
		time.Sleep(2 * time.Second)
		ticks := ticker.DummyTicks()
		log.Printf("Ticks - %+v", ticks)
		ticker.StoreTickInDB(ticks)
	}
}
