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
	groupNameLengthMin = 3
	groupNameLengthMax = 32
)

var (
	//ErrNilGroup - error when group model is nil
	ErrNilGroup = errors.New("Group model is empty")
	//ErrInvalidGroupNameLength - error when name lenght < groupNameLengthMin or > groupNameLengthMax
	ErrInvalidGroupNameLength = fmt.Errorf("Invalid Name length. It must be more than %v and less than %v", groupNameLengthMin, groupNameLengthMax)
	//ErrInvalidGroupID - error when id is not uuid
	ErrInvalidGroupID = errors.New("Invalid Group ID. It must be UUID")
)

//Group struct identificate router group
type Group struct {
	DefaultModel
	Name   string `db:"name"`
	Active bool   `db:"active"`
}

//MarshalJSON return Group model json
func (g Group) MarshalJSON() ([]byte, error) {
	group := model.GroupJSON{
		ID:        &g.ID,
		CreatedAt: getUnixTimeOrNil(&g.CreatedAt),
		UpdatedAt: getUnixTimeOrNil(&g.UpdatedAt),
		Name:      &g.Name,
		Active:    &g.Active,
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
	if groupJS.Name != nil {
		g.Name = *groupJS.Name
	}
	if groupJS.Active != nil {
		g.Active = *groupJS.Active
	}
	return nil
}

//Validate check model
func (g *Group) Validate() (err []error) {
	if g == nil {
		err = append(err, ErrNilGroup)
		return
	}
	if !validator.IsByteLength(g.Name, groupNameLengthMin, groupNameLengthMax) {
		err = append(err, ErrInvalidGroupNameLength)
	}
	return
}
