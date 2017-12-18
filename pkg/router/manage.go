package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.containerum.net/ch/api-gateway/pkg/model"

	"github.com/go-chi/chi"
)

//CreateManageRouter return manage handlers
func CreateManageRouter(router *Router) http.Handler {
	r := chi.NewRouter()
	// Router headers
	r.Get("/route", getAllRouter(router))
	r.Post("/route", createRouter(router))
	r.Get("/route/{id}", getRouter(router))
	r.Put("/route/{id}", updateRouter(router))
	r.Delete("/route/{id}", deleteRouter(router))
	// r.Get("/group/{group-id}/route", getAllRouter(router))
	return r
}

func getAllRouter(router *Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Name", "get_all_router")
		st := *router.store
		listners, _ := st.GetListenerList(&model.Listener{})
		w.WriteHeader(http.StatusOK)
		WriteJSON(w, listners)
	}
}

func getRouter(router *Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Name", "get_router")
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

func createRouter(router *Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Name", "create_router")
		st := *router.store
		errAnswer := make(map[string]interface{}, 1)
		decoder := json.NewDecoder(r.Body)

		var listenerNew model.Listener
		if err := decoder.Decode(&listenerNew); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			WriteJSON(w, err)
			return
		}

		if err := listenerNew.Valid(); err != nil {
			errAnswer["Error"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			WriteJSON(w, errAnswer)
			return
		}

		listenerSaved, err := st.CreateListener(&listenerNew)
		if err != nil {
			errAnswer["Error"] = fmt.Errorf("Can not add record: Listener: %v", listenerNew)
			w.WriteHeader(http.StatusInternalServerError)
			WriteJSON(w, errAnswer)
		}

		w.WriteHeader(http.StatusOK)
		WriteJSON(w, listenerSaved)
	}
}

func updateRouter(router *Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Name", "update_router")
		st := *router.store
		decoder := json.NewDecoder(r.Body)
		id := chi.URLParam(r, "id")
		if listener, err := st.GetListener(id); err != nil {
			w.WriteHeader(http.StatusNoContent)
		} else {
			var listenerNew model.Listener
			if err := decoder.Decode(&listenerNew); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if listenerNew.Valid() != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			listenerNew.ID = listener.ID
			if err := st.UpdateListener(&listenerNew); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
}

func deleteRouter(router *Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Name", "delete_router")
		st := *router.store
		id := chi.URLParam(r, "id")
		if err := st.DeleteListener(id); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
