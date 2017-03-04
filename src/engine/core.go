/*
* @Author: Ujjwal Kanth
*
* @Email: ujjwal.kanth@unbxd.com
* @Project: spock
* @File: src/engine/core.go
* @Description: Core creates Engine
 */

package engine

import "config"

// Engine defines a wrapper to persist in memory
type Engine struct {

	// config persist a reference to config.Conf
	config config.Conf

	// Virtual Servers
	vservers map[string]*VirtualServer
}

// Handlers returns the list of handlers extracted from Vservers
func (e *Engine) Handlers() []Handler {
	return nil
}

// New returns engine's new instance
func New(cf config.Conf) (*Engine, error) {
	return nil, nil
}
