/*
* @Author: Ujjwal Kanth
*
* @Email: ujjwal.kanth@unbxd.com
* @Project: spock
* @File: src/engine/vserver.go
* @Description:
 */

package engine

import (
	"config"

	"github.com/sirupsen/logrus"
)

// VirtualServer wraps configuration & path for a request
type VirtualServer struct {
	// Associated domain for the server
	Domain string

	// Logging
	Logger *logrus.Logger

	// Conf
	Config config.Conf

	// Plugin
	module Module

	// HealthCheck
	HC HealthCheck

	// Enabled
	Enabled bool

	// Warnings
	Warnings []string
}

// Init initailizes VirtualServer
func (v *VirtualServer) Init(cf config.Conf) error {
	return nil
}

// Handlers returns list of handlers assciated with this VirtualServer
func (v *VirtualServer) Handlers() []Handler {
	return nil
}

// NewVServer returns a new instance of VirtualServer
func NewVServer(cf config.Conf) *VirtualServer {
	return nil
}
