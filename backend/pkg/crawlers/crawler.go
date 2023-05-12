package crawlers

import (
	"context"
	"fmt"
	"time"

	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
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
	CrawlStockSymbols(ctx context.Context) ([]Ticker, error)
	CrawlStockHistory(ctx context.Context, symbol string, resolution Resolution, from int64, to int64) ([]Bar, error)
	CrawlStockSnapshot(ctx context.Context, symbols []string) ([]StockSnapshot, error)
	VNDMapExchange(exchange string) (string, error)
	VNDCrawlStockSymbols(ctx context.Context) ([]Ticker, error)
	VNDMapResolutionToString(resolution Resolution) (string, error)
	VNDCrawlStockHistory(ctx context.Context, symbol string, resolution Resolution, from int64, to int64) ([]Bar, error)
	SSIMapExchange(exchange string) (string, error)
	SSICrawlStockSymbols(ctx context.Context) ([]Ticker, error)
	SSIMapResolutionToString(resolution Resolution) (string, error)
	SSICrawlStockHistory(ctx context.Context, symbol string, resolution Resolution, from int64, to int64) ([]Bar, error)
	VNDCrawlStockSnapshot(ctx context.Context, symbols []string) ([]StockSnapshot, error)
}

type crawler struct {
	cfg    *config.Config
	logger logger.Logger
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

type StockSnapshot struct {
	Code            string
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

func NewCrawler(cfg *config.Config, logger logger.Logger) Crawler {
	return &crawler{cfg: cfg, logger: logger}
}

func (cr *crawler) CrawlStockSymbols(ctx context.Context) ([]Ticker, error) {
	switch cr.cfg.Crawler.Source {
	case "VND":
		return cr.VNDCrawlStockSymbols(ctx)
	case "SSI":
		return cr.SSICrawlStockSymbols(ctx)
	default:
		return nil, fmt.Errorf("not support crawler source: %v", cr.cfg.Crawler.Source)
	}
}

func (cr *crawler) CrawlStockHistory(ctx context.Context, symbol string, resolution Resolution, from int64, to int64) ([]Bar, error) {
	switch cr.cfg.Crawler.Source {
	case "VND":
		return cr.VNDCrawlStockHistory(ctx, symbol, resolution, from, to)
	case "SSI":
		return cr.SSICrawlStockHistory(ctx, symbol, resolution, from, to)
	default:
		return nil, fmt.Errorf("not support crawler source: %v", cr.cfg.Crawler.Source)
	}
}

func (cr *crawler) CrawlStockSnapshot(ctx context.Context, symbols []string) ([]StockSnapshot, error) {
	switch cr.cfg.Crawler.Source {
	case "VND":
		return cr.VNDCrawlStockSnapshot(ctx, symbols)
	default:
		return nil, fmt.Errorf("not support crawler source: %v", cr.cfg.Crawler.Source)
	}
}
