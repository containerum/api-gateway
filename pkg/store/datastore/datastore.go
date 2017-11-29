package datastore

import (
	"bitbucket.org/exonch/ch-gateway/pkg/store/migrations"
	"github.com/jinzhu/gorm"
)

type datastore struct {
	*gorm.DB
}

//New returns Store with gorm DB
func New(db *gorm.DB) interface{} {
	return &datastore{db}
}

//Init Create migration table
func (d *datastore) Init() error {
	// return d.AutoMigrate(&migrations.Migration{}).Error
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

//Init Create migration table
func (d *datastore) Up() error {
	migrations.RunMigrations(d.DB)
	return nil
}

//NOTE: Example work with gorm
// l := &model.Listener{}
// db.First(l)
//
// fmt.Print(*l)

// listener := &model.Listener{
// 	Name:        "newModel",
// 	StripPath:   true,
// 	ListenPath:  "/yy/*",
// 	UpstreamURL: "http://192.168.88.57:8888",
// 	Methods:     []string{"get", "post"},
// }
// db.NewRecord(listener)
// db.Create(listener)
