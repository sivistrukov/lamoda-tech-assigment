package v1

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"lamoda-tech-assigment/internal/services/usecases"
	"net/http"
)

type Handler struct {
	uc usecases.IUseCases
}

// NewHandler returns new instance of Handler
func NewHandler(uc usecases.IUseCases) Handler {
	return Handler{uc: uc}
}

// TODO: remove
func (h *Handler) Live(w http.ResponseWriter, r *http.Request) {
	// an example API handler
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

// InitializeRoutes is initializing handler routes.
func (h *Handler) InitializeRoutes(router *mux.Router) {
	router.HandleFunc("/live/", h.Live)

	products := router.PathPrefix("/products").Subrouter()
	products.HandleFunc("/", h.AddNewProduct).Methods(http.MethodPost)

	warehouses := router.PathPrefix("/warehouses").Subrouter()
	warehouses.HandleFunc("/", h.AddWarehouse).Methods(http.MethodPost)
	warehouses.HandleFunc("/{id:[0-9]+}/products/", h.GetWarehouseProducts).Methods(http.MethodGet)
	warehouses.HandleFunc("/{id:[0-9]+}/products/", h.AddProductsToWarehouse).Methods(http.MethodPost)
	warehouses.HandleFunc("/{id:[0-9]+}/products/{code}/quantity/", h.GetProductQuantity).Methods(http.MethodGet)
	warehouses.HandleFunc("/{id:[0-9]+}/reserve/", h.ReserveProductFromWarehouse).Methods(http.MethodPost)
	warehouses.HandleFunc("/{id:[0-9]+}/cancel-reservation/", h.CancelReservationFromWarehouse).Methods(http.MethodPost)
}
