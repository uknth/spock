/*
* @Author: Ujjwal Kanth
*
* @Email: ujjwal.kanth@unbxd.com
* @Project: spock
* @File: src/network/handler.go
* @Description: Exposes interface to implement Generic handlers.
 */

package network

import (
	"engine"
	"network/handler"
)

/*
	Handlers which are applicable irrespective of vserver.
	Definitions
		1. Monitor Handler
*/

// TO REGISTER ANY HANDLER ADD IT TO THIS HANDLERS VARIABLE
// List of all available handlers to be attached.
// This list is used by Server.bindDefaultHandlers() method
// to bind these handlers to a url path
var handlers = []engine.Handler{
	&handler.MonitorHandler{},
}
