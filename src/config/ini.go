/*
* @Author: Ujjwal Kanth
*
* @Email: ujjwal.kanth@unbxd.com
* @Project: spock
* @File: src/config/ini.go
* @Description: Default Implementation of Interfaces in Config
 */

package config

import (
	"github.com/go-ini/ini"
	"github.com/pkg/errors"
)

// deaultConfigFile Path for configuration
const defaultConfigFile = "spock.ini"

var (
	errInvalidNumOfParameters = errors.New("Invalid number of parameters")
	errInvalidParameter       = errors.New("Invalid parameter")
	errConfigNotLoaded        = errors.New("Config not loaded in cache")
)

var cf Conf

// New returns a new Conf object
func New(fileName string) (Conf, error) {
	cfg, err := ini.Load(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "Error loading Configuration")
	}
	cf = &conf{
		cfg: cfg,
	}
	return cf, nil
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

// Conf implementation
type conf struct {
	cfg *ini.File
}

func (c *conf) Sections() ([]Section, error) {
	var sections = make([]Section, len(c.cfg.Sections()))
	for _, s := range c.cfg.Sections() {
		sections = append(sections, &section{
			section: s,
		})
	}
	return sections, nil
}

func (c *conf) Section(name string) (Section, error) {
	return &section{
		section: c.cfg.Section(name),
	}, nil
}

func (c *conf) Reload() error {
	return c.cfg.Reload()
}

// section wraps around ini.Section
type section struct {
	section *ini.Section
}

func (s *section) HasKey(key string) bool {
	return s.section.Haskey(key)
}

func (s *section) Keys() []Key {
	var keys = make([]Key, len(s.section.Keys()))
	for _, k := range s.section.Keys() {
		keys = append(keys, &key{
			ikey: k,
		})
	}
	return keys
}

func (s *section) KeyStrings() []string {
	return s.section.KeyStrings()
}

func (s *section) Key(name string) Key {
	return &key{
		ikey: s.section.Key(name),
	}
}

// key wraps around ini.Key
type key struct {
	ikey *ini.Key
}

func (k *key) String() string {
	return k.ikey.String()
}
func (k *key) StringS(delim string) []string {
	return k.ikey.Strings(delim)
}
func (k *key) Int() (int, error) {
	return k.ikey.Int()
}
func (k *key) MustInt(defaultValue int) int {
	return k.ikey.MustInt(defaultValue)
}
func (k *key) Bool() (bool, error) {
	return k.ikey.Bool()
}
func (k *key) MustBool(defaultValue bool) bool {
	return k.ikey.MustBool(defaultValue)
}
func (k *key) Value() string {
	return k.ikey.Value()
}
