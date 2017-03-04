/*
* @Author: Ujjwal Kanth
*
* @Email: ujjwal.kanth@unbxd.com
* @Project: spock
* @File: src/network/mw/header.go
* @Description: Middleware to create header
 */

package mw

import (
	"config"
	"net/http"
)

// HeaderMiddleware adds custom application specific headers to the request
type HeaderMiddleware struct {
	skipper func(*http.Request) bool
}

// Init initializes `Header Middleware`
func (hm *HeaderMiddleware) Init(cf config.Conf) error { return nil }

// Skipper gives ability to skip a middleware absed on http.Request
func (hm *HeaderMiddleware) Skipper(fn func(*http.Request) bool) { hm.skipper = fn }

// ServeHTTP handles request
func (hm *HeaderMiddleware) ServeHTTP(
	rw http.ResponseWriter, r *http.Request,
	next http.HandlerFunc) {
	if hm.skipper != nil && hm.skipper(r) {
		next(rw, r)
		return
	}
	// TODO: Get server version from configuration
	rw.Header().Set("unbxd-lb-server", "Spock 0.1")
	rw.Header().Set("unbxd-lb-type", "Spock")
	next(rw, r)
}
