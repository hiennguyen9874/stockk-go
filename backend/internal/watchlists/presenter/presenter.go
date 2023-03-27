package presenter

import "time"

type WatchListCreate struct {
	Name    string   `json:"name" example:"Ngân hàng"`
	Tickers []string `json:"tickers" example:"VCB,TCB,MBB,ACB,CTG,BID"`
}

type WatchListResponse struct {
	Id        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name" example:"Ngân hàng"`
	Tickers   []string  `json:"tickers" example:"VCB,TCB,MBB,ACB,CTG,BID"`
	OwnerId   uint      `json:"owner_id"`
}

type WatchListUpdate struct {
	Name    *string   `json:"name" example:"Ngân hàng"`
	Tickers *[]string `json:"tickers" example:"VCB,TCB,MBB,ACB,CTG,BID"`
}
