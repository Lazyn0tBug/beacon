// models/case.go
package model

type Cases struct {
	GBN_MODEL
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (Cases) TableName() string {
	return "Cases"
}
