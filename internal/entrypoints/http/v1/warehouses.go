package v1

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"lamoda-tech-assigment/internal/adapters/postgresql"
	"lamoda-tech-assigment/internal/domain"
	"lamoda-tech-assigment/internal/entrypoints/http/shared"
	"net/http"
	"strconv"
)

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

func (h *Handler) AddProductsToWarehouse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	warehouseID, _ := strconv.ParseUint(params["id"], 10, 64)

	var payload []AddProduct
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

	shared.WriteJSON(http.StatusOK, nil, w)
}

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

	shared.WriteJSON(http.StatusOK, nil, w)
}

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

	response := map[string]uint{"quantity": quantity}
	shared.WriteJSON(http.StatusOK, &response, w)
}
