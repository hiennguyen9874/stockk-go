package presenter

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
	Exchanges              *[]DchartExchange   `json:"exchanges"`
	SupportedResolutions   *[]string           `json:"supported_resolutions"`
	Units                  *[]DchartUnit       `json:"units"`
	CurrencyCodes          *[]string           `json:"currency_codes"`
	SupportsMarks          *bool               `json:"supports_marks"`
	SupportsTime           *bool               `json:"supports_time"`
	SupportsTimescaleMarks *bool               `json:"supports_timescale_marks"`
	SymbolsTypes           *[]DchartSymbolType `json:"symbols_types"`
	SupportsSearch         *bool               `json:"supports_search"`
	SupportsGroupRequest   *bool               `json:"supports_group_request"`
}

type DchartLibrarySymbolInfo struct {
	Name                 string    `json:"name"`
	FullName             string    `json:"full_name"`
	BaseName             *[]string `json:"base_name"`
	Ticker               *string   `json:"ticker"`
	Description          string    `json:"description"`
	Type                 string    `json:"type"`
	Session              string    `json:"session"`
	SessionDisplay       *string   `json:"session_display"`
	Holidays             *string   `json:"holidays"`
	Corrections          *string   `json:"corrections"`
	Exchange             string    `json:"exchange"`
	ListedExchange       string    `json:"listed_exchange"`
	Timezone             string    `json:"timezone"`
	Format               string    `json:"format"`
	Pricescale           float32   `json:"pricescale"`
	Minmov               int32     `json:"minmov"`
	Fractional           *bool     `json:"fractional"`
	Minmove2             int32     `json:"minmove2"`
	HasIntraday          *bool     `json:"has_intraday"`
	SupportedResolutions []string  `json:"supported_resolutions"`
	IntradayMultipliers  *[]string `json:"intraday_multipliers"`
	HasSeconds           *bool     `json:"has_seconds"`
	HasTicks             *bool     `json:"has_ticks"`
	SecondsMultipliers   *[]string `json:"seconds_multipliers"`
	HasDaily             *bool     `json:"has_daily"`
	HasWeeklyAndMonthly  *bool     `json:"has_weekly_and_monthly"`
	HasEmptyBars         *bool     `json:"has_empty_bars"`
	HasNoVolume          *bool     `json:"has_no_volume"`
	VolumePrecision      *int32    `json:"volume_precision"`
	DataStatus           *string   `json:"data_status"`
	Expired              *bool     `json:"expired"`
	ExpirationDate       *int32    `json:"expiration_date"`
	Sector               *string   `json:"sector"`
	Industry             *string   `json:"industry"`
	CurrencyCode         *string   `json:"currency_code"`
	OriginalCurrencyCode *string   `json:"original_currency_code"`
	UnitId               *string   `json:"unit_id"`
	OriginalUnitId       *string   `json:"original_unit_id"`
	UnitConversionTypes  *[]string `json:"unit_conversion_types"`
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
