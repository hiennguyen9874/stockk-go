package presenter

type StockSnapshotResponse struct {
	Ticker      string  `json:"ticker"`
	RefPrice    float32 `json:"ref_price"`
	CeilPrice   float32 `json:"ceil_price"`
	FloorPrice  float32 `json:"floor_price"`
	TltVol      float32 `json:"tlt_vol"`
	TltVal      float32 `json:"tlt_val"`
	PriceB3     float32 `json:"price_b3"`
	PriceB2     float32 `json:"price_b2"`
	PriceB1     float32 `json:"price_b1"`
	VolB3       float32 `json:"vol_b3"`
	VolB2       float32 `json:"vol_b2"`
	VolB1       float32 `json:"vol_b1"`
	Price       float32 `json:"price"`
	Vol         float32 `json:"vol"`
	PriceS3     float32 `json:"price_s3"`
	PriceS2     float32 `json:"price_s2"`
	PriceS1     float32 `json:"price_s1"`
	VolS3       float32 `json:"vol_s3"`
	VolS2       float32 `json:"vol_s2"`
	VolS1       float32 `json:"vol_s1"`
	High        float32 `json:"high"`
	Low         float32 `json:"low"`
	BuyForeign  float32 `json:"buy_foreign"`
	SellForeign float32 `json:"sell_foreign"`
}
