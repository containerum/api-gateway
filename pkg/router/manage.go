package router

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"bitbucket.org/exonch/ch-gateway/pkg/model"

	log "github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
)

// TODO: move to sub-package
// TODO: rename router to gateway

//CreateManageRouter return manage handlers
func CreateManageRouter() http.Handler {
	r := chi.NewRouter()
	// Router headers
	r.Get("/route", getAllRouter)
	r.Get("/route/{id}", getRouter)
	r.Put("/route/{id}", updateRouter)
	r.Delete("/route/{id}", removeRouter)
	r.Get("/group/{group-id}/route", getAllRouter)
	r.Post("/group/{group-id}/route", addRouter)
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
		}).Error("getAllRouter error")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

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
	buf, _ := ioutil.ReadAll(r.Body)
	route, err := model.CreateRouterFromJSON(buf)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Add new router in DB
	if err = (*st).AddRouter(route); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(route.ConvertToJSON())
}

func updateRouter(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update route"))
}

func removeRouter(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Remove route"))
}
