package util

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// Entity an identifiable
type Entity interface {
	GetID() uint
}

// EntityImpl struct for DB objects
type EntityImpl struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// GetID returns the identifier of the resource
func (e *EntityImpl) GetID() uint {
	return e.ID
}

// DatabaseConfig configuration for the application database
type DatabaseConfig struct {
	Type    string `mapstructure:"type"`
	File    string `mapstructure:"file"`
	Migrate bool   `mapstructure:"migrate"`
}

// initializes the GORM database
func initDatabase(databaseConfig *DatabaseConfig) (*gorm.DB, error) {

	var err error
	var db *gorm.DB

	if databaseConfig.Type == "" {
		return nil, fmt.Errorf("Database type is required")
	}

	switch databaseConfig.Type {
	case "sqlite3":
		if databaseConfig.File == "" {
			return nil, fmt.Errorf("Database file is required")
		}
		db, err = gorm.Open(databaseConfig.Type, databaseConfig.File)
		if err != nil {
			return nil, fmt.Errorf("Error initializing database: %s", err.Error())
		}

	default:
		return nil, fmt.Errorf("database %s is not supported", databaseConfig.Type)
	}

	// if databaseConfig.Migrate {
	// 	db.AutoMigrate(&users.User{})
	// }

	return db, nil
}
