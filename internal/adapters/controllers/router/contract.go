package router

import (
	"github.com/gorilla/mux"
)

// Router :nodoc:
type Router interface {
	Route() *mux.Router
}
