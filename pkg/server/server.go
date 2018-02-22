package server

import (
	"crypto/md5"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi"

	"git.containerum.net/ch/api-gateway/pkg/model"
	"git.containerum.net/ch/api-gateway/pkg/proxy"
	"git.containerum.net/ch/api-gateway/pkg/router"
	"git.containerum.net/ch/api-gateway/pkg/store"
	clickhouse "git.containerum.net/ch/api-gateway/pkg/utils/clickhouselog"

	"git.containerum.net/ch/grpc-proto-files/auth"
	"git.containerum.net/ch/ratelimiter"

	log "github.com/Sirupsen/logrus"
)

//Server keeps all https servers and do grace reload
type Server struct {
	storeClient            *store.Store
	authClient             *auth.AuthClient
	rateClient             *ratelimiter.PerIPLimiter
	clickhouseLoggerClient *clickhouse.LogClient

	*sync.Mutex
	syncPeriod time.Duration
	servers    map[int]*httpInstance
}

type httpInstance struct {
	*http.Server
	*sync.Mutex
	port     int
	stop     bool
	closed   bool
	conCount int
}

const (
	minPort          = 15000
	maxPort          = 15100
	stopServerPeriod = time.Second * 15
)

var (
	lasthash  string
	curPort   int
	targetURL = "http://localhost"
)

//New return server instace
func New(d time.Duration) *Server {
	return &Server{
		Mutex:      &sync.Mutex{},
		syncPeriod: d,
		servers:    make(map[int]*httpInstance, 1),
	}
}

//Run starts grace http server
func (s *Server) Run(port int, certFile, keyFile string) error {
	go func() {
		err := s.runMainProxy(port, certFile, keyFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	}()
	for {
		if l, ok := s.hasChanges(); ok {
			serve, _ := s.createHTTPServer(l)
			s.Lock()
			s.servers[serve.port] = serve
			curPort = serve.port
			s.Unlock()
			s.grace(serve)
		}
		time.Sleep(s.syncPeriod)
	}
}

func (s *Server) runMainProxy(port int, certFile, keyFile string) error {
	mux := chi.NewRouter()
	mux.Mount("/*", proxy.ProxyHandler(targetURL, curPort))
	log.WithField("TagetURL", targetURL).WithField("Port", curPort).Debug("Main proxy handler builed")
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: mux,
	}
	if certFile != "" && keyFile != "" {
		log.WithField("Port", port).WithField("Protocol", "HTTPS").Debug("Run main proxy")
		return server.ListenAndServeTLS(certFile, keyFile)
	}
	log.WithField("Port", port).WithField("Protocol", "HTTP").Debug("Run main proxy")
	return server.ListenAndServe()
}

func (s *Server) grace(serve *httpInstance) {
	go serve.run()
	time.Sleep(stopServerPeriod)
	s.Lock()
	defer s.Unlock()
	for _, serv := range s.servers {
		if serv.port != serve.port {
			serv.Lock()
			defer serv.Unlock()
			serv.stop = true
		}
	}
}

func (inst *httpInstance) run() {
	inst.ConnState = func(con net.Conn, state http.ConnState) {
		inst.Lock()
		switch state {
		case http.StateNew:
			inst.conCount++
		case http.StateClosed:
			inst.conCount--
		}
		inst.Unlock()
	}
	go func(ins *httpInstance) {
		for {
			time.Sleep(stopServerPeriod)
			ins.Lock()
			if ins.stop && ins.conCount == 0 {
				inst.Close()
				inst.closed = true
			}
			ins.Unlock()
		}
	}(inst)
	inst.ListenAndServe()
}

func (s *Server) createHTTPServer(l *[]model.Listener) (*httpInstance, error) {
	route, err := router.CreateRouterNEW(*l, s.storeClient, s.authClient, s.rateClient, s.clickhouseLoggerClient)
	if err != nil {
		return nil, err
	}
	port := s.getPort()
	return &httpInstance{
		Server:   &http.Server{Addr: fmt.Sprintf(":%v", port), Handler: route},
		Mutex:    &sync.Mutex{},
		port:     port,
		conCount: 0,
	}, nil
}

func (s *Server) hasChanges() (*[]model.Listener, bool) {
	active := true
	listeners, err := (*s.storeClient).GetListenerList(&active)
	if err != nil {
		return nil, false //TODO: write err log
	}
	var hashdata []byte
	for _, l := range *listeners {
		key := fmt.Sprintf("%s_%s", l.Name, l.UpdatedAt.String())
		hashdata = append(hashdata, key...)
	}
	hash := fmt.Sprintf("%x", md5.Sum(hashdata))
	if hash != lasthash {
		lasthash = hash
		return listeners, true
	}
	return listeners, false
}

func (s *Server) getPort() int {
	max := minPort - 1
	for i := range s.servers {
		if max < i {
			max = i
		}
	}
	if max = max + 1; max > maxPort {
		max = minPort
	}
	return max
}
