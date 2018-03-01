package server

import (
	"fmt"
	slog "log"
	"net/http"
	"time"

	"git.containerum.net/ch/api-gateway/pkg2/model"
	middle "git.containerum.net/ch/api-gateway/pkg2/server/middleware"
	"git.containerum.net/ch/grpc-proto-files/auth"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

//Server keeps HTTP sever and it configs
type Server struct {
	http.Server
	routes *model.Routes
	config *model.Config
	auth   *auth.AuthClient
}

//New return configurated server with all handlers
func New(c *model.Config, r *model.Routes, a *auth.AuthClient) (*Server, error) {
	handlers, err := createHandler(c, r, a)
	if err != nil {
		return nil, err
	}
	return &Server{
		config: c,
		routes: r,
		auth:   a,
		Server: http.Server{
			Addr:              fmt.Sprintf("0.0.0.0:%v", c.Port),
			Handler:           handlers,
			ReadTimeout:       4 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
			WriteTimeout:      8 * time.Second,
			ErrorLog:          slog.New(log.New().Writer(), "server", 0),
		},
	}, nil
}

//Start return http or https ListenServer
func (s *Server) Start() error {
	if s.config.TLS.Enable {
		return s.Server.ListenAndServeTLS(s.config.TLS.Cert, s.config.TLS.Key)
	}
	return s.ListenAndServe()
}

func createHandler(c *model.Config, r *model.Routes, auth *auth.AuthClient) (http.Handler, error) {
	log.Debug(r.Routes)
	router := gin.New()
	limiter := middle.CreateLimiter(c.Rate.Limit)
	//Add middlewares
	router.Use(middle.Logger(), middle.Recovery(), middle.Cors())
	router.Use(limiter.Limit())
	router.Use(middle.ClearXHeaders(), middle.SetRequestID())
	router.Use(middle.CheckUserClientHeader(), middle.SetMainUserXHeaders())
	//Add routes
	for _, route := range r.Routes {
		if !route.Active {
			continue
		}
		if c.Auth.Enable {
			router.Handle(route.Method, route.Listen, middle.CheckAuth(route.Roles, auth), proxyHandler(route))
		} else {
			router.Handle(route.Method, route.Listen, proxyHandler(route))
		}
	}
	return router, nil
}
