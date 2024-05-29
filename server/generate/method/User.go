package method

import (
	"github.com/Lazyn0tBug/beacon/server/model"
	"gorm.io/gen"
)

type UserMethod interface {
	// GetUserList retrieves all users excluding the password field, also db.Omit("password").Find(&users)
	//
	// SELECT * FROM @@table OMIT (password)
	GetUserList() ([]*model.UserPublic, error)

	// get all users who role is above some role
	//
	// WHERE roleID > @roleID
	GetUsersAboveRole(roleID int) ([]gen.T, error)

	// SetActive sets the user as active and updates the updated_at timestamp.
	//
	// UPDATE @@table SET is_active = 1 WHERE id = @id
	SetActive(id uint) error

	// SetActive sets the user as active and updates the updated_at timestamp.
	//
	// UPDATE @@table SET is_active = 0 WHERE id = @id
	SetInActive(id uint) error
}
