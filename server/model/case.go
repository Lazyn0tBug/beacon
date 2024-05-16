// models/case.go
package model

type Case struct {
	GBN_MODEL
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (Case) TableName() string {
	return "Case"
}
