package store

import (
	"fmt"

	"git.containerum.net/ch/api-gateway/pkg/model"
	"github.com/jinzhu/gorm"
)

//TODO: Write DB logic
func newConnection(config model.DatabaseConfig) (*gorm.DB, error) {
	return gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.Address,
		config.Port,
		config.User,
		config.Database,
		config.Password,
	))
}
