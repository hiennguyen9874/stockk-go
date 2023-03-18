package presenter

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hiennguyen9874/stockk-go/pkg/httpErrors"
)

type DchartExchange struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Desc  string `json:"desc"`
}

type DchartUnit struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DchartSymbolType struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DchartConfig struct {
	Exchanges              *[]DchartExchange   `json:"exchanges,omitempty"`
	SupportedResolutions   *[]string           `json:"supported_resolutions,omitempty"`
	Units                  *[]DchartUnit       `json:"units,omitempty"`
	CurrencyCodes          *[]string           `json:"currency_codes,omitempty"`
	SupportsMarks          *bool               `json:"supports_marks,omitempty"`
	SupportsTime           *bool               `json:"supports_time,omitempty"`
	SupportsTimescaleMarks *bool               `json:"supports_timescale_marks,omitempty"`
	SymbolsTypes           *[]DchartSymbolType `json:"symbols_types,omitempty"`
	SupportsSearch         *bool               `json:"supports_search,omitempty"`
	SupportsGroupRequest   *bool               `json:"supports_group_request,omitempty"`
}

type DchartLibrarySymbolInfo struct {
	Name                 string    `json:"name"`
	FullName             string    `json:"full_name"`
	BaseName             *[]string `json:"base_name,omitempty"`
	Ticker               *string   `json:"ticker,omitempty"`
	Description          string    `json:"description"`
	Type                 string    `json:"type"`
	Session              string    `json:"session"`
	SessionDisplay       *string   `json:"session_display,omitempty"`
	Holidays             *string   `json:"holidays,omitempty"`
	Corrections          *string   `json:"corrections,omitempty"`
	Exchange             string    `json:"exchange"`
	ListedExchange       string    `json:"listed_exchange"`
	Timezone             string    `json:"timezone"`
	Format               string    `json:"format"`
	Pricescale           float32   `json:"pricescale"`
	Minmov               int32     `json:"minmov"`
	Fractional           *bool     `json:"fractional,omitempty"`
	Minmove2             int32     `json:"minmove2"`
	HasIntraday          *bool     `json:"has_intraday,omitempty"`
	SupportedResolutions []string  `json:"supported_resolutions"`
	IntradayMultipliers  *[]string `json:"intraday_multipliers,omitempty"`
	HasSeconds           *bool     `json:"has_seconds,omitempty"`
	HasTicks             *bool     `json:"has_ticks,omitempty"`
	SecondsMultipliers   *[]string `json:"seconds_multipliers,omitempty"`
	HasDaily             *bool     `json:"has_daily,omitempty"`
	HasWeeklyAndMonthly  *bool     `json:"has_weekly_and_monthly,omitempty"`
	HasEmptyBars         *bool     `json:"has_empty_bars,omitempty"`
	HasNoVolume          *bool     `json:"has_no_volume,omitempty"`
	VolumePrecision      *int32    `json:"volume_precision,omitempty"`
	DataStatus           *string   `json:"data_status,omitempty"`
	Expired              *bool     `json:"expired,omitempty"`
	ExpirationDate       *int32    `json:"expiration_date,omitempty"`
	Sector               *string   `json:"sector,omitempty"`
	Industry             *string   `json:"industry,omitempty"`
	CurrencyCode         *string   `json:"currency_code,omitempty"`
	OriginalCurrencyCode *string   `json:"original_currency_code,omitempty"`
	UnitId               *string   `json:"unit_id,omitempty"`
	OriginalUnitId       *string   `json:"original_unit_id,omitempty"`
	UnitConversionTypes  *[]string `json:"unit_conversion_types,omitempty"`
}

type DchartSearchSymbolResultItem struct {
	Symbol      string `json:"symbol"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Exchange    string `json:"exchange"`
	Ticker      string `json:"ticker"`
	Type        string `json:"type"`
}

type DchartHistoryFullDataResponse struct {
	Status string    `json:"s"`
	Time   []int64   `json:"t"`
	Close  []float64 `json:"c"`
	Open   []float64 `json:"o"`
	High   []float64 `json:"h"`
	Low    []float64 `json:"l"`
	Volume []int64   `json:"v"`
}

type DchartHistoryNoDataResponse struct {
	Status  string `json:"s"`
	NexTime *int64 `json:"nextTime,omitempty"`
}

type DchartHistoryErrorResponse struct {
	Status   string                  `json:"s"`
	ErrorMsg string                  `json:"errmsg"`
	Error    *httpErrors.ErrResponse `json:"-"`
}

func (e *DchartHistoryErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.Error.Status)
	return nil
}
