// models/case.go
package model

import "github.com/Lazyn0tBug/beacon/server/global"

type Case struct {
	global.GBN_MODEL
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (Case) TableName() string {
	return "Case"
}
