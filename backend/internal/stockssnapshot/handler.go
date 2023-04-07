package stockssnapshot

import "net/http"

type Handlers interface {
	GetStockSnapshotBySymbol() func(w http.ResponseWriter, r *http.Request)
}
