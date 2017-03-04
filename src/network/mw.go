/*
* @Author: Ujjwal Kanth
*
* @Email: ujjwal.kanth@unbxd.com
* @Project: spock
* @File: src/network/mw.go
* @Description: This file exposes global interface for Middleware
 */

package network

import (
	"engine"
	"network/mw"
)

// TO REGISTER ANY MIDDLEWARE ADD IT TO THE HF VARIABLE
// Middlewares are read and loaded from this variable
var middlewares = []engine.Middleware{
	&mw.HeaderMiddleware{},
	&mw.ReqIDMiddleware{},
}
