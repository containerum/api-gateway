package router

import (
	"fmt"
	"net/http"

	"bitbucket.org/exonch/ch-gateway/pkg/model"

	"github.com/go-chi/chi"

	log "github.com/Sirupsen/logrus"
)

//CreateManageRouter return manage handlers
func CreateManageRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getAllRouter)
	r.Post("/", addRouter)
	r.Get("/{id}", getRouter)
	r.Put("/{id}", updateRouter)
	r.Delete("/{id}", removeRouter)
	return r
}

func getAllRouter(w http.ResponseWriter, r *http.Request) {
	var rs *[]model.Router
	var err error
	active := r.URL.Query().Get("active")

	if active == "" {
		log.Debug("Call GetRoutesList in AllRouter")
		rs, err = (*st).GetRoutesList()
	} else {
		switch active {
		case "0", "t", "true":
			log.Debug("Call GetRoutesList active true in AllRouter")
			rs, err = (*st).GetRoutesListByActivation(true)
		case "1", "f", "false":
			log.Debug("Call GetRoutesList active false in getAllRouter")
			rs, err = (*st).GetRoutesListByActivation(false)
		default:
			err = fmt.Errorf("Couldn't parse Active: %v", active)
			log.WithField("Err", err).Debug("Write error in getAllRouter")
		}
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"Status": http.StatusInternalServerError,
			"Err":    err,
		}).Debug("getAllRouter bad answer")
		return
	}

	w.WriteHeader(http.StatusOK)
	for _, rr := range *rs {
		w.Write([]byte(rr.ID))
	}
	// TODO: Write answer in log
	log.WithFields(log.Fields{
		"Status": http.StatusOK,
	}).Debug("getAllRouter good answer")
}

func getRouter(w http.ResponseWriter, r *http.Request) {
	rs, err := (*st).GetRouter(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(rs.ID))
}

func addRouter(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create route"))
}

func updateRouter(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update route"))
}

func removeRouter(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Remove route"))
}
