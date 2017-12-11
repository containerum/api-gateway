package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"github.com/fatih/structs"
)

const (
	minNameLenght       = 3
	minListenPathLenght = 3
)

//Listener keeps proxy-router configs
type Listener struct {
	DefaultModel `structs:"-" json:"-"`
	Name         string `gorm:"not null" structs:"name"`
	nameSnake    string

	Roles []Role `structs:"roles"`
	//TODO: Add Group
	OAuth  bool `structs:"o_auth"`
	Active bool `structs:"active"`

	StripPath   bool   `structs:"strip_path"`
	ListenPath  string `gorm:"not null" structs:"listen_path"`
	UpstreamURL string `gorm:"not null" structs:"upstream_url"`
	Method      string `gorm:"not null" structs:"method"`

	Plugins []Plugin `structs:"plugins"`
}

func (l Listener) MarshalJSON() ([]byte, error) {
	listenerMap := structs.Map(l)
	for k, v := range l.getMap() {
		listenerMap[k] = v
	}
	return json.Marshal(listenerMap)
}

func (l *Listener) GetSnakeName() string {
	return l.nameSnake
}

func (l *Listener) AfterFind() (err error) {
	l.nameSnake = toSnake(l.Name)
	return nil
}

func (l *Listener) Valid() (err error) {
	if len(l.Name) < minNameLenght {
		return fmt.Errorf("Route name is too short. Min lenght: %d", minNameLenght)
	}
	if len(l.ListenPath) < minListenPathLenght {
		return fmt.Errorf("ListenPath is too short. Min lenght: %d", minListenPathLenght)
	}
	switch strings.ToLower(l.Method) {
	case "get", "post", "delete", "path", "options", "put":
		break
	default:
		return fmt.Errorf("Unsupported method: %s", l.Method)
	}
	return nil
}

func toSnake(in string) string {
	runes := []rune(in)
	length := len(runes)
	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}
	return string(out)
}
