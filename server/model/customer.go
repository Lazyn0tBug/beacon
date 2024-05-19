package model

import "github.com/Lazyn0tBug/beacon/server/global"

type Customer struct {
	global.GBN_MODEL
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

func (Customer) TableName() string {
	return "Customer"
}
