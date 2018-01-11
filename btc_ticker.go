package main

// BTCTicker allows mapping ticker data retrieved from gdax
type BTCTicker struct {
	Ask     string `json:"ask"`
	Bid     string `json:"bid"`
	Price   string `json:"price"`
	Size    string `json:"size"`
	Time    string `json:"time"`
	TradeID int64  `json:"trade_id"`
	Volume  string `json:"volume"`
}
