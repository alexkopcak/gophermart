package models

import "github.com/jackc/pgtype"

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
	OrderStatusWithDrawn  = "WITHDRAWN"
)

type Order struct {
	UserName int32            `json:"-"`
	Number   string           `json:"number"`
	Status   string           `json:"status"`
	Accrual  float32          `json:"accrual,omitempty"`
	Uploaded pgtype.Timestamp `json:"uploaded_at"`
}
