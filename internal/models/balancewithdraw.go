package models

type BalanceWithdraw struct {
	OrderID string  `json:"order"`
	Sum     float32 `json:"sum"`
}
