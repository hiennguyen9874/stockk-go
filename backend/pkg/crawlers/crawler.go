package crawlers

import (
	"context"
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
	VNDCrawlStockSymbols(ctx context.Context) ([]Ticker, error)
	VNDCrawlStockHistory(ctx context.Context, symbol string, resolution Resolution, from int64, to int64) ([]Bar, error)
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
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

func NewCrawler(cfg *config.Config) Crawler {
	return &crawler{cfg: cfg}
}
