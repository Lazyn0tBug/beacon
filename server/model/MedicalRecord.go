package model

import "github.com/Lazyn0tBug/beacon/server/global"

// MedicalRecord 模型表示病历信息。
type MedicalRecord struct {
	global.GBN_MODEL
	UserID     uint   `gorm:"comment:用户ID;column:user_id" json:"user_id"`
	DoctorID   uint   `gorm:"comment:医生ID;column:doctor_id" json:"doctor_id"`
	RecordTime string `gorm:"comment:病历记录时间;column:record_time" json:"record_time"`
	Diagnosis  string `gorm:"comment:诊断结果;column:diagnosis" json:"diagnosis"`
	Treatment  string `gorm:"comment:治疗方法;column:treatment" json:"treatment"`
	Medicine   string `gorm:"comment:药品用量;column:medicine" json:"medicine"`
	Checkup    string `gorm:"comment:检查结果;column:checkup" json:"checkup"`
	Advice     string `gorm:"comment:医嘱;column:advice" json:"advice"`
}

func (MedicalRecord) TableName() string {
	return "MedicalRecord"
}
