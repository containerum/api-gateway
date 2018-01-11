// +build ignore

package manage

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"git.containerum.net/ch/api-gateway/pkg/model"

	log "github.com/Sirupsen/logrus"
)

var (
	l1 = model.Listener{
		DefaultModel: model.DefaultModel{
			ID:        "11472afa-48ff-4583-b2e4-74119993f22a",
			CreatedAt: time.Date(2017, 12, 1, 12, 30, 0, 0, time.Local),
			UpdatedAt: time.Date(2017, 12, 1, 12, 40, 0, 0, time.Local),
		},
		Name:        "Router N1",
		OAuth:       newBool(false),
		Active:      newBool(false),
		Group:       g1,
		GroupRefer:  "7c7ebd79-8e9d-443e-8c20-d43de4e33b45",
		StripPath:   newBool(false),
		ListenPath:  "route1",
		UpstreamURL: "http://localhost",
		Method:      "GET",
	}
	l2 = model.Listener{
		DefaultModel: model.DefaultModel{
			ID:        "fdbcf79a-c254-48ec-b1df-60a0984fab5e",
			CreatedAt: time.Date(2017, 12, 1, 12, 30, 0, 0, time.Local),
			UpdatedAt: time.Date(2017, 12, 1, 12, 30, 0, 0, time.Local),
		},
		Name:        "Router N2",
		OAuth:       newBool(false),
		Active:      newBool(true),
		Group:       g1,
		GroupRefer:  "7c7ebd79-8e9d-443e-8c20-d43de4e33b45",
		StripPath:   newBool(false),
		ListenPath:  "route2",
		UpstreamURL: "http://localhost",
		Method:      "GET",
	}
	l3 = model.Listener{
		DefaultModel: model.DefaultModel{
			ID:        "126a743d-add6-4de9-971d-6691f316d530",
			CreatedAt: time.Date(2017, 12, 1, 12, 00, 0, 0, time.Local),
			UpdatedAt: time.Date(2017, 12, 1, 12, 00, 0, 0, time.Local),
		},
		Name:        "Router N3",
		OAuth:       newBool(true),
		Active:      newBool(false),
		Group:       g1,
		GroupRefer:  "7c7ebd79-8e9d-443e-8c20-d43de4e33b45",
		StripPath:   newBool(false),
		ListenPath:  "route3",
		UpstreamURL: "http://localhost",
		Method:      "POST",
	}
	l4 = model.Listener{
		DefaultModel: model.DefaultModel{
			ID:        "9fd10b6a-f58c-4d3f-90b9-0990ee80f684",
			CreatedAt: time.Date(2017, 13, 1, 12, 30, 0, 0, time.Local),
			UpdatedAt: time.Date(2017, 16, 1, 12, 40, 0, 0, time.Local),
		},
		Name:        "Router N4",
		OAuth:       newBool(true),
		Active:      newBool(true),
		Group:       g2,
		GroupRefer:  "87b0564f-0a54-4ba5-b2bb-ad8164919d77",
		StripPath:   newBool(false),
		ListenPath:  "route4",
		UpstreamURL: "http://localhost",
		Method:      "DELETE",
	}
	g1 = model.Group{
		DefaultModel: model.DefaultModel{
			ID:        "7c7ebd79-8e9d-443e-8c20-d43de4e33b45",
			CreatedAt: time.Date(2017, 13, 1, 12, 30, 0, 0, time.Local),
			UpdatedAt: time.Date(2017, 16, 1, 13, 50, 0, 0, time.Local),
		},
		Name: "Group N1",
	}
	g2 = model.Group{
		DefaultModel: model.DefaultModel{
			ID:        "87b0564f-0a54-4ba5-b2bb-ad8164919d77",
			CreatedAt: time.Date(2017, 13, 1, 12, 32, 0, 0, time.Local),
			UpdatedAt: time.Date(2017, 16, 1, 12, 32, 0, 0, time.Local),
		},
		Name: "Group N2",
	}
)

//GetAllRouter return listeners list
func (m manage) GetAllRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqName := "Get all router stub"
		listeners := append([]model.Listener{},
			l1, l2, l3, l4,
		)
		WriteAnswer(http.StatusOK, &listeners, nil, reqName, &w)
	}
}

//GetRouter return listeners by id
func (m manage) GetRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqName := "Get router stub"
		id := chi.URLParam(r, "id")
		switch id {
		case l1.ID:
			WriteAnswer(http.StatusOK, &l1, nil, reqName, &w)
		case l2.ID:
			WriteAnswer(http.StatusOK, &l2, nil, reqName, &w)
		case l3.ID:
			WriteAnswer(http.StatusOK, &l3, nil, reqName, &w)
		case l4.ID:
			WriteAnswer(http.StatusOK, &l4, nil, reqName, &w)
		default:
			e := errors.New("Unable to find listener")
			WriteAnswer(http.StatusBadRequest, nil, &[]error{e}, reqName, &w)
		}
	}
}

//CreateRouter return listeners list
func (m manage) CreateRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqName := "Create router stub"
		decoder := json.NewDecoder(r.Body)
		var listenerNew model.Listener
		if err := decoder.Decode(&listenerNew); err != nil {
			e := errors.New("Unable to decode listener")
			log.WithError(err).Error("Unable to decode listener")
			WriteAnswer(http.StatusBadRequest, nil, &[]error{e}, reqName, &w)
			return
		}
		if err := listenerNew.ValidateCreate(); len(err) != 0 {
			WriteAnswer(http.StatusBadRequest, nil, &err, reqName, &w)
			return
		}
		switch listenerNew.GroupRefer {
		case g1.ID:
			listenerNew.Group = g1
		case g2.ID:
			listenerNew.Group = g2
		default:
			e := errors.New("Unable to find group")
			WriteAnswer(http.StatusInternalServerError, nil, &[]error{e}, reqName, &w)
			return
		}
		listenerNew.CreatedAt = time.Now()
		listenerNew.UpdatedAt = time.Now()
		WriteAnswer(http.StatusOK, listenerNew, nil, reqName, &w)
	}
}

//UpdateRouter return listeners list
func (m manage) UpdateRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqName := "Update router stub"
		decoder := json.NewDecoder(r.Body)
		id := chi.URLParam(r, "id")
		if id != l1.ID || id != l2.ID || id != l3.ID || id != l4.ID {
			e := errors.New("Unable to find listener")
			WriteAnswer(http.StatusNoContent, nil, &[]error{e}, reqName, &w)
			return
		}
		var listenerNew model.Listener
		if err := decoder.Decode(&listenerNew); err != nil {
			WriteAnswer(http.StatusBadRequest, nil, &[]error{err}, reqName, &w)
			return
		}
		update, err := listenerNew.GetUpdateType(id)
		switch update {
		case model.ListenerUpdateFull:
		case model.ListenerUpdateActive:
		case model.ListenerUpdateOAuth:
		case model.ListenerUpdateNone:
			WriteAnswer(http.StatusBadRequest, nil, &err, reqName, &w)
		}
	}
}

//DeleteRouter return listeners list
func (m manage) DeleteRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func newBool(b bool) *bool {
	res := new(bool)
	*res = b
	return res
}
