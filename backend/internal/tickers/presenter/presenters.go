package presenter

type TickerResponse struct {
	Id        uint   `json:"id,omitempty"`
	Symbol    string `json:"symbol,omitempty"`
	Exchange  string `json:"exchange,omitempty"`
	FullName  string `json:"full_name,omitempty"`
	ShortName string `json:"short_name,omitempty"`
	Type      string `json:"type,omitempty"`
	IsActive  bool   `json:"is_active"`
}
