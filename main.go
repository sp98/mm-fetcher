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

	//APIAccessToken is the access token to connect to kite API
	APIAccessToken = "API_ACCESSTOKEN"

	//Stocks are the share markets stocks to fetch ticks from
	Stocks = "STOCKS"
	//IsTest runs fetcher with fake data for testing
	IsTest = "TEST_WITH_FAKE_DATA"
)

func main() {
	isTest := os.Getenv(IsTest)
	if isTest == "true" {
		log.Println("Starting Fetcher with fake data")
		utility.IsMarketOpen()
		dummySetup()
	} else {
		log.Println("Starting Fetcher")
		utility.IsMarketOpen()
		setup()
	}

}

func setup() {
	go utility.IsMarketClosed()
	tkr := ticker.NewTicker(os.Getenv(APIKey), os.Getenv(APISecret), os.Getenv(APIRequestToken), os.Getenv(APIAccessToken),
		utility.GetSubscriptions(os.Getenv(Stocks)), utility.GetStocks(os.Getenv(Stocks)))
	err := tkr.Connect()
	if err != nil {
		log.Fatalf("failed connecting to Kite API. %+v", err)
		return
	}

	log.Printf("Ticker - %+v", tkr)

	err = tkr.InitDB()
	if err != nil {
		log.Fatalf("failed intializding Database to Kite API. %+v", err)
		return
	}
	tkr.Start()

}

func dummySetup() {
	go utility.IsMarketClosed()
	//Run with dummy data when market is closed!
	testTicker := ticker.NewTicker(os.Getenv(APIKey), os.Getenv(APISecret), os.Getenv(APIRequestToken), os.Getenv(APIAccessToken),
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
