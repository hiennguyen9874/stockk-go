package crawlers

import (
	"time"

	"github.com/hiennguyen9874/stockk-go/config"
)

type Resolution int64

const (
	R1 Resolution = iota
	R5
	R15
	R30
	R60
	RD
	RW
	RM
)

type Crawler interface {
	VNDCrawlStockSymbols() ([]Ticker, error)
	VNDCrawlStockHistory(symbol string, resolution Resolution, from int64, to int64) ([]Bar, error)
	VNDMapResolutionToString(resolution Resolution) (string, error)
}

type crawler struct {
	cfg *config.Config
}

type Ticker struct {
	Symbol    string
	Exchange  string
	FullName  string
	ShortName string
	Type      string
}

type Bar struct {
	Time   time.Time
	Open   float32
	High   float32
	Low    float32
	Close  float32
	Volume float64
}

func NewCrawler(cfg *config.Config) Crawler {
	return &crawler{cfg: cfg}
}
