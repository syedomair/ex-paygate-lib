package models

// Merchant Public
type Merchant struct {
	ID   int    `json:"id" gorm:"column:id"`
	Key  string `json:"key" gorm:"column:key"`
	Name string `json:"name" gorm:"column:name"`
}

// TableName Public
func (Merchant) TableName() string {
	return "merchant"
}
