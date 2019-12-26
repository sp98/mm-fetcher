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
	KC              *kiteconnect.Client
	APIKey          string
	APISecret       string
	APIRequestToken string
	APIAccesToken   string
	Subscrptions    []uint32
	Stocks          [][]string
}

//NewTicker creates a new Ticker object
func NewTicker(apiKey, apiSecret, apiReqToken string, subs []uint32, stocks [][]string) *Ticker {
	return &Ticker{
		APIKey:          apiKey,
		APISecret:       apiSecret,
		APIRequestToken: apiReqToken,
		Subscrptions:    subs,
		Stocks:          stocks,
	}
}

//Connect to Kite API
func (t Ticker) Connect() error {
	kc := kiteconnect.New(t.APIKey)

	// Get user details and access token
	data, err := kc.GenerateSession(t.APIRequestToken, t.APISecret)
	if err != nil {
		return fmt.Errorf("error in generating kite session. %+v", err)
	}

	// Set access token
	t.APIAccesToken = data.AccessToken
	if t.APIAccesToken == "" {
		return fmt.Errorf("failed to get kite API access token")
	}
	kc.SetAccessToken(t.APIAccesToken)

	t.KC = kc
	return nil
}

//InitDB creates the Continues Queries for all the stocks across all the intervals
func (t Ticker) InitDB() error {

	DBUrl = os.Getenv("DB_URL")
	DBName = os.Getenv("DB_NAME")
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
			log.Println("Inside")
			err := db.CreateTickCQ(interval)
			if err != nil {
				return fmt.Errorf("error creating CQ for the isntrument: +%v. %+v", instrument, err)
			}
		}

	}
	return nil
}
