package v1

import (
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

// InitializeRoutes is initializing handler routes.
func (h *Handler) InitializeRoutes(router *mux.Router) {
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
