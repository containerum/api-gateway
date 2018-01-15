package store

import (
	"context"

	"git.containerum.net/ch/api-gateway/pkg/model"
)

// FromContext returns the Store associated with this context.
func FromContext(c context.Context) Store {
	return c.Value(key).(Store)
}

func GetListener(c context.Context, id string) (*model.Listener, error) {
	return FromContext(c).GetListener(id)
}

func GetListenerList(c context.Context, l *model.Listener) (*[]model.Listener, error) {
	return FromContext(c).GetListenerList(l)
}

func UpdateListener(c context.Context, l *model.Listener, utype model.ListenerUpdateType) error {
	return FromContext(c).UpdateListener(l, utype)
}

func CreateListener(c context.Context, l *model.Listener) (*model.Listener, error) {
	return FromContext(c).CreateListener(l)
}

func DeleteListener(c context.Context, id string) error {
	return FromContext(c).DeleteListener(id)
}

func GetGroupList(c context.Context, g *model.Group) (*[]model.Group, error) {
	return FromContext(c).GetGroupList(g)
}

func CreateGroup(c context.Context, g *model.Group) (*model.Group, error) {
	return FromContext(c).CreateGroup(g)
}
