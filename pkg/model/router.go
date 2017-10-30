package model

import (
	"encoding/json"
	"time"

	"fmt"
	uuid "github.com/satori/go.uuid"
)

//Router is struct for keep user proxy task
type Router struct {
	tableName struct{} `sql:"routers,alias:router"`

	ID   string `sql:"type:uuid" json:"-"`
	Name string

	Group   *Group
	GroupID string `json:"-"`

	Roles   []string `pg:",array"`
	OAuth   bool
	Active  bool
	Created time.Time

	StripPath   bool
	ListenPath  string
	UpstreamURL string
	Methods     []string `pg:",array"`

	BeforePlugins []string `pg:",array"`
	AfterPlugins  []string `pg:",array"`
}

//ConvertToJSON return JSON or empty if error
func (r *Router) ConvertToJSON() []byte {
	b, err := json.Marshal(r)
	if err != nil {
		return []byte{}
	}
	return b
}

//CreateDefaultRouter return default Router
func CreateDefaultRouter(name string) *Router {
	return &Router{
		ID:        uuid.NewV4().String(),
		Name:      name,
		Active:    false,
		OAuth:     true,
		Roles:     []string{"user"},
		StripPath: true,
		Created:   time.Now(),
		Group:     &Group{},
	}
}

//CreateRouterFromJSON return Router from JSON
func CreateRouterFromJSON(js []byte) (*Router, error) {
	r := CreateDefaultRouter("default")
	if err := json.Unmarshal(js, r); err != nil {
		return r, err
	}
	return r, nil
}

//BeforeInsert set GroupID
func (r *Router) BeforeInsert() {
	if r.Group != nil {
		fmt.Print(r.Group)
		r.GroupID = r.Group.ID
	}
}
