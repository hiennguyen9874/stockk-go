package dchart

import "net/http"

type Handlers interface {
	GetTime() func(w http.ResponseWriter, r *http.Request)
	GetConfig() func(w http.ResponseWriter, r *http.Request)
	GetSymbols() func(w http.ResponseWriter, r *http.Request)
	Search() func(w http.ResponseWriter, r *http.Request)
	History() func(w http.ResponseWriter, r *http.Request)
}
