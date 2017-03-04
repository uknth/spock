/*
* @Author: Ujjwal Kanth
*
* @Email: ujjwal.kanth@unbxd.com
* @Project: spock
* @File: src/srv/srv.go
* @Description: srv has all available external services
 */

package srv

import "config"

var services = []S{
	&eventTracker{
		name: "eventTracker",
	},
}

// S exposes service interface
type S interface {
	Init(config.Conf) error
	Name() string
}

// All returns all available services
func All() []S {
	return services
}
