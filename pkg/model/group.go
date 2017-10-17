package model

//Group struct ident router group
type Group struct {
	tableName struct{} `sql:"groups"`

	ID     string `sql:"type:uuid"`
	Name   string `sql:"type:varchar(32)"`
	Active bool
}
