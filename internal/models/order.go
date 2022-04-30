package models

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
	OrderStatusWithDrawn  = "WITHDRAWN"
)

type Order struct {
	UserName string  `json:"-"`
	Number   string  `json:"number"`
	Status   string  `json:"status"`
	Accrual  float32 `json:"accurual,omitempty"`
	Uploaded string  `json:"uploaded_at"`
}
