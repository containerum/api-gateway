// +build stub

package datastore

import (
	"time"

	"git.containerum.net/ch/api-gateway/pkg/model"

	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

var (
	groups map[string]model.Group = map[string]model.Group{
		"7c7ebd79-8e9d-443e-8c20-d43de4e33b45": model.Group{
			DefaultModel: model.DefaultModel{
				ID:        "7c7ebd79-8e9d-443e-8c20-d43de4e33b45",
				CreatedAt: time.Date(2017, 13, 1, 12, 30, 0, 0, time.Local),
				UpdatedAt: time.Date(2017, 16, 1, 13, 50, 0, 0, time.Local),
			},
			Name: "Group N1",
		},
		"87b0564f-0a54-4ba5-b2bb-ad8164919d77": model.Group{
			DefaultModel: model.DefaultModel{
				ID:        "87b0564f-0a54-4ba5-b2bb-ad8164919d77",
				CreatedAt: time.Date(2017, 13, 1, 12, 32, 0, 0, time.Local),
				UpdatedAt: time.Date(2017, 16, 1, 12, 32, 0, 0, time.Local),
			},
			Name: "Group N2",
		},
	}

	listeners map[string]model.Listener = map[string]model.Listener{
		"11472afa-48ff-4583-b2e4-74119993f22a": model.Listener{
			DefaultModel: model.DefaultModel{
				ID:        "11472afa-48ff-4583-b2e4-74119993f22a",
				CreatedAt: time.Date(2017, 12, 1, 12, 30, 0, 0, time.Local),
				UpdatedAt: time.Date(2017, 12, 1, 12, 40, 0, 0, time.Local),
			},
			Name:        "Router N1",
			OAuth:       newBool(false),
			Active:      newBool(false),
			Group:       groups["7c7ebd79-8e9d-443e-8c20-d43de4e33b45"],
			GroupRefer:  "7c7ebd79-8e9d-443e-8c20-d43de4e33b45",
			StripPath:   newBool(false),
			ListenPath:  "/route1",
			UpstreamURL: "http://localhost",
			Method:      "GET",
		},
		"fdbcf79a-c254-48ec-b1df-60a0984fab5e": model.Listener{
			DefaultModel: model.DefaultModel{
				ID:        "fdbcf79a-c254-48ec-b1df-60a0984fab5e",
				CreatedAt: time.Date(2017, 12, 1, 12, 30, 0, 0, time.Local),
				UpdatedAt: time.Date(2017, 12, 1, 12, 30, 0, 0, time.Local),
			},
			Name:        "Router N2",
			OAuth:       newBool(false),
			Active:      newBool(true),
			Group:       groups["7c7ebd79-8e9d-443e-8c20-d43de4e33b45"],
			GroupRefer:  "7c7ebd79-8e9d-443e-8c20-d43de4e33b45",
			StripPath:   newBool(false),
			ListenPath:  "/route2",
			UpstreamURL: "http://localhost",
			Method:      "GET",
		},
		"126a743d-add6-4de9-971d-6691f316d530": model.Listener{
			DefaultModel: model.DefaultModel{
				ID:        "126a743d-add6-4de9-971d-6691f316d530",
				CreatedAt: time.Date(2017, 12, 1, 12, 00, 0, 0, time.Local),
				UpdatedAt: time.Date(2017, 12, 1, 12, 00, 0, 0, time.Local),
			},
			Name:        "Router N3",
			OAuth:       newBool(true),
			Active:      newBool(false),
			Group:       groups["7c7ebd79-8e9d-443e-8c20-d43de4e33b45"],
			GroupRefer:  "7c7ebd79-8e9d-443e-8c20-d43de4e33b45",
			StripPath:   newBool(false),
			ListenPath:  "/route3",
			UpstreamURL: "http://localhost",
			Method:      "POST",
		},
		"9fd10b6a-f58c-4d3f-90b9-0990ee80f684": model.Listener{
			DefaultModel: model.DefaultModel{
				ID:        "9fd10b6a-f58c-4d3f-90b9-0990ee80f684",
				CreatedAt: time.Date(2017, 13, 1, 12, 30, 0, 0, time.Local),
				UpdatedAt: time.Date(2017, 16, 1, 12, 40, 0, 0, time.Local),
			},
			Name:        "Router N4",
			OAuth:       newBool(true),
			Active:      newBool(true),
			Group:       groups["87b0564f-0a54-4ba5-b2bb-ad8164919d77"],
			GroupRefer:  "87b0564f-0a54-4ba5-b2bb-ad8164919d77",
			StripPath:   newBool(false),
			ListenPath:  "/route4",
			UpstreamURL: "http://localhost",
			Method:      "DELETE",
		},
	}
)

//GetListener find Listener by ID
func (d *datastore) GetListener(id string) (*model.Listener, error) {
	reqName := "GetListener stub call"
	if l, ok := listeners[id]; !ok {
		log.WithError(ErrUnableFindListener).Debug(reqName)
		return nil, ErrUnableFindListener
	} else {
		log.Debug(reqName)
		return &l, nil
	}
}

//GetListenerList find all listeers
func (d *datastore) GetListenerList(l *model.Listener) (*[]model.Listener, error) {
	var listenersArr []model.Listener
	for _, l := range listeners {
		listenersArr = append(listenersArr, l)
	}
	log.Debug("GetListenerList stub call")
	return &listenersArr, nil
}

//UpdateListener updates model in DB
func (d *datastore) UpdateListener(l *model.Listener, utype model.ListenerUpdateType) error {
	reqName := "UpdateListener stub call"
	if listener, err := d.GetListener(l.ID); err != nil {
		log.WithError(err).Debug(reqName)
		return err
	} else {
		switch utype {
		case model.ListenerUpdateActive:
			listener.Active = l.Active
			log.WithField("UpdateType", "ListenerUpdateActive").Debug(reqName)
		case model.ListenerUpdateOAuth:
			listener.OAuth = l.OAuth
			log.WithField("UpdateType", "ListenerUpdateOAuth").Debug(reqName)
		case model.ListenerUpdateFull:
			listener.Name = l.Name
			listener.Method = l.Method
			listener.GroupRefer = l.GroupRefer
			listener.ListenPath = l.ListenPath
			listener.UpstreamURL = l.UpstreamURL
			log.WithField("UpdateType", "ListenerUpdateFull").Debug(reqName)
		}
		listener.UpdatedAt = time.Now()
		listeners[l.ID] = *listener
	}
	return nil
}

//TODO Check GroupID
//CreateListener create new listener in DB
func (d *datastore) CreateListener(l *model.Listener) (*model.Listener, error) {
	reqName := "CreateListener stub call"
	id := uuid.NewV4().String()
	listener := model.Listener{
		DefaultModel: model.DefaultModel{
			ID:        id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        l.Name,
		OAuth:       l.OAuth,
		Active:      l.Active,
		Group:       groups[l.GroupRefer],
		GroupRefer:  l.GroupRefer,
		StripPath:   l.StripPath,
		ListenPath:  l.ListenPath,
		UpstreamURL: l.UpstreamURL,
		Method:      l.Method,
	}
	listeners[id] = listener
	log.Debug(reqName)
	return &listener, nil
}

//DeleteListener delete listener in DB by ID
func (d *datastore) DeleteListener(id string) error {
	reqName := "DeleteListener stub call"
	if _, err := d.GetListener(id); err != nil {
		log.WithError(err).Debug(reqName)
		return err
	}
	delete(listeners, id)
	log.WithField("Id", id).Debug(reqName)
	return nil
}

func newBool(b bool) *bool {
	res := new(bool)
	*res = b
	return res
}
