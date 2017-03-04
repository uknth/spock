/*
* @Author: Ujjwal Kanth
*
* @Email: ujjwal.kanth@unbxd.com
* @Project: spock
* @File: src/network/handler/monitor.go
* @Description:
 */
package handler

import (
	"config"
	"engine"
	"net/http"

	"github.com/urfave/negroni"
)

const (
	healthCheckSection = "health-check"
	healthCheckPathOpt = "path"
)

// MonitorHandler exposes "/monitor" path
type MonitorHandler struct {
	path string
}

// Init initializes handler
func (mh *MonitorHandler) Init(cf config.Conf) error {
	// Read HealthCheck URL from config. Default value is /monitor
	mh.path = cf.Section(healthCheckSection).
		Key(healthCheckPathOpt).MustString("/monitor")
	return nil
}

// ServeHTTP is implemented for it to confer to http.Handler interface
func (mh *MonitorHandler) ServeHTTP(w http.ResponseWriter,
	r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("alive"))
	return
}

// Path returns the path to bind this handler to
func (mh *MonitorHandler) Path() string {
	return mh.path
}

// Methods returns the list of HTTP methods to which we need to bind
// this handler
func (mh *MonitorHandler) Methods() []string {
	return []string{engine.HEAD, engine.GET}
}

// Middlewares returns list of custom middlewares to be defined for this path
func (mh *MonitorHandler) Middlewares() []negroni.Handler {
	return nil
}

// Host returns the domain names it is applicable to
func (mh *MonitorHandler) Host() string {
	return engine.ALLDOMAIN
}
