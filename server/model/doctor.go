package model

import "github.com/Lazyn0tBug/beacon/server/global"

// Doctor 模型表示医生信息。
type Doctor struct {
	global.GBN_MODEL
	Name         string `gorm:"comment:医生名称;column:name" json:"name"`
	Gender       string `gorm:"comment:医生性别;column:gender" json:"gender"`
	Age          uint   `gorm:"comment:医生年龄;column:age" json:"age"`
	Phone        string `gorm:"comment:医生手机号;column:phone" json:"phone"`
	Email        string `gorm:"comment:医生邮箱;column:email" json:"email"`
	Title        string `gorm:"comment:医生职称;column:title" json:"title"`
	Specialty    string `gorm:"comment:医生擅长;column:specialty" json:"specialty"`
	Introduction string `gorm:"comment:医生简介;column:introduction" json:"introduction"`
	Avatar       string `gorm:"comment:医生头像;column:avatar" json:"avatar"`
	HospitalID   uint   `gorm:"comment:所属医院ID;column:hospital_id" json:"hospital_id"`
}

func (Doctor) TableName() string {
	return "Doctor"
}
