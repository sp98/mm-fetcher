package ticker

import (
	"fmt"
	"log"
	"os"

	"github.com/fetcher/pkg/store"

	kiteconnect "github.com/zerodhatech/gokiteconnect"
)

var (
	//DBUrl is the connection string for influx DB
	DBUrl = ""
	//DBName is the database name to store tick data
	DBName = ""
	//Intervals for storing stock data. Options 5m, 10m 15m
	Intervals = []string{"5m"}
)

//Ticker fetches Tick data
type Ticker struct {
	KC            *kiteconnect.Client
	APIKey        string
	APISecret     string
	APIAccesToken string
	Subscrptions  []uint32
	Stocks        [][]string
}

//NewTicker creates a new Ticker object
func NewTicker(apiKey, apiSecret, apiAccessToken string, subs []uint32, stocks [][]string) *Ticker {
	return &Ticker{
		APIKey:        apiKey,
		APISecret:     apiSecret,
		APIAccesToken: apiAccessToken,
		Subscrptions:  subs,
		Stocks:        stocks,
	}
}

//Connect to Kite API
func (t Ticker) Connect() error {
	kc := kiteconnect.New(t.APIKey)
	if t.APIAccesToken == "" {
		return fmt.Errorf("empty api access token env variable")
	}
	kc.SetAccessToken(t.APIAccesToken)

	//Set Kite Ticker client
	t.KC = kc
	return nil
}

//InitDB creates the Continues Queries for all the stocks across all the intervals
func (t Ticker) InitDB() error {

	DBUrl = os.Getenv("INFLUX_DB_URL")
	DBName = os.Getenv("TICK_STORE_DB_NAME")
	db := store.NewDB(DBUrl, DBName, "")
	err := db.CreateDB()
	if err != nil {
		return err
	}

	//Instrument Name, Sybmol, Token, Exchange, Interval

	//Create continuous queries
	for _, instrument := range t.Stocks {
		log.Println("Instrument - ", instrument)
		db := store.NewDB(DBUrl, DBName, "")
		db.Measurement = fmt.Sprintf("%s_%s", "ticks", instrument[2])
		for _, interval := range Intervals {
			err := db.CreateTickCQ(interval)
			if err != nil {
				return fmt.Errorf("error creating CQ for the isntrument: +%v. %+v", instrument, err)
			}
		}

	}
	return nil
}
