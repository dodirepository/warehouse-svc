package router

import (
	"net/http"

	lib "github.com/dodirepository/common-lib"
	"github.com/dodirepository/warehouse-svc/infrastructure/bootstrap"
	"github.com/dodirepository/warehouse-svc/internal/adapters/controllers/handlers"
	"github.com/gorilla/mux"
)

type router struct {
	router   *mux.Router
	services *bootstrap.Service
}

func NewRouter(services *bootstrap.Service) Router {
	return &router{
		router:   mux.NewRouter(),
		services: services,
	}
}

// Route :nodoc:
func (rtr *router) Route() *mux.Router {
	root := rtr.router.PathPrefix("/").Subrouter()
	in := root.PathPrefix("/in").Subrouter()
	v1 := in.PathPrefix("/v1").Subrouter()

	in.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		lib.Render("ok", http.StatusOK, w)
	}).Methods(http.MethodGet)

	whPutHandlers := handlers.NewUpdateStatus(rtr.services.Warehouse)
	whCreateHandlers := handlers.NewCreate(rtr.services.Warehouse)
	addDetails := handlers.NewAddProduct(rtr.services.Warehouse)
	tfHandlers := handlers.NewTransferProduct(rtr.services.Warehouse)
	wh := v1.PathPrefix("/warehouse").Subrouter()
	wh.HandleFunc("", whCreateHandlers.Handlers).Methods(http.MethodPost)
	wh.HandleFunc("/{id}/product", addDetails.Handlers).Methods(http.MethodPost)
	wh.HandleFunc("/status", whPutHandlers.Handlers).Methods(http.MethodPut)
	wh.HandleFunc("/{id}/product/move", tfHandlers.Handlers).Methods(http.MethodPost)
	return rtr.router
}
