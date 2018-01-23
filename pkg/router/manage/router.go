package manage

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.containerum.net/ch/api-gateway/pkg/model"

	"github.com/go-chi/chi"
)

const (
	getAllRouterMethod = "Get all router"
	getRouterMethod    = "Get router"
	createRouterMethod = "Create router"
	updateRouterMethod = "Update router"
	deleteRouterMethod = "Delete router"
)

var (
	//ErrUnableFindLisneners - error when unable get listeners from db
	ErrUnableFindLisneners = errors.New("Unable to find listeners")
	//ErrUnableFindLisnener - error when unable get listener from db
	ErrUnableFindLisnener = errors.New("Unable to find listener")
	//ErrUnableDecodeListener - error when unable to decode listener json
	ErrUnableDecodeListener = errors.New("Unable to decode listener")
)

//GetAllRouter return listeners list
func (m manage) GetAllRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		listeners, err := (*m.st).GetListenerList(nil)
		if err != nil {
			WriteAnswer(http.StatusBadRequest, getAllRouterMethod, &w, nil, ErrUnableFindLisneners)
			return
		}
		WriteAnswer(http.StatusOK, getAllRouterMethod, &w, listeners)
	}
}

//GetRouter return listeners by id
func (m manage) GetRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		listener, err := (*m.st).GetListener(id)
		if err != nil {
			WriteAnswer(http.StatusBadRequest, getRouterMethod, &w, nil, ErrUnableFindLisnener)
			return
		}
		WriteAnswer(http.StatusOK, getRouterMethod, &w, listener)
	}
}

//CreateRouter create listener id db
func (m manage) CreateRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var listenerNew model.Listener
		if err := decoder.Decode(&listenerNew); err != nil {
			WriteAnswer(http.StatusBadRequest, createRouterMethod, &w, nil, ErrUnableDecodeListener)
			return
		}
		if err := listenerNew.Validate(); len(err) != 0 {
			WriteAnswer(http.StatusBadRequest, createRouterMethod, &w, nil, err...)
			return
		}
		listenerSaved, err := (*m.st).CreateListener(&listenerNew)
		if err != nil {
			WriteAnswer(http.StatusInternalServerError, createRouterMethod, &w, nil, err)
			return
		}
		WriteAnswer(http.StatusOK, createRouterMethod, &w, listenerSaved)
	}
}

//UpdateRouter update router
func (m manage) UpdateRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		id := chi.URLParam(r, "id")
		var listenerNew model.Listener
		if _, err := (*m.st).GetListener(id); err != nil {
			WriteAnswer(http.StatusNoContent, updateRouterMethod, &w, nil, ErrUnableFindLisnener)
			return
		}
		if err := decoder.Decode(&listenerNew); err != nil {
			WriteAnswer(http.StatusBadRequest, updateRouterMethod, &w, nil, ErrUnableDecodeListener)
			return
		}
		if err := listenerNew.Validate(); len(err) != 0 {
			WriteAnswer(http.StatusBadRequest, createRouterMethod, &w, nil, err...)
			return
		}
		listenerNew.ID = id
		if err := (*m.st).UpdateListener(&listenerNew); err != nil {
			WriteAnswer(http.StatusInternalServerError, updateRouterMethod, &w, nil, err)
			return
		}
		WriteAnswer(http.StatusOK, updateRouterMethod, &w, listenerNew)
	}
}

//DeleteRouter delete router
func (m manage) DeleteRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if err := (*m.st).DeleteListener(id); err != nil {
			WriteAnswer(http.StatusInternalServerError, deleteRouterMethod, &w, nil, err)
			return
		}
		WriteAnswer(http.StatusOK, deleteRouterMethod, &w, nil)
	}
}
