package models

// Approve Public
type Approve struct {
	ID            int    `json:"id" gorm:"column:id"`
	MerchantID    int    `json:"merchant_id" gorm:"column:merchant_id"`
	CCNumber      string `json:"cc_number" gorm:"column:cc_number"`
	CCCVV         string `json:"cc_cvv" gorm:"column:cc_cvv"`
	CCMonth       string `json:"cc_month" gorm:"column:cc_month"`
	CCYear        string `json:"cc_year" gorm:"column:cc_year"`
	Currency      string `json:"currency" gorm:"column:currency"`
	Amount        string `json:"amount" gorm:"column:amount"`
	ApproveKey    string `json:"approve_key" gorm:"column:approve_key"`
	AmountBalance string `json:"amount_balance" gorm:"column:amount_balance"`
	Status        int    `json:"status" gorm:"column:status"`
	CreatedAt     string `json:"created_at" gorm:"column:created_at"`
}

// TableName Public
func (Approve) TableName() string {
	return "approve"
}
