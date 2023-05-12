package models

type StockSnapshot struct {
	Ticker          string
	BasicPrice      float32
	CeilingPrice    float32
	FloorPrice      float32
	AccumulatedVol  float32
	AccumulatedVal  float32
	MatchPrice      float32
	MatchQtty       float32
	HighestPrice    float32
	LowestPrice     float32
	BuyForeignQtty  float32
	SellForeignQtty float32
	ProjectOpen     float32
	CurrentRoom     float32
	FloorCode       string
	TotalRoom       float32
}
