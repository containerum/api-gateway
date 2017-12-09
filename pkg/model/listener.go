package model

import (
	"encoding/json"
	"unicode"

	"github.com/fatih/structs"
)

//Listener keeps proxy-router configs
type Listener struct {
	DefaultModel `structs:"-"`
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

// func (lm *Listener) BeforeUpdate() (err error) {
// 	fmt.Printf("\nUpdate: %v\n", lm.UpdatedAt)
// 	return
// }

func (l *Listener) AfterFind() (err error) {
	l.nameSnake = toSnake(l.Name)
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
