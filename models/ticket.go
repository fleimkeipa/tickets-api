package models

type Ticket struct {
	ID          int64  `json:"id" pg:",pk"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	Allocation  int    `json:"allocation"`
}

type CreateRequest struct {
	Name        string `json:"name" validate:"required,min=5,max=100"`
	Description string `json:"desc" validate:"max=500"`
	Allocation  int    `json:"allocation" validate:"required,gt=0"`
}

type PurchaseRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,gt=0"`
}
