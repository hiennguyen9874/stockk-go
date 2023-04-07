package tickers

import "net/http"

type Handlers interface {
	GetMulti() func(w http.ResponseWriter, r *http.Request)
	GetBySymbol() func(w http.ResponseWriter, r *http.Request)
	UpdateIsActiveBySymbol() func(w http.ResponseWriter, r *http.Request)
	SearchBySymbol() func(w http.ResponseWriter, r *http.Request)
}
