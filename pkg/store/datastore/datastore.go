// +build !stub

package datastore

import (
	"fmt"

	"git.containerum.net/ch/api-gateway/pkg/model"
	"git.containerum.net/ch/api-gateway/pkg/store/datastore/migrations"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //Gorm postgres driver
)

type datastore struct {
	*gorm.DB
}

//New returns Store with gorm DB
func New(config model.DatabaseConfig) (interface{}, error) {
	db, err := newConnection(config)
	if config.Debug {
		db.Debug()
		db.LogMode(true)
	}
	if config.SafeMigration {
		db.AutoMigrate(&model.Role{}, &model.Group{}, &model.Plugin{}, &model.Listener{}) //"Safety" migrations
	}
	return &datastore{db}, err
}

func newConnection(config model.DatabaseConfig) (*gorm.DB, error) {
	return gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.Address,
		config.Port,
		config.User,
		config.Database,
		config.Password,
	))
}

//Init Create migration table
func (d *datastore) Init() error {
	//Create table
	err := d.AutoMigrate(&migrations.Migration{}).Error
	if err != nil {
		return err
	}
	//Write first record if not exists
	var m migrations.Migration
	err = d.First(&m).Error
	if err != nil {
		if err.Error() == "record not found" {
			m.Dirty = false
			m.Version = 0
			d.NewRecord(&m)
			if err = d.Create(&m).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

//Version return curent DB Migration Version
func (d *datastore) Version() (int, error) {
	var m migrations.Migration
	err := d.First(&m).Error
	if err != nil {
		return 0, err
	}
	return m.Version, nil
}

//Up run all migration
func (d *datastore) Up() (int, error) {
	return migrations.RunMigrations(d.DB)
}

//Down run last down migration
func (d *datastore) Down() (int, error) {
	v, err := d.Version()
	if err != nil {
		return 0, err
	}
	return migrations.RunMigrationsDown(v, d.DB)
}
