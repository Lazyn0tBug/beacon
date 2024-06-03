package model

import "github.com/Lazyn0tBug/beacon/server/global"

type Customer struct {
	global.GBN_MODEL
	Name       string     `gorm:"comment:用户名称;column:name" json:"name"`
	Gender     string     `gorm:"comment:用户性别;column:gender" json:"gender"`
	Age        uint       `gorm:"comment:用户年龄;column:age" json:"age"`
	Phone      string     `gorm:"comment:用户手机号;column:phone" json:"phone"`
	Email      string     `gorm:"comment:用户邮箱;column:email" json:"email"`
	Address    string     `gorm:"comment:用户住址;column:address" json:"address"`
	IDCard     string     `gorm:"comment:用户身份证号;column:id_card" json:"id_card"`
	Avatar     string     `gorm:"comment:用户头像;column:avatar" json:"avatar"`
	MemberType MemberType `gorm:"comment:会员类型;column:member_type" json:"member_type"`
	Point      int        `gorm:"comment:积分;column:point" json:"point"`
}

type MemberLevel struct {
	global.GBN_MODEL
	Name        string `gorm:"comment:会员等级名称;column:name" json:"name"`
	Description string `gorm:"comment:会员等级描述;column:description" json:"description"`
}

type MemberType string

const (
	NonMember    MemberType = "非会员"
	Member       MemberType = "会员"
	SeniorMember MemberType = "高级会员"
)

func (Customer) TableName() string {
	return "Customer"
}
