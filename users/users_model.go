package users

import (
	"github.com/sterrasi/stepwise/util"
)

// User DTO
type User struct {
	util.EntityImpl
	UserName string `gorm:"type:varchar(20);index;not null"`

	FirstName  string `gorm:"type:varchar(20);index;not null"`
	MiddleName string `gorm:"type:varchar(20)"`
	LastName   string `gorm:"type:varchar(50);index;not null"`

	PrimaryEmail string `gorm:"type:varchar(35);unique_index;not null"`
	Organization string `gorm:"type:varchar(10);index;not null"`
	Sex          string `gorm:"type:varchar(1)"`
}

//type UserRegistration struct {
