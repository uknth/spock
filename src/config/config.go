/*
* @Author: Ujjwal Kanth
*
* @Email: ujjwal.kanth@unbxd.com
* @Project: spock
* @File: src/config/config.go
* @Description: Loads config for the server
 */

// Package config is a wrapper around `package ini`.
package config

// Conf provides method to acces conf
type Conf interface {
	// Sections returns the list of sections in the configuration
	Sections() ([]Section, error)

	// Section returns section
	Section(name string) Section

	// Reload reloads the config object
	Reload() error
}

// Section provides method to access section in conf
type Section interface {
	Name() string
	// HasKey checks if a particular key is present in config
	HasKey(key string) bool

	// Keys returns the list of keys present in the section
	Keys() []Key

	// KeyStrings returns the list of keys
	KeyStrings() []string

	// Key returns the key
	Key(key string) Key
}

// Key provides method to access keys
type Key interface {
	// String returns the string form of key
	String() string

	// StringS returns list of string divided by given delimiter.
	StringS(delim string) []string

	// MustString returns string value if present, if not the default value is
	// returned
	MustString(defaultValue string) string

	// Int returns int value, if the value is not int error is returned
	Int() (int, error)

	// MustInt returns the int value, if present, if not returns the default
	// value
	MustInt(defaultValue int) int

	// Bool returns boolean value, if the value is not present, error is
	// returned
	Bool() (bool, error)

	// MustBool returns bool value if present, if not present, the default
	// vaue is returned
	MustBool(defaultValue bool) bool

	// Value returns the exact value of the key
	Value() string
}

var cf Conf

// New returns Conf object based on filename & which parser to use
func New(filename string, parser ...string) (Conf, error) {
	switch parser[0] {
	case "INI":
		return NewINI(filename)
	case "TOML":
		return NewVIPERConf(filename)
	default:
		return NewINI(filename)
	}
}

// Reload reloads default configuration
func Reload() error {
	return cf.Reload()
}

// GetSafe returns config object, but does a check before it sends
func GetSafe() (Conf, error) {
	if cf != nil {
		return cf, nil
	}
	return nil, errConfigNotLoaded
}

// Get returns the config object in memory
func Get() Conf {
	return cf
}
