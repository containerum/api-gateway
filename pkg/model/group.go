package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	model "git.containerum.net/ch/json-types/gateway"

	validator "github.com/asaskevich/govalidator"
)

const (
	groupNameLengthMin = 4
	groupNameLengthMax = 32
)

//Group struct identificate router group
type Group struct {
	DefaultModel
	Name   string `gorm:"type:varchar(32);not null;unique"`
	Active *bool  `gorm:"not null"`
}

/* +++ JSON +++ */

//MarshalJSON return Group model json
func (g Group) MarshalJSON() ([]byte, error) {
	group := model.GroupJSON{
		ID:        &g.ID,
		CreatedAt: getUnixTimeOrNil(&g.CreatedAt),
		UpdatedAt: getUnixTimeOrNil(&g.UpdatedAt),
		DeletedAt: getUnixTimeOrNil(g.DeletedAt),
		Name:      &g.Name,
		Active:    g.Active,
	}
	return json.Marshal(group)
}

//UnmarshalJSON make Group from json
func (g Group) UnmarshalJSON(b []byte) error {
	var groupJS model.GroupJSON
	if err := json.Unmarshal(b, groupJS); err != nil {
		return err
	}
	if groupJS.ID != nil {
		g.ID = *groupJS.ID
	}
	if groupJS.CreatedAt != nil {
		g.CreatedAt = time.Unix(*groupJS.CreatedAt, 0)
	}
	if groupJS.UpdatedAt != nil {
		g.UpdatedAt = time.Unix(*groupJS.UpdatedAt, 0)
	}
	if groupJS.DeletedAt != nil {
		var unix = time.Unix(*groupJS.DeletedAt, 0)
		g.DeletedAt = &unix
	}
	if groupJS.Name != nil {
		g.Name = *groupJS.Name
	}
	g.Active = groupJS.Active
	return nil
}

/* +++ Validation +++ */
var (
	//ErrInvalidGroupNameLength - error when name lenght < groupNameLengthMin or > groupNameLengthMax
	ErrInvalidGroupNameLength = fmt.Errorf("Invalid Name length. It must be more than %v and less than %v", groupNameLengthMin, groupNameLengthMax)
	//ErrInvalidGroupID - error when id is not uuid
	ErrInvalidGroupID = errors.New("Invalid Group ID. It must be UUID")
)

//ValidateCreate check model before insert
func (g *Group) ValidateCreate() (err []error) {
	if !validator.IsByteLength(g.Name, groupNameLengthMin, groupNameLengthMax) {
		err = append(err, ErrInvalidGroupNameLength)
	}
	return
}

//ValidateUpdate check model before update
func (g *Group) ValidateUpdate() (err []error) {
	if !validator.IsUUID(g.ID) {
		err = append(err, ErrInvalidGroupID)
	}
	err = append(err, g.ValidateCreate()...)
	return
}
