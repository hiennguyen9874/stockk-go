package presenter

import "time"

type ClientCreate struct {
	CurrentTicker     *string `json:"current_ticker" example:"TCB"`
	CurrentResolution *string `json:"current_resolution" example:"D"`
}

type ClientResponse struct {
	Id                uint      `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	CurrentTicker     *string   `json:"current_ticker,omitempty" example:"TCB"`
	CurrentResolution *string   `json:"current_resolution,omitempty" example:"D"`
	OwnerId           uint      `json:"owner_id"`
}

type ClientUpdate struct {
	CurrentTicker     *string `json:"current_ticker" example:"TCB"`
	CurrentResolution *string `json:"current_resolution" example:"D"`
}
