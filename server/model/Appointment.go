package model

import "github.com/Lazyn0tBug/beacon/server/global"

// Appointment 模型表示预约信息。
type Appointment struct {
	global.GBN_MODEL
	UserID          uint   `gorm:"comment:用户ID;column:user_id" json:"user_id"`
	DoctorID        uint   `gorm:"comment:医生ID;column:doctor_id" json:"doctor_id"`
	AppointmentTime string `gorm:"comment:预约时间;column:appointment_time" json:"appointment_time"`
	Status          string `gorm:"comment:预约状态;column:status" json:"status"`
}

func (Appointment) TableName() string {
	return "Appointment"
}
