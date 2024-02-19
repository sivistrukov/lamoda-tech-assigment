package v1

import (
	"errors"
	"lamoda-tech-assigment/internal/entrypoints/http/shared"
	"net/http"
)

// AddNewProduct godoc
//
//	@Summary		Add new product
//	@Description	add new product in database
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			product	body		Product	true "New product"
//	@Success		201		{object}	shared.ResultResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Router			/v1/products/ [post]
func (h *Handler) AddNewProduct(w http.ResponseWriter, r *http.Request) {
	var payload Product
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

	_, err = h.uc.AddProduct(payload.Code, payload.Name, payload.Size)
	if err != nil {
		shared.WriteJSON(
			http.StatusInternalServerError,
			&shared.ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Details:    "Internal server error",
			}, w,
		)
		return
	}

	shared.WriteJSON(http.StatusCreated, &shared.ResultResponse{Ok: true}, w)
}
