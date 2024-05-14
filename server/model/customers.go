package model


type Customers struct {
	GBN_MODEL
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

func (Customers) TableName() string {
	return "Customers"
}
