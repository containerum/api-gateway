package server

import (
	"fmt"
	"net/http"

	"git.containerum.net/ch/api-gateway/pkg2/model"
	middle "git.containerum.net/ch/api-gateway/pkg2/server/middleware"
	"github.com/gin-gonic/gin"
)

//Server keeps HTTP sever and it configs
type Server struct {
	http.Server
	routes *model.Routes
	config *model.Config
}

//New return configurated server with all handlers
func New(c *model.Config, r *model.Routes) (*Server, error) {
	handlers, err := createHandler(c, r)
	if err != nil {
		return nil, err
	}
	return &Server{
		config: c,
		routes: r,
		Server: http.Server{
			Addr:    fmt.Sprintf(":%v", c.Port),
			Handler: handlers,
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

func createHandler(c *model.Config, r *model.Routes) (http.Handler, error) {
	router := gin.New()
	//Add middlewares
	router.Use(middle.Logger(), middle.Recovery(), middle.Cors())
	router.Use(middle.ClearXHeaders(), middle.SetRequestID())
	router.Use(middle.CheckUserClientHeader(), middle.SetMainUserXHeaders())
	//Add routes
	for _, route := range r.Routes {
		if !route.Active {
			continue
		}
		router.Handle(route.Method, route.Listen, func(c *gin.Context) {
			p := createProxy(&route)
			p.ServeHTTP(c.Writer, c.Request)
		})
	}
	return router, nil
}
