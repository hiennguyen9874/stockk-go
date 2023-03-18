package charts

import "net/http"

type Handlers interface {
	CreateOrUpdate() func(w http.ResponseWriter, r *http.Request)
	Get() func(w http.ResponseWriter, r *http.Request)
	Delete() func(w http.ResponseWriter, r *http.Request)
}
