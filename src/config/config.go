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
	Section(name string) (Section, error)

	// Reload reloads the config object
	Reload() error
}

// Section provides method to access section in conf
type Section interface {
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
