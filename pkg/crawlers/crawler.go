package crawlers

import (
	"github.com/hiennguyen9874/stockk-go/config"
)

type Crawler interface {
	VNDCrawlStockSymbols() ([]Ticker, error)
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

func NewCrawler(cfg *config.Config) Crawler {
	return &crawler{cfg: cfg}
}
