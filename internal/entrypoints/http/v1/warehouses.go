package v1

import (
	"errors"
	"fmt"
	"lamoda-tech-assigment/internal/adapters/postgresql"
	"lamoda-tech-assigment/internal/domain"
	"lamoda-tech-assigment/internal/entrypoints/http/shared"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// AddWarehouse godoc
//
//	@Summary		Add new warehouse
//	@Description	add new warehouse in database
//	@Tags			warehouses
//	@Accept			json
//	@Produce		json
//	@Param			warehouse	body		CreateWarehouse	true	"New warehouse"
//	@Success		201			{object}	Warehouse
//	@Failure		400			{object}	shared.ErrorResponse
//	@Router			/v1/warehouses/ [post]
func (h *Handler) AddWarehouse(w http.ResponseWriter, r *http.Request) {
	var payload CreateWarehouse
	err := shared.DecodeJSON(r, &payload)
	if err != nil {
		var errResp *shared.ErrorResponse
		if errors.As(err, &errResp) {
			shared.WriteJSON(errResp.StatusCode, errResp, w)
			return
		}
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	warehouse, err := h.uc.AddWarehouse(payload.Name)
	if err != nil {
		shared.WriteJSON(
			http.StatusInternalServerError,
			&shared.ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Details:    "Internal server error",
			},
			w,
		)
		return
	}

	response := Warehouse{
		ID:          warehouse.ID,
		Name:        warehouse.Name,
		IsAvailable: warehouse.IsAvailable,
	}

	shared.WriteJSON(http.StatusCreated, &response, w)
}

// AddProductsToWarehouse godoc
//
//	@Summary		Add products
//	@Description	add products to warehouse
//	@Tags			warehouses
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int					true	"Warehouse ID"
//	@Param			products	body		[]ProductQuantity	true	"Products"
//	@Success		201			{array}		ProductQuantity
//	@Failure		400			{object}	shared.ErrorResponse
//	@Failure		404			{object}	shared.ErrorResponse
//	@Router			/v1/warehouses/{id}/products/ [post]
func (h *Handler) AddProductsToWarehouse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	warehouseID, _ := strconv.ParseUint(params["id"], 10, 64)

	var payload []ProductQuantity
	err := shared.DecodeJSON(r, &payload)
	if err != nil {
		var errResp *shared.ErrorResponse
		if errors.As(err, &errResp) {
			shared.WriteJSON(errResp.StatusCode, errResp, w)
			return
		}
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	if len(payload) < 1 {
		errResp := shared.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Details:    "payload is empty",
		}
		shared.WriteJSON(errResp.StatusCode, errResp, w)
		return
	}

	products := make([]domain.Product, len(payload))
	for i, p := range payload {
		products[i] = domain.Product{
			Code:     p.Code,
			Quantity: p.Quantity,
		}
	}

	result, err := h.uc.AddProductToWarehouse(uint(warehouseID), products...)
	if err != nil {
		if errors.Is(err, postgresql.EntityNotFoundError) {
			shared.WriteJSON(
				http.StatusNotFound,
				&shared.ErrorResponse{
					StatusCode: http.StatusNotFound,
					Details:    fmt.Sprintf("Warehouse with id %v not found", warehouseID),
				}, w,
			)
			return
		}
		if errors.Is(err, postgresql.ForeignKeyError) {
			shared.WriteJSON(
				http.StatusNotFound,
				&shared.ErrorResponse{
					StatusCode: http.StatusBadRequest,
					Details:    "the list contains non-existent product code",
				}, w,
			)
			return
		}
		shared.WriteJSON(
			http.StatusInternalServerError,
			&shared.ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Details:    "Internal server error",
			}, w,
		)
		return
	}

	shared.WriteJSON(http.StatusOK, &result, w)
}

// ReserveProductFromWarehouse godoc
//
//	@Summary		Reserve products
//	@Description	reserve products in warehouse
//	@Tags			warehouses
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int						true	"Warehouse ID"
//	@Param			products	body		[]ReservationRequest	true	"Products"
//	@Success		200			{object}	shared.ResultResponse
//	@Failure		400			{object}	shared.ErrorResponse
//	@Failure		404			{object}	shared.ErrorResponse
//	@Router			/v1/warehouses/{id}/reserve/ [post]
func (h *Handler) ReserveProductFromWarehouse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	warehouseID, _ := strconv.ParseUint(params["id"], 10, 64)

	var payload []ReservationRequest
	err := shared.DecodeJSON(r, &payload)
	if err != nil {
		var errResp *shared.ErrorResponse
		if errors.As(err, &errResp) {
			shared.WriteJSON(errResp.StatusCode, errResp, w)
			return
		}
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	products := make([]domain.Product, len(payload))
	for i, command := range payload {
		products[i] = domain.Product{Code: command.Code, Quantity: command.Quantity}
	}

	err = h.uc.ReserveProductsInWarehouse(uint(warehouseID), products...)
	if err != nil {
		if errors.Is(err, domain.ProductNotFoundError) {
			shared.WriteJSON(
				http.StatusNotFound,
				&shared.ErrorResponse{
					StatusCode: http.StatusBadRequest,
					Details:    "the list contains non-existent product code",
				}, w,
			)
			return
		}

		var quantityErr *domain.NotEnoughProductQuantityError
		if errors.As(err, &quantityErr) {
			shared.WriteJSON(
				http.StatusBadRequest,
				&shared.ErrorResponse{
					StatusCode: http.StatusBadRequest,
					Details:    err.Error(),
				}, w,
			)
			return
		}

		shared.WriteJSON(
			http.StatusInternalServerError,
			&shared.ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Details:    "internal server error",
			}, w,
		)
		return
	}

	shared.WriteJSON(http.StatusOK, &shared.ResultResponse{Ok: true}, w)
}

// CancelReservationFromWarehouse godoc
//
//	@Summary		Cancel reservation products
//	@Description	cancel reservation reserve products in warehouse
//	@Tags			warehouses
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int						true	"Warehouse ID"
//	@Param			products	body		[]ReservationRequest	true	"Products"
//	@Success		200			{object}	shared.ResultResponse
//	@Failure		400			{object}	shared.ErrorResponse
//	@Failure		404			{object}	shared.ErrorResponse
//	@Router			/v1/warehouses/{id}/cancel-reservation/ [post]
func (h *Handler) CancelReservationFromWarehouse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	warehouseID, _ := strconv.ParseUint(params["id"], 10, 64)

	var payload []ReservationRequest
	err := shared.DecodeJSON(r, &payload)
	if err != nil {
		var errResp *shared.ErrorResponse
		if errors.As(err, &errResp) {
			shared.WriteJSON(errResp.StatusCode, errResp, w)
			return
		}
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	products := make([]domain.Product, len(payload))
	for i, command := range payload {
		products[i] = domain.Product{Code: command.Code, Quantity: command.Quantity}
	}

	err = h.uc.CancelReservationInWarehouse(uint(warehouseID), products...)
	if err != nil {
		if errors.Is(err, domain.ProductNotFoundError) {
			shared.WriteJSON(
				http.StatusNotFound,
				&shared.ErrorResponse{
					StatusCode: http.StatusBadRequest,
					Details:    "the list contains unreserved items",
				}, w,
			)
			return
		}

		var quantityErr *domain.NotEnoughProductQuantityError
		if errors.As(err, &quantityErr) {
			shared.WriteJSON(
				http.StatusBadRequest,
				&shared.ErrorResponse{
					StatusCode: http.StatusBadRequest,
					Details:    err.Error(),
				}, w,
			)
			return
		}

		shared.WriteJSON(
			http.StatusInternalServerError,
			&shared.ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Details:    "internal server error",
			}, w,
		)
		return
	}

	shared.WriteJSON(http.StatusOK, &shared.ResultResponse{Ok: true}, w)
}

// GetWarehouseProducts godoc
//
//	@Summary		Get products
//	@Description	get products stored in warehouse
//	@Tags			warehouses
//	@Produce		json
//	@Param			id	path		int	true	"Warehouse ID"
//	@Success		200	{array}		domain.Product
//	@Failure		404	{object}	shared.ErrorResponse
//	@Router			/v1/warehouses/{id}/products/ [get]
func (h *Handler) GetWarehouseProducts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	warehouseID, _ := strconv.ParseUint(params["id"], 10, 64)

	products, err := h.uc.GetWarehouseProducts(uint(warehouseID))
	if err != nil {
		if errors.Is(err, postgresql.EntityNotFoundError) {
			shared.WriteJSON(
				http.StatusNotFound,
				&shared.ErrorResponse{
					StatusCode: http.StatusNotFound,
					Details:    fmt.Sprintf("Warehouse with id %v not found", warehouseID),
				},
				w,
			)
			return
		}
		shared.WriteJSON(
			http.StatusInternalServerError,
			&shared.ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Details:    "Internal server error",
			},
			w,
		)
		return
	}

	shared.WriteJSON(http.StatusOK, &products, w)
}

// GetProductQuantity godoc
//
//	@Summary		Get product's quantity
//	@Description	get product's quantity stored in warehouse
//	@Tags			warehouses
//	@Produce		json
//	@Param			id		path		int	true	"Warehouse ID"
//	@Param			code	path		int	true	"Product code"
//	@Success		200		{array}		ProductQuantity
//	@Failure		404		{object}	shared.ErrorResponse
//	@Router			/v1/warehouses/{id}/products/{code}/quantity [get]
func (h *Handler) GetProductQuantity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	warehouseID, _ := strconv.ParseUint(params["id"], 10, 64)
	code := params["code"]
	quantity, err := h.uc.GetProductQuantityInWarehouse(uint(warehouseID), code)
	if err != nil {
		if errors.Is(err, domain.ProductNotFoundError) {
			shared.WriteJSON(
				http.StatusNotFound,
				&shared.ErrorResponse{
					StatusCode: http.StatusNotFound,
					Details:    fmt.Sprintf("Product with code %s not found", code),
				},
				w,
			)
			return
		}
		if errors.Is(err, postgresql.EntityNotFoundError) {
			shared.WriteJSON(
				http.StatusNotFound,
				&shared.ErrorResponse{
					StatusCode: http.StatusNotFound,
					Details:    fmt.Sprintf("Warehouse with id %v not found", warehouseID),
				},
				w,
			)
			return
		}
		shared.WriteJSON(
			http.StatusInternalServerError,
			&shared.ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Details:    "Internal server error",
			},
			w,
		)
		return
	}

	response := map[string]any{"code": code, "quantity": quantity}
	shared.WriteJSON(http.StatusOK, &response, w)
}
