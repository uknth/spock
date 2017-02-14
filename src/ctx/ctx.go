/*
* @Author: Ujjwal Kanth
*
* @Email: ujjwal.kanth@unbxd.com
* @Project: spock
* @File: src/ctx/ctx.go
* @Description: Exposes application level context, initialized at startup
 */

package ctx

import (
	"config"
	"log"
	"srv"

	"github.com/pkg/errors"
)

/*
   Go doesn't provide reflection based loading so simplest way to
   add contexts in startup is to add to this array, instead of loading
   from a configuration file.
*/

// contains all available contexts to be loaded on startup
var ctxs = []C{}

// C exposes application level context's interface
type C interface {
	Init(config.Conf) error
	Name() string
}

// Load loads the context at startup
func Load(cf config.Conf) error {
	// Load Application Context
	log.Println("Application Ctx -")
	for _, c := range ctxs {
		log.Println("\tInitializing " + c.Name())
		err := c.Init(cf)
		if err != nil {
			return errors.Wrap(err, "Error initializing"+c.Name())
		}
	}

	// Load Services
	log.Println("Service Ctx-")
	for _, s := range srv.All() {
		log.Println("\tIntializing " + s.Name())
		err := s.Init(cf)
		if err != nil {
			return errors.Wrapf(err, "Error Initializing Service: "+s.Name())
		}
	}
	return nil
}
