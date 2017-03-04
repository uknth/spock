/*
* @Author: Ujjwal Kanth
*
* @Email: ujjwal.kanth@unbxd.com
* @Project: spock
* @File: src/srv/notify.go
* @Description: notification package to handle events
 */

package srv

import "config"

type eventTracker struct {
	name string
}

// Initialize the configuration here
func (al *eventTracker) Init(cf config.Conf) error {
	// Initialize tracker/action
	return nil
}

func (al *eventTracker) Name() string {
	return al.name
}
