package users

import (
	"github.com/jinzhu/gorm"
	"github.com/sterrasi/stepwise/util"
)

var (
	db *gorm.DB
)

type UserAttributes struct {
	UserID      uint `gorm:"index:att;not null"`
	AttributeID uint `gorm:"index:att;not null"`
}

// TableName for users
func (User) TableName() string {
	return "users"
}

func newInstance() interface{} {
	return &User{}
}

// Checks the User table for a user with the given email.
// Returns true when the user already exists
func notAlreadyRegistered(email string) (bool, error) {
	user := &User{}
	if err := db.Where(&User{PrimaryEmail: email}).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return true, nil
		} else {
			return false, err
		}
	}
	return false, nil
}

// GetUsers returns a page of users
func GetUsers(offset int, limit int) ([]*User, error) {
	users := make([]*User, limit)
	if err := db.Offset(offset).Limit(limit).Find(users).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, util.ErrNotFound
		} else {
			return nil, err
		}
	}
	return users, nil
}

// GetUser returns a specific user
func GetUser(id int) (util.Entity, error) {
	user := &User{}
	uid := uint(id)
	if err := db.First(user, uid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, util.ErrNotFound
		} else {
			return nil, err
		}
	}
	return user, nil
}

// UpdateUser updates a specified user
func UpdateUser(user interface{}) error {

	return nil
}

// PatchUser does a partial update on the specified user
func PatchUser(user interface{}) error {

	return nil
}

// DeleteUser deletes the user with the specified ID
func DeleteUser(id int) error {

	return nil
}

// CreateUser creates a new user using the specified model
func CreateUser(user interface{}) (util.Entity, error) {

	return nil, nil
}
