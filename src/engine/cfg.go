/*
* @Author: Ujjwal Kanth
*
* @Email: ujjwal.kanth@unbxd.com
* @Project: spock
* @File: src/engine/cfg.go
* @Description: C stores engine's config
 */

package engine

import (
	"config"
	"io/ioutil"

	"strings"

	"github.com/pkg/errors"
)

// Config Folders
var siteAvailableFolders = "./conf"

var (
	errInvalidConfig     = errors.New("Invalid Configuration file(s).")
	errServerNameMissing = errors.New("Server name not found")
)

// Configuration constants
const (
	virtualServerSec = "virtual-server"
	serverNameOpt    = "server-name"
	accessLogOpt     = "access-log"
	errorLogOpt      = "error-log"
	indexFileOpt     = "index"
	pathPrefixSec    = "path-"
	moduleOpt        = "module"
	defaultModule    = "file"
)

// Walk the directory conf and get the list of virtual servers

// path wraps properties specific to path for a virtual server
type path struct {
	uri           string
	module        Module
	configSection config.Section
}

func (p *path) String() string {
	return p.uri
}

// VServerConfig wraps config.Conf to expose virtual server specific
// configuration.
type VServerConfig struct {
	fileName string
	config   config.Conf
	paths    []path
}

// Init initializes the object
func (vsc *VServerConfig) Init() (*VServerConfig, error) {
	// Get all sections
	sections, err := vsc.config.Sections()
	if err != nil {
		return nil, errors.Wrap(err, "Error getting sections")
	}
	for _, section := range sections {
		// Path
		uri := strings.TrimPrefix(section.Name(), pathPrefixSec)
		// Module
		module := NewModule(section.Key(moduleOpt).
			MustString(defaultModule))
		// Also add this to our structure
		vsc.paths = append(vsc.paths, path{
			uri:           uri,
			module:        module,
			configSection: section,
		})
	}
	return vsc, nil
}

func (vsc *VServerConfig) domain() string {
	return vsc.config.Section(virtualServerSec).Key(serverNameOpt).String()
}

// Domain returns the domain name
func (vsc *VServerConfig) Domain() (string, error) {
	if !vsc.config.Section(virtualServerSec).HasKey(serverNameOpt) {
		return "", errServerNameMissing
	}
	return vsc.domain(), nil
}

// AccessLog returns the log file for the virtualserver
func (vsc *VServerConfig) AccessLog() string {
	return vsc.config.Section(virtualServerSec).Key(accessLogOpt).
		MustString(vsc.domain() + "_access.log")
}

// ErrorLog returns the log file to log Error for virtualserver
func (vsc *VServerConfig) ErrorLog() string {
	return vsc.config.Section(virtualServerSec).Key(errorLogOpt).
		MustString(vsc.domain() + "_error.log")
}

// Paths returns the associated path for this balancer
func (vsc *VServerConfig) Paths() []string {
	var paths = make([]string, 0)
	for _, path := range vsc.paths {
		paths = append(paths, path.String())
	}
	return paths
}

// VServerConfigs returns a new object of Config
func VServerConfigs() ([]VServerConfig, error) {
	var vsconfs = make([]VServerConfig, 0)
	// Read Directory Files
	files, err := ioutil.ReadDir(siteAvailableFolders)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading directory")
	}
	// We have list of files, lets read the config from them
	for _, file := range files {
		cf, err := config.New(siteAvailableFolders + file.Name())
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing config"+file.Name())
		}
		vs, err := (&VServerConfig{
			config:   cf,
			fileName: file.Name(),
			paths:    make([]path, 0),
		}).Init()
		if err != nil {
			return nil, errors.Wrap(err, "Error creating vserverconfig")
		}
		vsconfs = append(vsconfs, *vs)
	}
	return vsconfs, nil
}
