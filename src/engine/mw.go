/*
* @Author: Ujjwal Kanth
*
* @Email: ujjwal.kanth@unbxd.com
* @Project: spock
* @File: src/engine/mw.go
* @Description:
 */

package engine

import (
	"config"
	"net/http"

	"github.com/urfave/negroni"
)

// Middleware wraps echo's `MiddlewareFunc` and adds a provision for
// Init to be called here.
// There are situations where we might need some pre-loads or configurations to
// initialize objects before Handle could handle the request.
// Inherently echo doesn't provide a mechanism to instantiate objects at server
// startup. To facilitate this, Init() method is used.
type Middleware interface {
	// Init initializes middleware
	Init(cf config.Conf) error
	// Skipper gives an ability to conditionally skip a middleware
	Skipper(func(*http.Request) bool)
	// Negroni's handler
	// Refer: github.com/urfave/negroni/blob/master/negroni.go#L15
	negroni.Handler
}
