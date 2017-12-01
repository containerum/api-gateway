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
