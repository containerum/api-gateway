package model

import (
	"time"
)

//Router is struct for keep user proxy task
type Router struct {
	tableName struct{} `sql:"routers,alias:router"`

	ID   string `sql:"type:uuid"`
	Name string

	Group   *Group
	GroupID string

	Roles   []string
	OAuth   bool
	Active  bool
	Created time.Time

	StripPath   bool
	ListenPath  string
	UpstreamURL string
	Methods     []string

	BeforePlugins []string
	AfterPlugins  []string
}

//CreateDefaultRouter return default Router
func CreateDefaultRouter() *Router {
	return &Router{
		Active:    false,
		OAuth:     true,
		Roles:     []string{"user"},
		StripPath: true,
	}
}

//BeforeInsert set GroupID
func (r *Router) BeforeInsert() {
	r.GroupID = r.Group.ID
}
