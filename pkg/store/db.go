package store

import (
	"fmt"
	"log"
	"time"

	"github.com/fetcher/pkg/utility"

	client "github.com/orourkedd/influxdb1-client/client"
	kiteticker "github.com/zerodhatech/gokiteconnect/ticker"
)

//DB is the influx db struct
type DB struct {
	Address     string
	Name        string
	Measurement string
}

//NewDB returns instance of an InfluxDB struct
func NewDB(address string, name string, measurement string) *DB {
	return &DB{
		Address:     address,
		Name:        name,
		Measurement: measurement,
	}

}

//GetClient creates a new Influx DB client
func (db *DB) GetClient() (client.Client, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: db.Address,
	})
	if err != nil {
		log.Fatalln("Error on creating Influx DB client: ", err)
		return nil, err
	}
	return c, nil
}

func (db DB) executeQuery(query client.Query) (*client.Response, error) {
	dbClient, _ := db.GetClient()
	defer dbClient.Close()

	response, err := dbClient.Query(query)

	if err != nil && response.Error() != nil {
		log.Fatalln("Error executing Query - ", err)
		return nil, err
	}

	return response, nil

}

func (db DB) executePointWrite(bp client.BatchPoints, measurement string, tags map[string]string, fields map[string]interface{}, t time.Time) error {
	dbClient, _ := db.GetClient()
	defer dbClient.Close()
	pt, err := client.NewPoint(measurement, tags, fields, t)
	if err != nil {
		log.Fatalln("Error creating new point - ", err)
		return nil
	}

	bp.AddPoint(pt)

	// Write the batch
	if err := dbClient.Write(bp); err != nil {
		log.Fatalln("Error writing batch - ", err)
		return nil
	}
	return nil
}

//CreateDB creates a new database
func (db DB) CreateDB() error {
	dbClient, _ := db.GetClient()
	defer dbClient.Close()
	query := client.Query{
		Command: fmt.Sprintf(createDB, db.Name),
	}
	_, err := db.executeQuery(query)
	if err != nil {
		//TODO: Handle dabase already exists error.
		return fmt.Errorf("Error creating database. %+v", err)
	}
	return nil
}

//CurrentDate returns the date in a specified format
func CurrentDate(format string) string {
	currentTime := time.Now()
	return currentTime.Format(format)
}

// CreateTickCQ creates a continuous query on Tick Measurement.
func (db DB) CreateTickCQ(tradeInterval string) error {
	cqMeasurement := fmt.Sprintf("%s_%s", db.Measurement, tradeInterval)
	query := fmt.Sprintf(tickCQTime, cqMeasurement, db.Measurement, tradeInterval)
	cquery := fmt.Sprintf(tickCQ, cqMeasurement, db.Name, query)
	q := client.NewQuery(cquery, db.Name, "")
	_, err := db.executeQuery(q)
	if err != nil {
		return fmt.Errorf("error creating continuous query. %+v ", err)
	}
	return nil

}

//StoreTick saves tick data in influx db
func (db DB) StoreTick(tickData *kiteticker.Tick) error {

	buyDepth, SellDepth := utility.ParseDepth(tickData.Depth)
	dbClient, _ := db.GetClient()
	defer dbClient.Close()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db.Name,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	tick := *tickData
	fields := map[string]interface{}{
		"LastPrice":          tick.LastPrice,
		"LastTradedQuantity": tick.LastTradedQuantity,
		"TotalBuyQuantity":   tick.TotalBuyQuantity,
		"TotalSellQuantity":  tick.TotalSellQuantity,
		"VolumeTraded":       tick.VolumeTraded,
		"TotalBuy":           tick.TotalBuy,
		"TotalSell":          tick.TotalSell,
		"AverageTradePrice":  tick.AverageTradePrice,
		"BuyDepth":           buyDepth,
		"SellDepth":          SellDepth,
		"Open":               tick.OHLC.Open,
		"High":               tick.OHLC.High,
		"Low":                tick.OHLC.Low,
		"Close":              tick.OHLC.Close,
	}
	tags := map[string]string{
		// "InstrumentToken": fmt.Sprint(tick.InstrumentToken),
	}

	err = db.executePointWrite(bp, db.Measurement, tags, fields, tickData.Timestamp.Time)
	if err != nil {
		return fmt.Errorf("Error storing tick data to db. %+v", err)
	}
	return nil

}
