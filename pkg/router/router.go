package router

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"git.containerum.net/ch/api-gateway/pkg/model"
	"git.containerum.net/ch/api-gateway/pkg/proxy"
	"git.containerum.net/ch/api-gateway/pkg/router/middleware"
	"git.containerum.net/ch/api-gateway/pkg/store"

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
	listeners   map[string]*model.Listener
	store       *store.Store
	authClient  *auth.AuthClient
	statsClient *statsd.Statter
	rateClient  **ratelimiter.PerIPLimiter
	stopSync    chan struct{}
}

//CreateRouter create and return HTTP handle router
func CreateRouter() *Router {
	return &Router{
		Mux:        chi.NewRouter(),
		Mutex:      &sync.Mutex{},
		rateClient: new(*ratelimiter.PerIPLimiter),
		listeners:  make(map[string]*model.Listener),
		stopSync:   make(chan struct{}, 1),
	}
}

//InitRoutes create main routes
func (r *Router) InitRoutes() {
	//Init middleware
	r.Use(middleware.ClearXHeaders)
	r.Use(middleware.Logger(r.statsClient))
	r.Use(middleware.RequestID)
	r.Use(middleware.Rate(r.rateClient))
	r.Use(chimid.Recoverer)
	// TODO: Add compression middleware

	r.NotFound(noRouteHandler())              //Init Not Found page handler
	r.HandleFunc("/", rootRouteHandler())     //Init root route handler
	r.Mount("/manage", CreateManageRouter(r)) //Manage handlers
	// router.HandleFunc("/status", statusHandler()) //Init status page handler
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
func (r *Router) Start(syncPeriod time.Duration) {
	st := *r.store
	listeners, err := st.GetListenerList(&model.Listener{Active: true})
	if err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Error("GetListenerList failed in router.Start")
	} else {
		for _, listener := range *listeners {
			r.addRoute(listener)
		}
	}

	//Run Sync
	go func(t time.Duration, stop chan struct{}) {
		for {
			time.Sleep(t)
			select {
			case <-stop:
				log.Debug("STOP SYNC")
				return
			default:
				r.Synchronize()
			}
		}
	}(syncPeriod, r.stopSync)
}

//Synchronize check route updates and accept it
func (r *Router) Synchronize() {
	st := *r.store
	listeners, err := st.GetListenerList(&model.Listener{Active: true})
	if err != nil {
		log.WithError(err)
	}

	listenersUpdate := make(map[string]model.Listener) //Listeners to Update
	listenersNew := make(map[string]model.Listener)    //Listeners to Create
	listenersDelete := make(map[string]model.Listener) //Listeners to Delete

	//Deep copy old routes to delete map
	for k, v := range r.listeners {
		listenersDelete[k] = *v
	}

	//Find routes to update, create or delete
	for _, listener := range *listeners {
		if listenerOld, ok := r.listeners[listener.ID]; ok {
			if listenerOld.UpdatedAt != listener.UpdatedAt {
				listenersUpdate[listener.ID] = listener
			}
		} else {
			listenersNew[listener.ID] = listener
			log.Debug(listener)
		}
		delete(listenersDelete, listener.ID)
	}

	log.WithFields(log.Fields{
		"Update": listenersUpdate,
		"New":    listenersNew,
		"Delete": listenersDelete,
	}).Debug("Starting Sync")

	r.updateRoutes(&listenersNew, &listenersUpdate, &listenersDelete)
}

//WriteJSON render JSON
func WriteJSON(w http.ResponseWriter, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	w.Write(jsonBytes)
	return nil
}

//addRoute append new http route
func (r *Router) addRoute(target model.Listener) {
	r.Lock()
	defer r.Unlock()

	//Add handler
	method := strings.ToUpper(target.Method)
	log.Debug(method)
	log.Debug(target)

	if target.OAuth {
		r.With(middleware.CheckAuthToken(r.authClient)).MethodFunc(method, target.ListenPath, func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("X-Request-Name", target.GetSnakeName())
			buildProxy(&target, w, req)
		})
	} else {
		r.MethodFunc(method, target.ListenPath, func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("X-Request-Name", target.GetSnakeName())
			buildProxy(&target, w, req)
		})
	}
	r.listeners[target.ID] = &target

	log.WithFields(log.Fields{
		"ListenPath":  target.ListenPath,
		"Method":      method,
		"Roles":       target.Roles,
		"Active":      target.Active,
		"Name":        target.Name,
		"UpstreamURL": target.UpstreamURL,
		"Auth":        target.OAuth,
	}).Debug("New route builded")
}

//updateRoute update, delete or create new routes. If only create, then append to exist chi, else create new chi route
func (r *Router) updateRoutes(listenersNew *map[string]model.Listener, listenersUpdate *map[string]model.Listener, listenersDelete *map[string]model.Listener) {
	//append route to exist chi
	if len(*listenersUpdate) == 0 && len(*listenersDelete) == 0 {
		for _, listener := range *listenersNew {
			if ok := r.Match(chi.NewRouteContext(), listener.Method, listener.ListenPath); !ok {
				r.addRoute(listener)
			} else {
				log.Debug(listener)
			}
		}
		return
	}

	//Make new routes and replace it
	r.stopSync <- struct{}{} //Stop sync
	route := CreateRouter()
	route.RegisterStore(r.store)
	route.RegisterRatelimiter(*r.rateClient)
	route.RegisterAuth(r.authClient)
	route.RegisterStatsd(r.statsClient)
	route.InitRoutes()

	for _, listener := range *listenersNew {
		if ok := r.Match(chi.NewRouteContext(), listener.Method, listener.ListenPath); !ok {
			r.addRoute(listener)
			delete(*listenersNew, listener.ID)
		}
	}
	for _, listener := range *listenersUpdate {
		if ok := r.Match(chi.NewRouteContext(), listener.Method, listener.ListenPath); !ok {
			r.addRoute(listener)
			delete(*listenersUpdate, listener.ID)
		}
	}

	route.Start(time.Second * 10)
	*r = *route //Change mux
	log.Debug("NEW MUX")
}

// TODO: Run before plugins
// TODO: Run after plugins
func buildProxy(target *model.Listener, w http.ResponseWriter, req *http.Request) {
	p := proxy.CreateProxy(target)
	p.ServeHTTP(w, req)
}
