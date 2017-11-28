package router

import (
	"net/http"
	"strings"
	"sync"

	"bitbucket.org/exonch/ch-gateway/pkg/model"
	"bitbucket.org/exonch/ch-gateway/pkg/proxy"
	"bitbucket.org/exonch/ch-gateway/pkg/router/middleware"
	"bitbucket.org/exonch/ch-gateway/pkg/store"
	"github.com/cactus/go-statsd-client/statsd"
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
var statter *statsd.Statter

//CreateRouter create and return HTTP handle router
func CreateRouter(s *store.Store, std *statsd.Statter) *Router {
	st, statter = s, std

	//Create default router
	r := chi.NewRouter()
	router := &Router{r, &sync.Mutex{}, s}

	//Init middleware
	middleware.Statter = std
	router.Use(middleware.Logger)
	router.Use(middleware.ClearXHeaders)
	router.Use(middleware.RequestID)
	// TODO: Add compression middleware

	//Init Not Found page handler
	router.NotFound(noRouteHandler())
	//Init root route handler
	router.HandleFunc("/", rootRouteHandler())
	//Init status pahe handler
	router.HandleFunc("/status", statusHandler())
	//Init manage handlers
	router.Mount("/manage", CreateManageRouter())

	return router
}

//AddRoute append new http route
func (r *Router) AddRoute(target *model.Router) {
	// TODO: Add rate limit
	r.Lock()
	defer r.Unlock()

	for _, method := range target.Methods {
		method = strings.ToUpper(method)
		r.MethodFunc(method, target.ListenPath, func(w http.ResponseWriter, req *http.Request) {
			buildRoute(target, w, req)
		})
		log.WithFields(log.Fields{
			"ListenPath":  target.ListenPath,
			"Method":      method,
			"Roles":       target.Roles,
			"Active":      target.Active,
			"Name":        target.Name,
			"UpstreamURL": target.UpstreamURL,
		}).Debug("Route build")
	}
}

func buildRoute(target *model.Router, w http.ResponseWriter, req *http.Request) {
	p := proxy.CreateProxy(target)
	// TODO: Run before plugins
	p.ServeHTTP(w, req)
	// TODO: Run after plugins
}
