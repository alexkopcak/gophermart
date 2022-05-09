package models

import (
	"github.com/jackc/pgtype"
)

type Withdrawals struct {
	OrderID     string           `json:"order"`
	Sum         float32          `json:"sum"`
	ProcessedAt pgtype.Timestamp `json:"processed_at"`
}
