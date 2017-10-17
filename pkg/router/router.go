package router

import (
	"net/http"
	"strings"
	"sync"

	"bitbucket.org/exonch/ch-gateway/pkg/model"
	"bitbucket.org/exonch/ch-gateway/pkg/router/middleware"
	"bitbucket.org/exonch/ch-gateway/pkg/store"
	"github.com/go-chi/chi"

	log "github.com/Sirupsen/logrus"
)

//Router keeps all http routes
type Router struct {
	*chi.Mux
	*sync.Mutex
	store *store.Store
}

var st *store.Store

//CreateRouter create and return HTTP handle router
func CreateRouter(s *store.Store) *Router {
	r := chi.NewRouter()
	router := &Router{r, &sync.Mutex{}, s}
	st = s

	//Init middlewares
	router.Use(middleware.ClearXHeaders)
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	// TODO: Add compression middleare

	//Init Not Found page handler
	router.NotFound(noRouteHandler())
	//Init root route handler
	router.HandleFunc("/", rootRouteHandler())
	//Init manage handlers
	router.Mount("/manage", CreateManageRouter())

	return router
}

//AddRoute append new http route
func (r *Router) AddRoute(target *model.Router) {
	r.Lock()
	for _, method := range target.Methods {
		method = strings.ToUpper(method)
		r.MethodFunc(method, target.ListenPath, func(w http.ResponseWriter, req *http.Request) {
			buildRoute(target, method, w, req)
		})
		log.WithFields(log.Fields{
			"ListenPath": target.ListenPath,
			"Method":     method,
		}).Debug("Route builded")
	}
	r.Unlock()
}

func buildRoute(target *model.Router, method string, w http.ResponseWriter, req *http.Request) {

}
