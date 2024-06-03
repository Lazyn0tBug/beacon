// models/user.go
package model

import (
	"time"

	"github.com/Lazyn0tBug/beacon/server/global"
	"github.com/Lazyn0tBug/beacon/server/utils"
	"gorm.io/gorm"
)

type User struct {
	global.GBN_MODEL
	// global.GVA_MODEL
	Username string `gorm:"column:username;size:255;not null;comment:用户登录名" "json:"userName" gorm:"comment:用户登录名"` // 用户登录名
	Password string `gorm:"column:password;size:32;not null;comment:The hashed password of the user" json:"-"`
	Nickname string `gorm:"column:nickname;size:255;comment:用户昵称" json:"nickName" json:"nickName"` // 用户昵称
	RoleID   uint   `gorm:"column:role_id;not null;comment:角色ID" json:"roleID"`                    // 角色ID
	// SideMode    string         `json:"sideMode" gorm:"default:dark;comment:用户侧边主题"`                                          // 用户侧边主题
	// HeaderImg   string         `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
	// BaseColor   string         `json:"baseColor" gorm:"default:#fff;comment:基础颜色"`                                           // 基础颜色
	// ActiveColor string         `json:"activeColor" gorm:"default:#1890ff;comment:活跃颜色"`                                      // 活跃颜色
	// AuthorityId uint           `json:"authorityId" gorm:"default:888;comment:用户角色ID"`                                        // 用户角色ID
	// Authority   SysAuthority   `json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:用户角色"`
	// Authorities []SysAuthority `json:"authorities" gorm:"many2many:sys_user_authority;"`
	// Phone       string         `json:"phone"  gorm:"comment:用户手机号"`                     // 用户手机号
	Email    string `gorm:"column:email;size:255;uniqueindex;not null;comment:用户邮箱" json:"email"`               // 用户邮箱
	IsActive int    `gorm:"column:is_active; size:1;not null;default:1;comment:用户是否被冻结 1正常 0冻结" json:"enable" ` //用户是否被冻结 1正常 2冻结
}

func (user User) ToPublic() UserPublic {
	return UserPublic{
		UserID:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		RoleID:   user.RoleID,
		Email:    user.Email,
		IsActive: user.IsActive,
	}
}

type UserPublic struct {
	UserID   uint64
	Username string
	Nickname string
	RoleID   uint
	Email    string
	IsActive int
}

func (User) TableName() string {
	return "User"
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if err := user.HashPassword(); err != nil {
		return err
	}
	return nil
}

// HashPassword hashes the password using bcrypt.
func (user *User) HashPassword() error {
	passwordHash, err := utils.BcryptHash(string(user.Password))
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	return nil
}

// CheckPassword checks if the provided password matches the stored password.
func (user *User) CheckPassword(password string) bool {
	return utils.BcryptCheck(user.Password, password)
}

// SetActive sets the user as active.
func (user *User) SetActive() {
	user.IsActive = 1
	user.UpdatedAt = time.Now()
}
