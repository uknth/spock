/*
* @Author: Ujjwal Kanth
*
* @Email: ujjwal.kanth@unbxd.com
* @Project: spock
* @File: src/network/mw/reqid.go
* @Description: request id middleware
 */

// Package mw implements all default middleware package used by echo
// server
package mw

import (
	"config"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

const xRequestID = "X-Request-ID"

// ReqIDMiddleware sets a request id for a request
type ReqIDMiddleware struct {
	skipper func(*http.Request) bool
}

// Init initializes a request
func (rm *ReqIDMiddleware) Init(cf config.Conf) error { return nil }

// Skipper exposes method to skip a middleware
func (rm *ReqIDMiddleware) Skipper(fn func(*http.Request) bool) { rm.skipper = fn }

func (rm *ReqIDMiddleware) ServeHTTP(
	rw http.ResponseWriter, r *http.Request,
	next http.HandlerFunc,
) {
	if rm.skipper != nil && rm.skipper(r) {
		next(rw, r)
		return
	}

	if rid := r.Header.Get(xRequestID); rid == "" {
		rid = uuid.NewV4().String()
		r.Header.Set(xRequestID, rid)
		rw.Header().Set(xRequestID, rid)
	}
	next(rw, r)
}
