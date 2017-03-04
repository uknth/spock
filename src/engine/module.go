/*
* @Author: Ujjwal Kanth
*
* @Email: ujjwal.kanth@unbxd.com
* @Project: spock
* @File: src/engine/plugin.go
* @Description:
 */
package engine

// Module provides methods to plug and perform actions on request
// for a particular VirtualServer
type Module interface {
}

// NewModule returns a new object which implements Module interface
func NewModule(key string) Module {
	return nil
}
