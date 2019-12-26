package store

var (
	createDB = "CREATE DATABASE %s"
	tickCQ   = "CREATE CONTINUOUS QUERY %s ON %s BEGIN %s END"

	orderQuery = "SELECT * FROM Orders WHERE TradeDate=~/%s/"
	tradeQuery = `select * from trade where InstrumentToken='%s' ORDER BY time DESC limit 1`

	firstCandleStickQuery = `select * from %s limit 1`
	maxHighQuery          = "SELECT max(High) as Highest from %s"
	minLowQuery           = "SELECT min(Low) as Lowest from %s"
	ticksQuery            = "SELECT * FROM %s ORDER BY time DESC"
	ohlcQuery             = "SELECT * FROM %s ORDER BY time DESC limit 1"
	tickCQTime            = `SELECT FIRST(LastPrice) as Open, MAX(LastPrice) as High, MIN(LastPrice) as Low, LAST(LastPrice) as Close INTO %s FROM %s GROUP BY time(%s)`
)

/*
	tickCQTime            = `SELECT FIRST(LastPrice) as Open, MAX(LastPrice) as High,
							 MIN(LastPrice) as Low, LAST(LastPrice) as Close,
							 last(AverageTradePrice) as AverageTradePrice, mean(TotalBuyQuantity) as TotalBuyQuantity,
							 mean(TotalSellQuantity) as TotalSellQuantity,
							 last(VolumeTraded) as VolumeTraded,
							 last(TotalBuy) as TotalBuy,
							 last(TotalSell) as TotalSell,
							 INTO %s FROM %s GROUP BY time(%s)` */
