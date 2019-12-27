package utility

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	kiteticker "github.com/zerodhatech/gokiteconnect/ticker"
)

const (
	MarketOpenTime       = "%s 9:00:00"
	MarketCloseTime      = "%s 15:30:00"
	MarketActualOpenTime = "%s 09:13:00 MST"
	TstringFormat        = "2006-01-02 15:04:05"
	LayOut               = "2006-01-02 15:04:05"
	InfluxLayout         = "2006-01-02T15:04:05Z"
)

//IsWithInMarketOpenTime tells whether current time is withing market time and not on weekends
func IsWithInMarketOpenTime() (bool, error) {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	motString := fmt.Sprintf(MarketOpenTime, time.Now().Format("2006-01-02"))
	mot, err := time.ParseInLocation("2006-01-02 15:04:05", motString, loc)
	if err != nil {
		return false, fmt.Errorf("error parsing market open time. %+v", err)
	}

	mctString := fmt.Sprintf(MarketCloseTime, time.Now().Format("2006-01-02"))
	mct, err := time.ParseInLocation("2006-01-02 15:04:05", mctString, loc)
	if err != nil {
		return false, fmt.Errorf("error parsing market open time. %+v", err)
	}

	currentTime := time.Now()
	if currentTime.After(mot) && currentTime.Before(mct) && currentTime.Weekday() != 6 && currentTime.Weekday() != 7 {
		return true, nil
	}
	return false, nil

}

//IsMarketOpen is an infinite loop that runs untill the market is open
func IsMarketOpen() {
	for {
		t, _ := IsWithInMarketOpenTime()
		if t {
			log.Println(" withhin market open time.")
			break
		}
		log.Println("market is not open yet. waiting... ")
		time.Sleep(10 * time.Second)
	}

}

//IsMarketClosed is finite loop that panics when the market is closed
func IsMarketClosed() {
	for {
		t, _ := IsWithInMarketOpenTime()
		if !t {
			log.Println("stop fetcher as market is closed.")
			panic(1)
		}
		log.Println("market is still open. keep fetching...")
		time.Sleep(3 * time.Second)
	}
}

//GetSubscriptions returns the list of subscription IDs
func GetSubscriptions(stock string) []uint32 {
	token := []uint32{}

	stocks := strings.Split(stock, ",")

	for _, s := range stocks {
		sSlice := strings.Split(s, ";")
		if len(sSlice) >= 4 {
			token = append(token, getUnit32(sSlice[2]))
		}
	}

	return token
}

//GetStocks returns formated 2-D array of stocks
func GetStocks(stock string) [][]string {
	result := [][]string{}
	stocks := strings.Split(stock, ",")

	for _, s := range stocks {
		sSlice := strings.Split(s, ";")
		result = append(result, sSlice)
	}

	return result
}

//ParseDepth parses the market depth into string
func ParseDepth(depth kiteticker.Depth) (string, string) {
	var buyDepth string
	var sellDepth string
	for _, bd := range depth.Buy {
		buyDepth = fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v", buyDepth, ";", bd.Price, ";", bd.Quantity, ";", bd.Orders, ",")
	}

	for _, sd := range depth.Sell {
		sellDepth = fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v", buyDepth, ";", sd.Price, ";", sd.Quantity, ";", sd.Orders, ",")
	}

	return buyDepth, sellDepth
}
func getUnit32(str string) uint32 {
	u, _ := strconv.ParseUint(str, 10, 32)
	return uint32(u)
}
