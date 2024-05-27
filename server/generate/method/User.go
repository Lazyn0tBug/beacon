package method

import (
	"github.com/Lazyn0tBug/beacon/server/model"
	"gorm.io/gen"
)

type UserMethod interface {
	gen.Method

	// GetUserList retrieves all users excluding the password field, also db.Omit("password").Find(&users)
	//
	// SELECT * FROM @@table OMIT (password)
	GetUserList() ([]*model.UserPublic, error)

	// SetActive sets the user as active and updates the updated_at timestamp.
	//
	// UPDATE @@table SET is_active = 1, updated_at = now() WHERE id = @id
	SetActive(id uint) error

	// SetActive sets the user as active and updates the updated_at timestamp.
	//
	// UPDATE @@table SET is_active = 0, updated_at = now() WHERE id = @id
	SetInActive(id uint) error
}
