package bidder

import (
	"bidding/api"
	appstore "bidding/store"
	"net/http"

	"github.com/go-chi/chi"
)

// store holds shared store conn from the api
var store *appstore.Conn

// Init initializes all the v1 routes
func Init(r chi.Router) {
	// FIXME: not ideal way to manage store, still its a workaround
	store = api.Store
	r.Method(http.MethodPost, "/register", api.Handler(bidderRegisterHandler))
	r.Method(http.MethodGet, "/", api.Handler(getBiddersHandler))
}
