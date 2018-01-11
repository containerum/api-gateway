// +build stub

package datastore

import (
	"git.containerum.net/ch/api-gateway/pkg/model"

	"github.com/jinzhu/gorm"
)

type datastore struct {
	*gorm.DB
}

var (
	version = 1
)

//New returns Stub Store
func New(config model.DatabaseConfig) (interface{}, error) {
	return &datastore{&gorm.DB{}}, nil
}

//Init Create migration table
func (d *datastore) Init() error {
	return nil
}

//Version return curent DB Migration Version
func (d *datastore) Version() (int, error) {
	return version, nil
}

//Up run all migration
func (d *datastore) Up() (int, error) {
	version++
	return version, nil
}

//Down run last down migration
func (d *datastore) Down() (int, error) {
	version--
	return version, nil
}
