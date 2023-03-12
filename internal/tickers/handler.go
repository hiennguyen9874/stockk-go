package tickers

import "net/http"

type Handlers interface {
	Get() func(w http.ResponseWriter, r *http.Request)
	GetMulti() func(w http.ResponseWriter, r *http.Request)
	Delete() func(w http.ResponseWriter, r *http.Request)
}
