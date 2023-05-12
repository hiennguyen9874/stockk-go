package presenter

type StockSnapshotResponse struct {
	Ticker          string  `json:"ticker"`
	BasicPrice      float32 `json:"basic_price"`
	CeilingPrice    float32 `json:"ceiling_price"`
	FloorPrice      float32 `json:"floor_price"`
	AccumulatedVol  float32 `json:"accumulated_vol"`
	AccumulatedVal  float32 `json:"accumulated_val"`
	MatchPrice      float32 `json:"match_price"`
	MatchQtty       float32 `json:"match_qtty"`
	HighestPrice    float32 `json:"highest_price"`
	LowestPrice     float32 `json:"lowest_price"`
	BuyForeignQtty  float32 `json:"buy_foreign_qtty"`
	SellForeignQtty float32 `json:"sell_foreign_qtty"`
	ProjectOpen     float32 `json:"project_open"`
	CurrentRoom     float32 `json:"current_room"`
	FloorCode       string  `json:"floor_code"`
	TotalRoom       float32 `json:"total_room"`
}
