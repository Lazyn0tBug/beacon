// models/user.go
package model

import (
	"time"

	"github.com/Lazyn0tBug/beacon/server/utils"
	"gorm.io/gorm"
)

type Users struct {
	GBN_MODEL
	// global.GVA_MODEL
	Username string `json:"userName" gorm:"comment:用户登录名"`             // 用户登录名
	Password string `json:"-"  gorm:"comment:用户登录密码"`                  // 用户登录密码
	NickName string `json:"nickName" gorm:"default:系统用户;comment:用户昵称"` // 用户昵称
	RoleId   uint   `gorm:"comment:角色ID"`
	// SideMode    string         `json:"sideMode" gorm:"default:dark;comment:用户侧边主题"`                                          // 用户侧边主题
	// HeaderImg   string         `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
	// BaseColor   string         `json:"baseColor" gorm:"default:#fff;comment:基础颜色"`                                           // 基础颜色
	// ActiveColor string         `json:"activeColor" gorm:"default:#1890ff;comment:活跃颜色"`                                      // 活跃颜色
	// AuthorityId uint           `json:"authorityId" gorm:"default:888;comment:用户角色ID"`                                        // 用户角色ID
	// Authority   SysAuthority   `json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:用户角色"`
	// Authorities []SysAuthority `json:"authorities" gorm:"many2many:sys_user_authority;"`
	// Phone       string         `json:"phone"  gorm:"comment:用户手机号"`                     // 用户手机号
	Email    string `json:"email"  gorm:"comment:用户邮箱"`                      // 用户邮箱
	IsActive int    `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 0冻结"` //用户是否被冻结 1正常 2冻结
}

func (Users) TableName() string {
	return "Users"
}

func (u *Users) BeforeCreate(tx *gorm.DB) (err error) {
	return u.HashPassword()
}

// HashPassword hashes the password using bcrypt.
func (u *Users) HashPassword() error {
	passwordHash, err := utils.BcryptHash(string(u.Password))
	if err != nil {
		return err
	}
	u.Password = string(passwordHash)
	return nil
}

// CheckPassword checks if the provided password matches the stored password.
func (u *Users) CheckPassword(password string) bool {
	return utils.BcryptCheck(u.Password, password)
}

// SetActive sets the user as active.
func (u *Users) SetActive() {
	u.IsActive = 1
	u.UpdatedAt = time.Now()
}
