package manage

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"

	"git.containerum.net/ch/api-gateway/pkg/model"
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
		reqName := "Get all router"
		listeners, err := (*m.st).GetListenerList(&model.Listener{})
		if err != nil {
			WriteAnswer(http.StatusBadRequest, nil, &[]error{ErrUnableFindLisneners}, reqName, &w)
			return
		}
		WriteAnswer(http.StatusOK, listeners, nil, reqName, &w)
	}
}

//GetRouter return listeners by id
func (m manage) GetRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqName := "Get router"
		id := chi.URLParam(r, "id")
		listener, err := (*m.st).GetListener(id)
		if err != nil {
			WriteAnswer(http.StatusBadRequest, nil, &[]error{ErrUnableFindLisnener}, reqName, &w)
			return
		}
		WriteAnswer(http.StatusOK, listener, nil, reqName, &w)
	}
}

//CreateRouter create listener id db
func (m manage) CreateRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqName := "Create router"
		decoder := json.NewDecoder(r.Body)
		var listenerNew model.Listener
		if err := decoder.Decode(&listenerNew); err != nil {
			WriteAnswer(http.StatusBadRequest, nil, &[]error{ErrUnableDecodeListener}, reqName, &w)
			return
		}
		if err := listenerNew.ValidateCreate(); len(err) != 0 {
			WriteAnswer(http.StatusBadRequest, nil, &err, reqName, &w)
			return
		}
		listenerSaved, err := (*m.st).CreateListener(&listenerNew)
		if err != nil {
			WriteAnswer(http.StatusInternalServerError, nil, &[]error{err}, reqName, &w)
			return
		}
		WriteAnswer(http.StatusOK, listenerSaved, nil, reqName, &w)
	}
}

//TODO: write updating method
//UpdateRouter update router
func (m manage) UpdateRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqName := "Update router"
		// decoder := json.NewDecoder(r.Body)
		id := chi.URLParam(r, "id")
		if _, err := (*m.st).GetListener(id); err != nil {
			WriteAnswer(http.StatusNoContent, nil, &[]error{ErrUnableFindLisnener}, reqName, &w)
			return
		}
		// update, err := listenerNew.GetUpdateType(id)

	}
}

//DeleteRouter delete router
func (m manage) DeleteRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqName := "Delete router"
		id := chi.URLParam(r, "id")
		if err := (*m.st).DeleteListener(id); err != nil {
			WriteAnswer(http.StatusInternalServerError, nil, &[]error{err}, reqName, &w)
			return
		}
		WriteAnswer(http.StatusOK, nil, nil, reqName, &w)
	}
}
