package models

// Ledger Public
type Ledger struct {
	ID         int    `json:"id" gorm:"column:id"`
	MerchantID int    `json:"merchant_id" gorm:"column:merchant_id"`
	ApproveID  int    `json:"approve_id" gorm:"column:approve_id"`
	Amount     string `json:"amount" gorm:"column:amount"`
	ActionType string `json:"action_type" gorm:"column:action_type"`
	CreatedAt  string `json:"created_at" gorm:"column:created_at"`
}

// TableName Public
func (Ledger) TableName() string {
	return "ledger"
}
