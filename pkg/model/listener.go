package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"git.containerum.net/ch/api-gateway/pkg/utils/snake"
	model "git.containerum.net/ch/json-types/gateway"

	validator "github.com/asaskevich/govalidator"
)

//ListenerUpdateType keeps method type for update
type ListenerUpdateType int

const (
	listenerNameLengthMin        = 4
	listenerNameLengthMax        = 128
	listenerListenPathLengthMin  = 3
	listenerListenPathLengthMax  = 128
	listenerUpstreamURLLengthMin = 3
	listenerUpstreamURLLengthMax = 128

	//ListenerUpdateNone when nothing to update
	ListenerUpdateNone ListenerUpdateType = iota
	//ListenerUpdateFull when update all params
	ListenerUpdateFull
	//ListenerUpdateActive when update only active
	ListenerUpdateActive
	//ListenerUpdateOAuth when update only Oauth
	ListenerUpdateOAuth
)

var (
	listenerMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
)

//Listener keeps proxy-router configs
type Listener struct {
	DefaultModel
	Name      string `gorm:"not null"`
	nameSnake string

	// Roles  []Role
	OAuth  *bool
	Active *bool

	Group      Group `gorm:"ForeignKey:GroupRefer"`
	GroupRefer string

	StripPath   *bool
	ListenPath  string `gorm:"not null"`
	UpstreamURL string `gorm:"not null"`
	Method      string `gorm:"not null"`

	// Plugins []Plugin
}

//GetSnakeName return Listerner name in snake case
func (l *Listener) GetSnakeName() string {
	return l.nameSnake
}

/* +++ JSON +++ */

//MarshalJSON return Listerner model json
func (l Listener) MarshalJSON() ([]byte, error) {
	listener := model.ListenerJSON{
		ID:        &l.ID,
		CreatedAt: getUnixTimeOrNil(&l.CreatedAt),
		UpdatedAt: getUnixTimeOrNil(&l.UpdatedAt),
		DeletedAt: getUnixTimeOrNil(l.DeletedAt),
		Name:      &l.Name,
		OAuth:     l.OAuth,
		Active:    l.Active,
		Group: &model.GroupJSON{
			ID:        &l.Group.ID,
			CreatedAt: getUnixTimeOrNil(&(l.Group).CreatedAt),
			UpdatedAt: getUnixTimeOrNil(&(l.Group).UpdatedAt),
			DeletedAt: getUnixTimeOrNil(l.Group.DeletedAt),
			Name:      &l.Group.Name,
			Active:    l.Group.Active,
		},
		StripPath:   l.StripPath,
		ListenPath:  &l.ListenPath,
		UpstreamURL: &l.UpstreamURL,
		Method:      &l.Method,
	}
	return json.Marshal(listener)
}

//UnmarshalJSON make Listener from json
func (l *Listener) UnmarshalJSON(b []byte) error {
	listenerJS := &model.ListenerJSON{}
	if err := json.Unmarshal(b, listenerJS); err != nil {
		return err
	}
	if listenerJS.ID != nil {
		l.ID = *listenerJS.ID
	}
	if listenerJS.CreatedAt != nil {
		l.CreatedAt = time.Unix(*listenerJS.CreatedAt, 0)
	}
	if listenerJS.UpdatedAt != nil {
		l.UpdatedAt = time.Unix(*listenerJS.UpdatedAt, 0)
	}
	if listenerJS.DeletedAt != nil {
		var unix = time.Unix(*listenerJS.DeletedAt, 0)
		l.DeletedAt = &unix
	}
	if listenerJS.Name != nil {
		l.Name = *listenerJS.Name
	}
	if listenerJS.ListenPath != nil {
		l.ListenPath = *listenerJS.ListenPath
	}
	if listenerJS.UpstreamURL != nil {
		l.UpstreamURL = *listenerJS.UpstreamURL
	}
	if listenerJS.Method != nil {
		l.Method = *listenerJS.Method
	}
	if listenerJS.GroupID != nil {
		l.GroupRefer = *listenerJS.GroupID
	}
	l.OAuth = listenerJS.OAuth
	l.Active = listenerJS.Active
	l.StripPath = listenerJS.StripPath
	return nil
}

//AfterFind runs after select in DB, write nameSnake value
func (l *Listener) AfterFind() (err error) {
	l.nameSnake = snake.StrToSnake(l.Name)
	return nil
}

/* +++ Validation +++ */
var (
	//ErrNilListener - error when listener model is nil
	ErrNilListener = errors.New("Listener model is empty")
	//ErrInvalidListenerNameLength - error when name lenght < listenerNameLengthMin or > listenerNameLengthMax
	ErrInvalidListenerNameLength = fmt.Errorf("Invalid Name length. It must be more than %v and less than %v", listenerNameLengthMin, listenerNameLengthMax)
	//ErrInvalidListenerID - error when id is not uuid
	ErrInvalidListenerID = errors.New("Invalid Listener ID. It must be UUID")
	//ErrInvalidListenerMethod - error when method not supported
	ErrInvalidListenerMethod = fmt.Errorf("Invalid method. It must be one of %v", listenerMethods)
	//ErrInvalidListenerGroupRefer - error when group refer is not uuid
	ErrInvalidListenerGroupRefer = errors.New("Invalid Group refer ID in Listener. It must be UUID")
	//ErrInvalidListenerListenPath - error when listen path < listenerListenPathLengthMin or > listenerListenPathLengthMax
	ErrInvalidListenerListenPath = fmt.Errorf("Invalid Listen Path length. It must be more than %v and less than %v", listenerListenPathLengthMin, listenerListenPathLengthMax)
	//ErrInvalidListenerUpstreamURL - error when upstream lenght < listenerUpstreamURLLengthMin or > listenerUpstreamURLLengthMax
	ErrInvalidListenerUpstreamURL = fmt.Errorf("Invalid Upstream URL length. It must be more than %v and less than %v", listenerUpstreamURLLengthMin, listenerUpstreamURLLengthMax)
	//ErrInvalidListenerActive - error when Active is nil
	ErrInvalidListenerActive = errors.New("Param Active is empty")
)

//ValidateCreate check model before insert
func (l *Listener) ValidateCreate() (err []error) {
	if l == nil {
		err = append(err, ErrNilListener)
		return
	}
	if !validator.IsByteLength(l.Name, listenerNameLengthMin, listenerNameLengthMax) {
		err = append(err, ErrInvalidListenerNameLength)
	}
	if !validator.IsIn(strings.ToUpper(l.Method), listenerMethods...) {
		err = append(err, ErrInvalidListenerMethod)
	}
	if !validator.IsUUID(l.GroupRefer) {
		err = append(err, ErrInvalidListenerGroupRefer)
	}
	if !validator.IsByteLength(l.ListenPath, listenerListenPathLengthMin, listenerListenPathLengthMax) {
		err = append(err, ErrInvalidListenerListenPath)
	}
	if !validator.IsByteLength(l.UpstreamURL, listenerUpstreamURLLengthMin, listenerUpstreamURLLengthMax) {
		err = append(err, ErrInvalidListenerUpstreamURL)
	}
	return
}

//ValidateUpdate check model before update
func (l *Listener) ValidateUpdate(id string) (err []error) {
	if !validator.IsUUID(id) {
		err = append(err, ErrInvalidListenerID)
	}
	l.ID = id
	err = append(err, l.ValidateCreate()...)
	return
}

//ValidateUpdateActive check if possible to update Active field
func (l *Listener) ValidateUpdateActive(id string) (err []error) {
	if l == nil {
		err = append(err, ErrNilListener)
		return
	}
	if !validator.IsUUID(id) {
		err = append(err, ErrInvalidListenerID)
	}
	l.ID = id
	if l.Active == nil {
		err = append(err, ErrInvalidListenerActive)
	}
	return
}

//ValidateUpdateOAuth check if possible to update OAuth field
func (l *Listener) ValidateUpdateOAuth(id string) (err []error) {
	if l == nil {
		err = append(err, ErrNilListener)
		return
	}
	if !validator.IsUUID(id) {
		err = append(err, ErrInvalidListenerID)
	}
	l.ID = id
	if l.OAuth == nil {
		err = append(err, ErrInvalidListenerActive)
	}
	return
}

//GetUpdateType return update method
func (l *Listener) GetUpdateType(id string) (ListenerUpdateType, []error) {
	var err []error
	if err = l.ValidateUpdate(id); len(err) == 0 {
		return ListenerUpdateFull, err
	}
	if err = l.ValidateUpdateActive(id); len(err) == 0 {
		return ListenerUpdateActive, err
	}
	if err = l.ValidateUpdateOAuth(id); len(err) == 0 {
		return ListenerUpdateOAuth, err
	}
	return ListenerUpdateNone, err
}
