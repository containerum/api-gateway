package router

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/exonch/ch-gateway/pkg/model"

	"github.com/go-chi/chi"
)

//CreateManageRouter return manage handlers
func CreateManageRouter(router *Router) http.Handler {
	r := chi.NewRouter()
	// Router headers
	r.Get("/route", getAllRouter(router))
	r.Get("/route/{id}", getRouter(router))
	r.Put("/route/{id}", updateRouter(router))
	r.Get("/group/{group-id}/route", getAllRouter(router))
	return r
}

func getAllRouter(router *Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		st := *router.store
		listners, _ := st.GetListenerList(&model.Listener{})
		w.WriteHeader(http.StatusOK)
		WriteJSON(w, listners)
	}
}

func getRouter(router *Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if l, ok := router.listeners[id]; ok {
			w.WriteHeader(http.StatusOK)
			WriteJSON(w, l)
		} else {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
}

func updateRouter(router *Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		st := *router.store
		decoder := json.NewDecoder(r.Body)
		id := chi.URLParam(r, "id")
		if listener, ok := router.listeners[id]; ok {
			var listenerNew model.Listener
			if err := decoder.Decode(&listenerNew); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			listener.Method = listenerNew.Method
			if err := st.UpdateListener(listener); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
}
