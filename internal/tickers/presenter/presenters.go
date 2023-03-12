package presenter

import (
	"github.com/google/uuid"
)

type TickerResponse struct {
	Id        uuid.UUID `json:"id,omitempty"`
	Symbol    string    `json:"symbol,omitempty"`
	Exchange  string    `json:"exchange,omitempty"`
	FullName  string    `json:"full_name,omitempty"`
	ShortName string    `json:"short_name,omitempty"`
	Type      string    `json:"type,omitempty"`
	IsActive  bool      `json:"is_active"`
}
