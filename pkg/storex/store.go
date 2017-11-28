package storex

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"bitbucket.org/exonch/ch-gateway/pkg/model"
)

type Storex interface {
}

func New(config model.DatabaseConfig) (Storex, error) {
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.Address,
		config.Port,
		config.User,
		config.Database,
		config.Password,
	))

	db.LogMode(true)
	db.AutoMigrate(&model.Role{}, &model.Listener{}) //"Save" migrations

	listener := &model.Listener{
		Name:        "newModel",
		StripPath:   true,
		ListenPath:  "/yy/*",
		UpstreamURL: "http://192.168.88.57:8888",
		//Methods:     []string{"get", "post"},
	}
	db.NewRecord(listener)
	db.Create(listener)

	return db, err
}
