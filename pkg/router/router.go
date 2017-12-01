package router

import (
	"net/http"
	"strings"
	"sync"

	"bitbucket.org/exonch/ch-gateway/pkg/model"
	"bitbucket.org/exonch/ch-gateway/pkg/proxy"
	"bitbucket.org/exonch/ch-gateway/pkg/router/middleware"
	"bitbucket.org/exonch/ch-gateway/pkg/store"

	"git.containerum.net/ch/grpc-proto-files/auth"
	"git.containerum.net/ch/ratelimiter"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/go-chi/chi"

	log "github.com/Sirupsen/logrus"
)

//Router keeps all http routes
type Router struct {
	*chi.Mux
	*sync.Mutex
	store       *store.Store
	authClient  *auth.AuthClient
	statsClient *statsd.Statter
	rateClient  **ratelimiter.PerIPLimiter
}

var st *store.Store
var statter *statsd.Statter

//CreateRouter create and return HTTP handle router
func CreateRouter() *Router {
	//Create default router
	r := chi.NewRouter()
	x := &Router{
		Mux:        r,
		Mutex:      &sync.Mutex{},
		rateClient: new(*ratelimiter.PerIPLimiter),
	}

	//Init middleware
	r.Use(middleware.ClearXHeaders)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Rate(x.rateClient))
	// TODO: Add compression middleware

	r.NotFound(noRouteHandler())          //Init Not Found page handler
	r.HandleFunc("/", rootRouteHandler()) //Init root route handler
	// //Init status pahe handler
	// router.HandleFunc("/status", statusHandler())
	//Init manage handlers
	// router.Mount("/manage", CreateManageRouter())

	return x
}

//RegisterStore registre store interface in router
func (r *Router) RegisterStore(s *store.Store) {
	r.store = s
}

//RegisterAuth registre auth interface in router
func (r *Router) RegisterAuth(c *auth.AuthClient) {
	r.authClient = c
}

//RegisterStatsd registre statsd interface in router
func (r *Router) RegisterStatsd(s *statsd.Statter) {
	r.statsClient = s
}

//RegisterRatelimiter registre statsd interface in router
func (r *Router) RegisterRatelimiter(l *ratelimiter.PerIPLimiter) {
	*(r.rateClient) = l
}

//Start init all active routes
func (r *Router) Start() {
	st := *r.store
	listeners, err := st.GetListenerList(&model.Listener{Active: true})
	if err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Error("GetListenerList failed in router.Start")
	} else {
		for _, l := range *listeners {
			r.addRoute(&l)
		}
	}
}

//addRoute append new http route
func (r *Router) addRoute(target *model.Listener) {
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

func buildRoute(target *model.Listener, w http.ResponseWriter, req *http.Request) {
	p := proxy.CreateProxy(target)
	// TODO: Run before plugins
	p.ServeHTTP(w, req)
	// TODO: Run after plugins
}
