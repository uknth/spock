/*
* @Author: Ujjwal Kanth
*
* @Email: ujjwal.kanth@unbxd.com
* @Project: spock
* @File: src/engine/handler.go
* @Description:
 */

// Package engine wraps entire
package engine

import (
	"config"

	"net/http"

	"github.com/urfave/negroni"
)

// Request Methods
const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	HEAD    = "HEAD"
	DELETE  = "DELETE"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
)

const (
	ALLDOMAIN = "*"
)

// Handler wraps http.Handler and adds a provision for
// Init to be called here.
// There are situations where we might need some pre-loads or configurations to
// initialize objects before Handle could handle the request.
// Inherently echo doesn't provide a mechanism to instantiate objects at server
// startup. To facilitate this, Init() method is used.
type Handler interface {
	// Composes http.Handler
	// Implement ServeHTTP(ResponseWriter, *Request)
	// Refer: https://golang.org/pkg/net/http/#Handler
	http.Handler

	// Init initializes handler
	Init(config.Conf) error

	// Path returns the path to bind this handler to
	Path() string

	// Host returns the host name to match
	Host() string

	// Methods returns the list of HTTP methods to which we need to bind
	// this handler
	Methods() []string

	// Middlewares returns list of custom middlewares to be defined for this path
	Middlewares() []negroni.Handler
}

// Handlers return list of handlers by aggregating from Virtual Servers
func Handlers() []Handler {
	// Returns list of all handlers by aggregating vservers
	// TODO: make API `return vserv.Handlers();`
	return nil
}
