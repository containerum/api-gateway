package router

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"bitbucket.org/exonch/ch-gateway/pkg/router/middleware"
	"github.com/go-chi/chi"
)

//Router keeps all http routes
type Router struct {
	*chi.Mux
	handlers map[string]*http.HandlerFunc
	*sync.Mutex
}

//CreateRouter create and return HTTP handle router
func CreateRouter() *Router {
	r := chi.NewRouter()
	router := &Router{r, make(map[string]*http.HandlerFunc, 0), &sync.Mutex{}}

	//Init middlewares
	router.Use(middleware.Logger)
	//Init Not Found page handler
	router.NotFound(noRouteHandler())
	//Init root route handler
	router.HandleFunc("/", rootRouteHandler())

	return router
}

//AddRoute append new http route
func (r *Router) AddRoute(id string) {
	r.Lock()
	r.Get(fmt.Sprintf("/%v", id), func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("X1"))
	})
	r.Unlock()
}

func noRouteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 2)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No route"))
	}
}

func rootRouteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello! I am API Gateway."))
	}
}
