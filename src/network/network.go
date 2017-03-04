/*
* @Author: Ujjwal Kanth
*
* @Email: ujjwal.kanth@unbxd.com
* @Project: spock
* @File: src/network/network.go
* @Description: network handles request
 */

// Package network is responsible for networking layer
package network

import (
	"config"
	"encoding/json"
	"engine"
	"log"
	"time"

	"reflect"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/thoas/stats"
	"github.com/tylerb/graceful"
	"github.com/urfave/negroni"
)

// Defines the state of the server
const (
	UP = iota
	DOWN
)

const (
	// Config Section
	networkSection = "network"

	// Network Contants
	portOpt = "port"

	// Control Panel Constants
	controlPanelEnabledOpt = "control-enabled"
	controlPortOpt         = "control-port"
	controlUserNameOpt     = "control-username"
	controlPasswordOpt     = "control-password"

	// Log Constants
	networkLogPathOpt = "log-path"

	// Grace Constants
	gracefulEnabled = "graceful-enabled"
	gracefulTimeout = 30

	// Profiler
	profilerEnabled = "profile-enabled"

	// Default Values
	defaultPort        = "8080"
	defaultControlPort = 9090
)

// Server describes the server handling request/response
// It also wraps http.Server
type Server struct {
	Name   string
	Port   string
	State  int
	Grace  bool
	stat   *stats.Stats
	MW     *negroni.Negroni
	Router *mux.Router
	Serf   *http.Server
}

// Init initializes Server
func (s *Server) Init(cf config.Conf) error {
	section := cf.Section(networkSection)

	// Intialize Multiplexer
	s.Router = mux.NewRouter()
	// Initialize Server
	s.Serf = &http.Server{}
	// Port number to run
	s.Serf.Addr = section.Key(portOpt).MustString(defaultPort)
	// Graceful Enabled
	s.Grace = section.Key(gracefulEnabled).MustBool(false)
	// Stat
	s.stat = stats.New()

	// Initialize Global Middlewares
	err := s.initMiddlewares(cf)
	if err != nil {
		return errors.Wrap(err, "Initializing Middlewares")
	}
	// Initialize Handlers
	err = s.initHandlers(cf)
	if err != nil {
		return errors.Wrap(err, "Initialzing Handler")
	}
	return nil
}

// initMiddlewares
func (s *Server) initMiddlewares(cf config.Conf) error {
	log.Println("Initializing Middlewares:")
	for _, mw := range middlewares {
		name := reflect.TypeOf(mw).String()
		log.Println(" --- " + name)
		err := mw.Init(cf)
		if err != nil {
			return errors.Wrap(err, "Error Initialzing "+name)
		}
	}
	return nil
}

// initHandlers
func (s *Server) initHandlers(cf config.Conf) error {
	log.Println("Initializing Handlers:")
	for _, mh := range handlers {
		name := reflect.TypeOf(mh).String()
		log.Println(" --- " + name)
		err := mh.Init(cf)
		if err != nil {
			return errors.Wrap(err, "Error Initializing "+name)
		}
	}
	return nil
}

/*
	Standard Bind Methods
*/

// Bind binds default handlers and middlewares to interface
func (s *Server) Bind(cf config.Conf) {
	section := cf.Section(networkSection)
	// Bind Default Middleware
	if section.Key(profilerEnabled).MustBool(false) {

	}
	// Bind middleware
	s.bindDefaultMiddlewares()
	// Bind stat Handler
	s.bindStatHandler()
	// Bind Handlers
	s.bindDefaultHandlers()
}

// bindStatHandler binds statistics for the node on /stats
func (s *Server) bindStatHandler() {
	s.Router.HandleFunc("/stats", func(w http.ResponseWriter,
		r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if b, err := json.Marshal(
			s.stat.Data(),
		); err == nil {
			w.Write(b)
			return
		}
		http.Error(w, "Error in parsing json",
			http.StatusInternalServerError)
		return
	})
}

// bindDefaultHandlers binds default handlers that are defined in handler.go's
// `handlers` var
func (s *Server) bindDefaultHandlers() {
	log.Println("Bind Handlers: ")
	for _, mh := range handlers {
		name := reflect.TypeOf(mh).String()
		log.Println(" ---", name, " Path: ", mh.Path(),
			" Methods: ", mh.Methods())
		if mh.Host() != engine.ALLDOMAIN {
			s.Router.Handle(mh.Path(), mh).
				Host(mh.Host()).
				Methods(mh.Methods()...)
		} else {
			s.Router.Handle(mh.Path(), mh).Methods(mh.Methods()...)
		}
	}
}

// bindDefaultMiddlewares binds default middlewares that are required in
// application
func (s *Server) bindDefaultMiddlewares() {
	// Default Negroni Classic
	s.MW = negroni.New(
		negroni.NewRecovery(),
		negroni.NewStatic(http.Dir("public")),
		s.stat,
	)
	// Just put whatever MW that you want here
	// in csv format
	log.Println("Bind CUSTOM Middlewares")
	for _, mw := range middlewares {
		s.MW.Use(mw) // Negroni Handler
	}
}

// Wire bundles everything
func (s *Server) Wire() {
	// Bind Mux's router as default handler
	s.MW.UseHandler(s.Router)
	// Bind Negroni's handler in to http.Server.Handler
	s.Serf.Handler = s.MW
}

// New returns a new network instance
func New(cf config.Conf) (*Server, error) {
	s := &Server{
		Name:  "Sever/Spock 0.1",
		State: DOWN,
	}

	// Run Init
	err := s.Init(cf)
	if err != nil {
		return nil, err
	}

	// Run bind
	s.Bind(cf)

	// Wire everything together
	s.Wire()

	// Return object if everything is alright
	return s, nil
}

// Start starts the server
func Start(server *Server) error {
	if server.Grace {
		return (&graceful.Server{
			Timeout: gracefulTimeout * time.Second,
			Server:  server.Serf,
		}).ListenAndServe()
	}
	return http.ListenAndServe(
		server.Serf.Addr,
		server.Serf.Handler,
	)
}
