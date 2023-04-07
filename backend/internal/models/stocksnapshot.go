package models

type StockSnapshot struct {
	Ticker      string
	RefPrice    float32
	CeilPrice   float32
	FloorPrice  float32
	TltVol      float32
	TltVal      float32
	PriceB3     float32
	PriceB2     float32
	PriceB1     float32
	VolB3       float32
	VolB2       float32
	VolB1       float32
	Price       float32
	Vol         float32
	PriceS3     float32
	PriceS2     float32
	PriceS1     float32
	VolS3       float32
	VolS2       float32
	VolS1       float32
	High        float32
	Low         float32
	BuyForeign  float32
	SellForeign float32
}
