package router

import (
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"git.containerum.net/ch/api-gateway/pkg/model"
	"git.containerum.net/ch/api-gateway/pkg/proxy"
	"git.containerum.net/ch/api-gateway/pkg/router/middleware"
	"git.containerum.net/ch/api-gateway/pkg/store"

	clickhouse "git.containerum.net/ch/api-gateway/pkg/utils/clickhouselog"
	"git.containerum.net/ch/grpc-proto-files/auth"
	"git.containerum.net/ch/ratelimiter"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/go-chi/chi"
	chimid "github.com/go-chi/chi/middleware"

	log "github.com/Sirupsen/logrus"
)

//Router keeps all http routes
type Router struct {
	*chi.Mux
	*sync.Mutex
	listeners              map[string]*model.Listener
	store                  *store.Store
	authClient             *auth.AuthClient
	statsClient            *statsd.Statter
	rateClient             *ratelimiter.PerIPLimiter
	clickhouseLoggerClient *clickhouse.LogClient
	stopSync               chan struct{}
}

const (
	syncPeriod = time.Second * 30
)

var (
	//ErrUnbaleGetListenersForSync - error when unable to get listeners
	ErrUnbaleGetListenersForSync = errors.New("Unable to get listeners for Synchronization")
)

//CreateRouter create and return HTTP handle router
func CreateRouter(router *Router) *Router {
	if router == nil {
		return &Router{
			Mux:                    chi.NewRouter(),
			Mutex:                  &sync.Mutex{},
			rateClient:             &ratelimiter.PerIPLimiter{},
			clickhouseLoggerClient: &clickhouse.LogClient{},
			listeners:              make(map[string]*model.Listener),
			stopSync:               make(chan struct{}, 1),
		}
	}
	return &Router{
		Mux:                    chi.NewRouter(),
		Mutex:                  &sync.Mutex{},
		store:                  router.store,
		rateClient:             router.rateClient,
		statsClient:            router.statsClient,
		authClient:             router.authClient,
		clickhouseLoggerClient: router.clickhouseLoggerClient,
		listeners:              make(map[string]*model.Listener),
		stopSync:               make(chan struct{}, 1),
	}
}

//InitRoutes create main routes
//TODO: Add compression middleware
func (r *Router) InitRoutes() {
	//Init middleware
	r.Use(middleware.ClearXHeaders)
	r.Use(middleware.TranslateUserXHeaders)
	r.Use(middleware.CheckRequiredXHeaders)
	r.Use(middleware.Logger(r.statsClient, r.clickhouseLoggerClient))
	r.Use(middleware.RequestID)
	r.Use(middleware.Rate(r.rateClient))
	r.Use(chimid.Recoverer)

	r.With(middleware.CheckAuthToken(r.authClient)).Mount("/manage", CreateManageRouter(r)) //Manage handlers
	r.NotFound(noRouteHandler())                                                            //Init Not Found page handler
	r.HandleFunc("/", rootRouteHandler())                                                   //Init root route handler
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
	r.rateClient = l
}

//RegisterClickhouseLogger registre clickhouse logger client
func (r *Router) RegisterClickhouseLogger(cl *clickhouse.LogClient) {
	r.clickhouseLoggerClient = cl
}

//Start init all active routes
func (r *Router) Start() {
	go r.runSynchronization()
}

//Stop -stops sync
func (r *Router) Stop() {
	r.stopSync <- struct{}{}
}

func (r *Router) runSynchronization() {
	for {
		select {
		case <-r.stopSync:
			log.Info("Synchronization stopped")
			return
		default:
			r.synchronize()
		}
		time.Sleep(syncPeriod)
	}
}

//Synchronize check route updates and accept it
func (r *Router) synchronize() {
	st := *r.store
	active := true
	listenersUpdate := make(map[string]model.Listener) //Listeners to Update
	listenersNew := make(map[string]model.Listener)    //Listeners to Create
	listenersDelete := make(map[string]model.Listener) //Listeners to Delete
	listeners, err := st.GetListenerList(&active)
	if err != nil {
		log.WithError(err).Warn(ErrUnbaleGetListenersForSync)
		return
	}
	copyListenersMap(r.listeners, &listenersDelete)
	//Find routes to update, create or delete
	for _, listener := range *listeners {
		if listenerOld, ok := r.listeners[listener.ID]; ok {
			if listenerOld.UpdatedAt.Unix() != listener.UpdatedAt.Unix() {
				listenersUpdate[listener.ID] = listener
			}
		} else {
			listenersNew[listener.ID] = listener
		}
		delete(listenersDelete, listener.ID)
	}
	log.WithFields(log.Fields{
		"New routes":    len(listenersNew),
		"Update routes": len(listenersUpdate),
		"Delete routes": len(listenersDelete),
	}).Debug("Synchronization")
	r.updateRoutes(&listenersNew, &listenersUpdate, &listenersDelete)
}

//updateRoute update, delete or create new routes. If only create, then append to exist chi, else create new chi route
func (r *Router) updateRoutes(listenersNew *map[string]model.Listener, listenersUpdate *map[string]model.Listener, listenersDelete *map[string]model.Listener) {
	//append route to exist chi
	if len(*listenersUpdate) == 0 && len(*listenersDelete) == 0 {
		for _, listener := range *listenersNew {
			if ok := r.Match(chi.NewRouteContext(), listener.Method, listener.ListenPath); !ok {
				r.addRoute(listener)
			}
		}
		return
	}
	//Make new routes and replace old
	r.Stop()
	route := CreateRouter(r)
	route.InitRoutes()
	for _, listener := range appendListenersMap(*listenersNew, *listenersUpdate) {
		if ok := r.Match(chi.NewRouteContext(), listener.Method, listener.ListenPath); !ok {
			r.addRoute(listener)
		}
	}
	route.Start()
	*r = *route //Change mux
	log.WithField("Route count", len(route.Routes())).Info("New mux started")
}

//addRoute append new http route
func (r *Router) addRoute(target model.Listener) {
	r.Lock()
	defer r.Unlock()
	method := strings.ToUpper(target.Method)
	if target.OAuth {
		r.With(middleware.CheckAuthToken(r.authClient)).MethodFunc(method, target.ListenPath, func(w http.ResponseWriter, req *http.Request) {
			buildProxy(&target, w, req)
		})
	} else {
		r.MethodFunc(method, target.ListenPath, func(w http.ResponseWriter, req *http.Request) {
			buildProxy(&target, w, req)
		})
	}
	r.listeners[target.ID] = &target
	log.WithFields(log.Fields{
		"ListenPath":  target.ListenPath,
		"Method":      method,
		"Statsd":      target.GetSnakeName(),
		"Name":        target.Name,
		"UpstreamURL": target.UpstreamURL,
		"Auth":        target.OAuth,
		"StripPath":   target.StripPath,
		"Group":       target.Group.Name,
	}).Info("New route")
}

// TODO: Run before plugins
// TODO: Run after plugins
func buildProxy(target *model.Listener, w http.ResponseWriter, req *http.Request) {
	w.Header().Set("X-Request-Name", target.GetSnakeName())
	w.Header().Set("X-Gateway-ID", target.ID)
	w.Header().Set("X-Upstream", target.UpstreamURL)
	p := proxy.CreateProxy(target, w.Header())
	p.ServeHTTP(w, req)
}

func copyListenersMap(src map[string]*model.Listener, dst *map[string]model.Listener) {
	for k, v := range src {
		(*dst)[k] = *v
	}
}

func appendListenersMap(arrs ...map[string]model.Listener) (res map[string]model.Listener) {
	for _, arr := range arrs {
		for k, v := range arr {
			res[k] = v
		}
	}
	return
}
