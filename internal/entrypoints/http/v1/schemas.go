package v1

type Warehouse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	IsAvailable bool   `json:"isAvailable"`
}

type CreateWarehouse struct {
	Name string `json:"name"`
}

type Product struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Size string `json:"size"`
}

type AddProduct struct {
	Code     string `json:"code"`
	Quantity uint   `json:"quantity"`
}

type ReservationRequest struct {
	Code     string `json:"code"`
	Quantity uint   `json:"quantity"`
}
