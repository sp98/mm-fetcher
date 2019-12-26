package utility

import (
	"fmt"
	"strconv"
	"strings"

	kiteticker "github.com/zerodhatech/gokiteconnect/ticker"
)

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
