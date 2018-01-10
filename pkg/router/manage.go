package router

import (
	"net/http"

	"git.containerum.net/ch/api-gateway/pkg/router/manage"

	"github.com/go-chi/chi"
)

//CreateManageRouter return manage handlers
func CreateManageRouter(router *Router) http.Handler {
	r := chi.NewRouter()
	// Router handlers
	// r.Get("/route", getAllRouter(router))
	// r.Post("/route", createRouter(router))
	// r.Get("/route/{id}", getRouter(router))
	// r.Put("/route/{id}", updateRouter(router))
	// r.Delete("/route/{id}", deleteRouter(router))
	// Group handlers
	// r.Get("/group", getAllGroup(router))
	// r.Post("/group", createGroup(router))
	m := manage.NewManager(router.store)
	r.Get("/route", m.GetAllRouter())
	r.Post("/route", m.CreateRouter())
	r.Get("/route/{id}", m.GetRouter())
	r.Put("/route/{id}", m.UpdateRouter())
	return r
}

//
// func getAllRouter(router *Router) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		reqName := "Get all router"
// 		st := *router.store
// 		listners, err := st.GetListenerList(&model.Listener{})
// 		if err != nil {
// 			WriteAnswer(http.StatusInternalServerError, nil, &[]error{errors.New("Unable to find listeners")}, reqName, w)
// 			return
// 		}
// 		WriteAnswer(http.StatusOK, listners, nil, reqName, w)
// 	}
// }
//
// func getRouter(router *Router) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		reqName := "Get router"
// 		id := chi.URLParam(r, "id")
// 		if l, ok := router.listeners[id]; ok {
// 			WriteAnswer(http.StatusOK, l, nil, reqName, w)
// 		} else {
// 			e := errors.New("Unable to find listener")
// 			WriteAnswer(http.StatusBadRequest, nil, &[]error{e}, reqName, w)
// 		}
// 	}
// }
//
// func createRouter(router *Router) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		reqName := "Create router"
// 		st := *router.store
// 		decoder := json.NewDecoder(r.Body)
// 		var listenerNew model.Listener
// 		if err := decoder.Decode(&listenerNew); err != nil {
// 			e := errors.New("Unable to decode listener")
// 			WriteAnswer(http.StatusBadRequest, nil, &[]error{e}, reqName, w)
// 			return
// 		}
// 		if err := listenerNew.ValidateCreate(); len(err) != 0 {
// 			WriteAnswer(http.StatusBadRequest, nil, &err, reqName, w)
// 			return
// 		}
// 		listenerSaved, err := st.CreateListener(&listenerNew)
// 		if err != nil {
// 			WriteAnswer(http.StatusInternalServerError, nil, &[]error{err}, reqName, w)
// 			return
// 		}
// 		WriteAnswer(http.StatusOK, listenerSaved, nil, reqName, w)
// 	}
// }
//
// func updateRouter(router *Router) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		reqName := "Update router"
// 		st := *router.store
// 		decoder := json.NewDecoder(r.Body)
// 		id := chi.URLParam(r, "id")
// 		if listener, err := st.GetListener(id); err != nil {
// 			WriteAnswer(http.StatusNoContent, nil, &[]error{err}, reqName, w)
// 		} else {
// 			var listenerNew model.Listener
// 			if err := decoder.Decode(&listenerNew); err != nil {
// 				WriteAnswer(http.StatusBadRequest, nil, &[]error{err}, reqName, w)
// 				return
// 			}
// 			listenerNew.ID = listener.ID
// 			//Full update
// 			var validErr []error
// 			if err := listenerNew.ValidateUpdate(); len(err) == 0 {
// 				if err := st.UpdateListener(&listenerNew, model.ListenerUpdateFull); err != nil {
// 					WriteAnswer(http.StatusInternalServerError, nil, &[]error{err}, reqName, w)
// 					return
// 				}
// 				validErr = err
// 			}
//
// 			//Activate update
// 			if err := listenerNew.ValidateUpdateActive(); len(err) != 0 {
// 				if err := st.UpdateListener(&listenerNew, model.ListenerUpdateActive); err != nil {
// 					WriteAnswer(http.StatusInternalServerError, nil, &[]error{err}, reqName, w)
// 					return
// 				}
// 				validErr = err
// 			}
//
// 			if len(validErr) != 0 {
// 				WriteAnswer(http.StatusBadRequest, nil, &validErr, reqName, w)
// 				return
// 			}
// 			WriteAnswer(http.StatusOK, listenerNew, nil, reqName, w)
// 		}
// 	}
// }

//
// func deleteRouter(router *Router) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("X-Request-Name", "delete_router")
// 		st := *router.store
// 		id := chi.URLParam(r, "id")
// 		if err := st.DeleteListener(id); err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}
// 		w.WriteHeader(http.StatusOK)
// 	}
// }

// func getAllGroup(router *Router) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("X-Request-Name", "get_all_group")
// 		st := *router.store
// 		groups, err := st.GetGroupList(&model.Group{})
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		w.WriteHeader(http.StatusOK)
// 		// WriteJSON(w, groups)
// 	}
// }
//
// func createGroup(router *Router) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("X-Request-Name", "create_group")
// 		st := *router.store
// 		errAnswer := make(map[string]interface{}, 1)
// 		decoder := json.NewDecoder(r.Body)
//
// 		var groupNew model.Group
// 		if err := decoder.Decode(&groupNew); err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			// WriteJSON(w, err)
// 			return
// 		}
//
// 		if err := groupNew.ValidateCreate(); len(err) != 0 {
// 			errAnswer["Error"] = err
// 			w.WriteHeader(http.StatusBadRequest)
// 			// WriteJSON(w, errAnswer)
// 			return
// 		}
//
// 		groupSaved, err := st.CreateGroup(&groupNew)
// 		if err != nil {
// 			errAnswer["Error"] = fmt.Errorf("Can not add record: Group: %v", groupNew)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			// WriteJSON(w, errAnswer)
// 		}
//
// 		w.WriteHeader(http.StatusOK)
// 		// WriteJSON(w, groupSaved)
// 	}
// }
